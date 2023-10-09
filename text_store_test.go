package main

import (
	badger "github.com/dgraph-io/badger/v4"
	"testing"
)

func TestMakeSlug(t *testing.T) {
	helloWorldString := "hello world"
	someData := []byte(helloWorldString)
	slug := GenerateTextSlug(someData)
	// produces some output
	if len(slug) == 0 {
		t.Fatalf(`Slug generated from "%s" is empty ("%s")`, helloWorldString, slug)
	}
	// produces the same output repeatedly
	slug2 := GenerateTextSlug([]byte(helloWorldString))
	if len(slug2) == 0 {
		t.Fatalf(`Slug generated from "%s" is empty ("%s")`, helloWorldString, slug2)
	}
	if string(slug) != string(slug2) {
		t.Fatalf(`First invocation generated "%s" yet second invocation on the same data produced "%s"`, slug2, slug)
	}
}

func FuzzMakeSlug(f *testing.F) {
	f.Fuzz(func(t *testing.T, someString string) {
		slug := GenerateTextSlug([]byte(someString))
		if len(slug) == 0 {
			t.Fatalf(`Slug generated from "%s" is empty ("%s")`, someString, slug)
		}
	})
}

func TestLookupSavePastebins(t *testing.T) {
	const examplePaste string = `
Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
`
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	if err != nil {
		t.Fatalf("Failure to create Badger in-memory database for testing: %v", err)
	}
	defer db.Close()
	slug, err := SavePastebin(db, []byte(examplePaste))
	if err != nil {
		t.Fatal(err)
	}
	if len(slug) == 0 {
		t.Fatalf(`Slug generated from "%s" is empty: %s`, examplePaste, slug)
	}
	retrievedPaste, err := LookupPastebin(db, slug)
	if err != nil {
		t.Fatal(err)
	}
	if retrievedPaste != examplePaste {
		t.Fatalf(`Expected: "%s", Got: "%s"`, examplePaste, retrievedPaste)
	}
}
