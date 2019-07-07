package main

import (
	"compress/bzip2"
	"encoding/xml"
	//	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	file, err := os.Open("file.bz")
	if err != nil {
		log.Fatal(err)
	}
	out, err := os.Create("wikipagelinks")
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(out, "", 0)
	r := bzip2.NewReader(file)
	d := xml.NewDecoder(r)
	re, _ := regexp.Compile("\\[\\[([^|\\]]*)(?:\\|[^\\]]+)?\\]\\]")
	counter := 0
	for {
		t, err := d.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		switch t := t.(type) {
		case xml.StartElement:
			if t.Name.Local == "title" {
				var text string
				if err := d.DecodeElement(&text, &t); err != nil {
					log.Fatal(err)
				}
				//fmt.Println(text)
				//out.WriteString(text)
				//out.WriteString("\n")
				logger.Output(2, text)
			} else if t.Name.Local == "text" {
				var text string
				if err := d.DecodeElement(&text, &t); err != nil {
					log.Fatal(err)
				}
				articles := make(map[string]bool)
				for _, match := range re.FindAllStringSubmatch(text, -1) {
					if strings.HasPrefix(match[1], "Category:") {
						continue
					}
					if _, exists := articles[match[1]]; !exists {
						//fmt.Println(match[1])
						//out.WriteString("\t")
						//out.WriteString(match[1])
						//out.WriteString("\n")
						logger.Output(2, "\t"+match[1])
						articles[match[1]] = true
					}
				}
				if counter == 2000 {
					return
				}
				counter++
			}
		}
	}
}
