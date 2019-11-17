package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/probonopd/go-scribus"
	"reflect"
)

func main() {
	// Load Scribus file
	doc, err := scribus.NewScribusDocumentFromFile("Document-1.sla")
	if err != nil {
		fmt.Println("Error", err)
	}

	// Pretty-print the ScribusDocument struct
	spew.Dump(doc)

	// Find all PAGEOBJECT elements
	zzzGetPageObjects(doc.DOCUMENT)
}

// FIXME: This does not work because
// cannot range over document.DOCUMENT (type scribus.DOCUMENT)
func GetPageObjects(document scribus.ScribusDocument) []*scribus.PAGEOBJECT {
	var pageobjects []*scribus.PAGEOBJECT
	// for elem := range document.DOCUMENT {
	// 	fmt.Println(elem)
	// }
	return pageobjects
}

// FIXME: This finds PAGEOBJECTs but how can I return references to them?
func zzzGetPageObjects(obj interface{}) []*scribus.PAGEOBJECT {
	var pageobjects []*scribus.PAGEOBJECT
	s := reflect.Indirect(reflect.ValueOf(obj))
	for i := 0; i < s.NumField(); i++ {
		f := s.Type().Field(i)
		field := s.Field(i)
		if fieldIsExported(f) { // Exported-check must be evaluated first to avoid panic.
			fmt.Println(f.Type, f.Name)
			if f.Name == "PAGEOBJECT" {
				fmt.Println("FOUND PAGEOBJECT *** but how can I get a pointer to that struct? ***")
				//pageobjects = append(pageobjects, obj.(*scribus.PAGEOBJECT)) // How can I get the pointer to the PAGEOBJECT?
			}
			if field.Kind() == reflect.Struct {
				// fmt.Println("recursing")
				zzzGetPageObjects(field)
			}
		}
	}
	return pageobjects
}

func fieldIsExported(field reflect.StructField) bool {
	// log.Println(field.Name)
	return field.Name[0] >= 65 == true && field.Name[0] <= 90 == true
}
