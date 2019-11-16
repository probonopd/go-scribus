package scribus

import (
	"os"

	"testing"
)

func TestReadScribusFile(t *testing.T) {
	scribusDocument, err := NewScribusDocumentFromFile("Document-1.sla")
	if err != nil {
		t.Errorf("error: %v", err)
	}

	// Check if we got the correct version back
	if scribusDocument.Version != "1.5.1.svn" {
		t.Errorf("scribusDocument.Version was incorrect, got: %v, want: %v.", scribusDocument.Version, "1.5.1.svn")
	}

}

func TestWriteScribusFile(t *testing.T) {
	scribusDocument, err := NewScribusDocumentFromFile("Document-1.sla")
	if err != nil {
		t.Errorf("error: %v", err)
	}

	testfile := "test.sla"
	err = scribusDocument.WriteScribusFile(testfile)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	// Test if the written file exists
	if _, err := os.Stat(testfile); os.IsNotExist(err) {
		t.Errorf("error: %v does not exist", testfile)
	}

	// Test if we can read the written file
	scribusDocument, err = NewScribusDocumentFromFile(testfile)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	// Check if we got the correct version back
	if scribusDocument.Version != "1.5.1.svn" {
		t.Errorf("scribusDocument.Version was incorrect, got: %v, want: %v.", scribusDocument.Version, "1.5.1.svn")
	}

}
