package scribus

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestEditScribusDocumentLowLevel(t *testing.T) {

	// Read in test Scribus file
	xmlFile, err := os.Open("Document-1.sla")
	if err != nil {
		t.Errorf("error: %v", err)
	}

	defer xmlFile.Close()
	byteValue, _ := ioutil.ReadAll(xmlFile)
	var scribusDocument SCRIBUSUTF8NEW
	err = xml.Unmarshal(byteValue, &scribusDocument)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	// Print information from the XML
	t.Logf("Scribus version %v", scribusDocument.Version)

	// Check if we got the correct version back
	if scribusDocument.Version != "1.5.1.svn" {
		t.Errorf("scribusDocument.Version was incorrect, got: %v, want: %v.", scribusDocument.Version, "1.5.1.svn")
	}

	// Print all page objects, copy them, and change their text
	// See https://stackoverflow.com/a/28041994 for why to do it with "i"
	for i := range scribusDocument.DOCUMENT.PAGEOBJECT {
		t.Logf("PAGEOBJECT %v", scribusDocument.DOCUMENT.PAGEOBJECT[i].StoryText.ITEXT.CH)
		// Change its text
		scribusDocument.DOCUMENT.PAGEOBJECT[i].StoryText.ITEXT.CH = "Changed"
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
			t.Errorf("Sum was incorrect, got: %d, want: %d.", oldxpos, 140)
		}
		scribusDocument.DOCUMENT.PAGEOBJECT[i+1].XPOS = strconv.Itoa(oldxpos * 2)
		// Change its text
		scribusDocument.DOCUMENT.PAGEOBJECT[i+1].StoryText.ITEXT.CH = "Changed Copy"
	}

	// Write out changed file
	if xmlstring, err := xml.MarshalIndent(scribusDocument, "", "    "); err == nil {
		xmlstring = []byte(xml.Header + strings.Replace(string(xmlstring), "&#xA;", "", -1)) // FIXME: https://forum.golangbridge.org/t/read-xml-change-values-write-back-crippled-file/16253/4
		_ = ioutil.WriteFile("Changed-Document-1.sla", xmlstring, 0644)
	}

}
