package scribus

import (
	"fmt"
	"os"

	"testing"
)

func TestReadScribusFile(t *testing.T) {
	document, err := NewScribusDocumentFromFile("Document-1.sla")
	if err != nil {
		t.Errorf("error: %v", err)
	}

	// Check if we got the correct version back
	if document.Version != "1.5.1.svn" {
		t.Errorf("document.Version was incorrect, got: %v, want: %v.", document.Version, "1.5.1.svn")
	}

}

func TestWriteScribusFile(t *testing.T) {
	document, err := NewScribusDocumentFromFile("Document-1.sla")
	if err != nil {
		t.Errorf("error: %v", err)
	}

	// Make a simple text change
	for _, po := range document.DOCUMENT.GetPageObjectsWithText("One") {
		po.StoryText.ITEXT[0].CH = "Changed"
		fmt.Println(po.StoryText.ITEXT[0].CH)
	}

	// Change some bullet points
	for _, po := range document.DOCUMENT.GetPageObjectsWithText("Three") {
		po.StoryText.ChangeBulletPoints([]string{"AAA", "BBB", "CCC"})
	}

	// Change a picture
	for i := range document.DOCUMENT.PAGEOBJECT { // Need to use 'i' so that we can edit the original rather than a copy
		if document.DOCUMENT.PAGEOBJECT[i].PTYPE == "2" { // Assuming that 2 means picture
			document.DOCUMENT.PAGEOBJECT[i].PFILE = "/usr/share/icons/hicolor/16x16/apps/software-properties.png"
		}
	}

	// Write back Scribus file
	testfile := "test.sla"
	err = document.WriteScribusFile(testfile)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	// Test if the written file exists
	if _, err := os.Stat(testfile); os.IsNotExist(err) {
		t.Errorf("error: %v does not exist", testfile)
	}

	// Test if we can read the written file
	document, err = NewScribusDocumentFromFile(testfile)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	// Check if we got the correct version back
	if document.Version != "1.5.1.svn" {
		t.Errorf("document.Version was incorrect, got: %v, want: %v.", document.Version, "1.5.1.svn")
	}

}
