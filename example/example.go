package main

import (
	"fmt"
	"github.com/probonopd/go-scribus"
)

func main() {
	// Load Scribus file
	doc, err := scribus.NewScribusDocumentFromFile("Document-1.sla")
	if err != nil {
		fmt.Println("Error", err)
	}

	for i := range doc.DOCUMENT.PAGEOBJECT {
		fmt.Println(doc.DOCUMENT.PAGEOBJECT[i].StoryText.ITEXT.CH)
	}

	for _, po := range doc.DOCUMENT.GetPageObjectsWithText("One") {
		po.StoryText.ITEXT.CH = "Foo"
		fmt.Println(po.StoryText.ITEXT.CH)
	}

}
