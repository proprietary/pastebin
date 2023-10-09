package main

import (
	"crypto/sha256"
	"errors"
	badger "github.com/dgraph-io/badger/v4"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

type Slug string

const MAX_PASTEBIN_BYTES int = 10 << 20

func openDb() *badger.DB {
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

func SavePastebin(db *badger.DB, text []byte) (Slug, error) {
	if utf8.ValidString(string(text)) == false {
		return "", ErrNotAUtf8String
	}
	if len(text) > MAX_PASTEBIN_BYTES {
		return "", ErrMaxPastebinSizeExceeded
	}
	slug := GenerateTextSlug(text)
	err := db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(slug), text)
	})
	return slug, err
}

func LookupPastebin(db *badger.DB, slug Slug) (string, error) {
	var textData []byte = nil
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(slug))
		if err != nil {
			return err
		}
		textData = make([]byte, item.ValueSize())
		err = item.Value(func(val []byte) error {
			copy(textData, val)
			return nil
		})
		return err
	})
	if err != nil {
		return "", err
	}
	// assume the data is valid if it was persisted to begin with
	text := string(textData)
	return text, err
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
)
