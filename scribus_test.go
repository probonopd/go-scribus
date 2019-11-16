package scribus

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
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

func TestEditScribusDocument(t *testing.T) {

	scribusDocument, err := NewScribusDocumentFromFile("Document-1.sla")
	if err != nil {
		t.Errorf("error: %v", err)
	}

	// Print all page objects, copy them, and change their text
	// See https://stackoverflow.com/a/28041994 for why to do it with "i"
	for i, _ := range scribusDocument.DOCUMENT.PAGEOBJECT {
		t.Logf("PAGEOBJECT %v", scribusDocument.DOCUMENT.PAGEOBJECT[i].StoryText.ITEXT[0].CH)
		scribusDocument.ChangeTextOfPageObject(i, "Changed")                      // Change its text
		scribusDocument.DuplicatePageObject(i)                                    // Duplicate the object
		oldxpos, _ := strconv.Atoi(scribusDocument.DOCUMENT.PAGEOBJECT[i+1].XPOS) // Get the position of the object on the page
		oldypos, _ := strconv.Atoi(scribusDocument.DOCUMENT.PAGEOBJECT[i+1].YPOS) // Get the position of the object on the page
		if oldxpos != 140 {
			t.Errorf("oldxpos was incorrect, got: %d, want: %d.", oldxpos, 140)
		}
		scribusDocument.MovePageObject(i+1, oldxpos*2, oldypos)     // Change the position of the copy on the page so that we can see it
		scribusDocument.ChangeTextOfPageObject(i+1, "Changed Copy") // Change the text of the copy
	}

	// Write out changed file
	scribusDocument.WriteScribusFile("Changed-Document-HighLevel-1.sla")
	if err != nil {
		t.Errorf("error: %v", err)
	}
}

func TestEditScribusDocumentLowLevel(t *testing.T) {

	// Read in test Scribus file
	xmlFile, err := os.Open("Document-1.sla")
	if err != nil {
		t.Errorf("error: %v", err)
	}

	defer xmlFile.Close()
	byteValue, _ := ioutil.ReadAll(xmlFile)
	var scribusDocument ScribusDocument
	err = xml.Unmarshal(byteValue, &scribusDocument)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	// Check if we got the correct version back
	if scribusDocument.Version != "1.5.1.svn" {
		t.Errorf("scribusDocument.Version was incorrect, got: %v, want: %v.", scribusDocument.Version, "1.5.1.svn")
	}

	// Print all page objects, copy them, and change their text
	// See https://stackoverflow.com/a/28041994 for why to do it with "i"
	for i := range scribusDocument.DOCUMENT.PAGEOBJECT {
		t.Logf("PAGEOBJECT %v", scribusDocument.DOCUMENT.PAGEOBJECT[i].StoryText.ITEXT[0].CH)
		// Change its text
		scribusDocument.DOCUMENT.PAGEOBJECT[i].StoryText.ITEXT[0].CH = "Changed"
		// Insert a copied PAGEOBJECT, https://stackoverflow.com/a/51311878
		// Increase the slice scribusDocument.DOCUMENT.PAGEOBJECT by 1 (by copying a PAGEOBJECT)
		scribusDocument.DOCUMENT.PAGEOBJECT = append(scribusDocument.DOCUMENT.PAGEOBJECT, scribusDocument.DOCUMENT.PAGEOBJECT[i])
		// Move the objects thereafter, to make space for the object we want to insert
		copy(scribusDocument.DOCUMENT.PAGEOBJECT[i+1:], scribusDocument.DOCUMENT.PAGEOBJECT[i:])
		// Insert the object
		scribusDocument.DOCUMENT.PAGEOBJECT[i+1] = scribusDocument.DOCUMENT.PAGEOBJECT[i]
		// Change its position on the page so that we can see it
		oldxpos, _ := strconv.Atoi(scribusDocument.DOCUMENT.PAGEOBJECT[i+1].XPOS)

		if oldxpos != 140 {
			t.Errorf("oldxpos was incorrect, got: %d, want: %d.", oldxpos, 140)
		}
		scribusDocument.DOCUMENT.PAGEOBJECT[i+1].XPOS = strconv.Itoa(oldxpos * 2)
		// Change its text
		scribusDocument.DOCUMENT.PAGEOBJECT[i+1].StoryText.ITEXT[0].CH = "Changed Copy"
	}

	// Write out changed file
	if xmlstring, err := xml.MarshalIndent(scribusDocument, "", "    "); err == nil {
		xmlstring = []byte(xml.Header + strings.Replace(string(xmlstring), "&#xA;", "", -1)) // FIXME: https://forum.golangbridge.org/t/read-xml-change-values-write-back-crippled-file/16253/4
		_ = ioutil.WriteFile("Changed-Document-1.sla", xmlstring, 0644)
	}

}
