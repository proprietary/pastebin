package text_store

import (
	"crypto/sha256"
	"errors"
	badger "github.com/dgraph-io/badger/v4"
	"github.com/proprietary/pastebin/pastebin_record"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

type Slug string

const MAX_PASTEBIN_BYTES int = 10 << 20

func OpenDb() *badger.DB {
	dbPath := os.Getenv("SO_LIBHACK_PASTE__DB_PATH")
	if len(dbPath) == 0 {
		log.Fatal(`Missing environment variable for path of database: "SO_LIBHACK_PASTE__DB_PATH"`)
	}
	db, err := badger.Open(badger.DefaultOptions(dbPath))
	if err != nil {
		log.Fatal("Fail to open database at \"", dbPath, "\":", err)
	}
	return db
}

func SavePastebin(db *badger.DB, text []byte, exp time.Time) (Slug, error) {
	if utf8.ValidString(string(text)) == false {
		return "", ErrNotAUtf8String
	}
	if len(text) > MAX_PASTEBIN_BYTES {
		return "", ErrMaxPastebinSizeExceeded
	}
	now := time.Now()
	if now.After(exp) {
		return "", ErrInvalidExpiration
	}
	slug := GenerateTextSlug(text)
	pb := pastebin_record.PastebinRecord{
		Body:               string(text),
		TimeCreated:        timestamppb.New(now),
		Expiration:         timestamppb.New(exp),
		MimeType:           nil,
		SyntaxHighlighting: nil,
	}
	pastebinBytes, err := proto.Marshal(pb.ProtoReflect().Interface())
	if err != nil {
		return "", err
	}
	entry := badger.NewEntry([]byte(slug), pastebinBytes).WithTTL(exp.Sub(now))
	err = db.Update(func(txn *badger.Txn) error {
		return txn.SetEntry(entry)
	})
	return slug, err
}

func LookupPastebin(db *badger.DB, slug Slug) (string, error) {
	var pbData []byte = nil
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(slug))
		if err != nil {
			return err
		}
		pbData = make([]byte, item.ValueSize())
		err = item.Value(func(val []byte) error {
			copy(pbData, val)
			return nil
		})
		return err
	})
	if err != nil {
		return "", err
	}
	// assume the data is a valid structure if it was persisted to begin with
	var pb pastebin_record.PastebinRecord
	err = proto.Unmarshal(pbData, &pb)
	if err != nil {
		return "", err
	}
	return pb.GetBody(), err
}

/// KeepalivePastebin issues a new expiration for a pastebin.
func KeepalivePastebin(db *badger.DB, slug Slug, newExpiry time.Time) error {
	var pbData []byte = nil
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(slug))
		if err != nil {
			return err
		}
		pbData = make([]byte, item.ValueSize())
		err = item.Value(func(val []byte) error {
			copy(pbData, val)
			return nil
		})
		return err
	})
	if err != nil {
		return err
	}
	var pb pastebin_record.PastebinRecord
	if err = proto.Unmarshal(pbData, &pb); err != nil {
		return err
	}
	if newExpiry.Before(pb.Expiration.AsTime()) {
		return ErrInvalidExpiration
	}
	pb.Expiration = timestamppb.New(newExpiry)
	updatedPb, err := proto.Marshal(&pb)
	if err != nil {
		return err
	}
	err = db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(slug), updatedPb)
	})
	return err
}

// / RunMaintenance maintains integrity of the storage, including removing expired pastes.
func RunMaintenance(db *badger.DB) error {
	now := time.Now()
	// iterate through kv's and delete items according to various criteria
	err := db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true
		it := txn.NewIterator(opts)
		defer it.Close()
		thisPastebin := new(pastebin_record.PastebinRecord)
		for it.Rewind(); it.Valid(); it.Next() {
			shouldDelete := false
			err := it.Item().Value(func(pastebinBytes []byte) error {
				err := proto.Unmarshal(pastebinBytes, thisPastebin)
				if err != nil {
					// not this type of protobuf? it's trash; mark it for deletion
					shouldDelete = true
					return nil
				}
				exp := thisPastebin.GetExpiration()
				shouldDelete = exp != nil && exp.AsTime().After(now)
				return nil
			})
			if err != nil {
				return err
			}
			if shouldDelete {
				err := txn.Delete(it.Item().Key())
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	// Run Badger GC
	// See: https://dgraph.io/docs/badger/get-started/#garbage-collection
	const BADGER_GC_DURATION time.Duration = 5 * time.Minute
	ticker := time.NewTicker(BADGER_GC_DURATION)
	defer ticker.Stop()
	for range ticker.C {
	again:
		err := db.RunValueLogGC(0.7)
		if err == nil {
			goto again
		}
	}
	return nil
}

// / GenerateTextSlug uniquely encodes a body of bytes into a short string identifier.
func GenerateTextSlug(text []byte) Slug {
	h := sha256.New()
	h.Write(text)
	sum := h.Sum(nil)
	return makeSlug(sum)
}

// / makeSlug encodes arbitrary bytes to a short string.
func makeSlug(digest []byte) Slug {
	var availableCharacters []rune = []rune("acdefghjklmnpqrstuvwxyz023456789")
	const slugLength int = 8
	var outputString strings.Builder
	var count int = 0
	for _, b := range digest {
		if count >= slugLength {
			break
		}
		choice := int(b) % len(availableCharacters)
		outputString.WriteRune(availableCharacters[choice])
		count++
	}
	return Slug(outputString.String())
}

var (
	ErrNotAUtf8String          = errors.New("Not a UTF-8 string")
	ErrMaxPastebinSizeExceeded = errors.New("Pastebin text too long")
	ErrInvalidExpiration       = errors.New("Attempt to create a pastebin that should have already expired")
)
