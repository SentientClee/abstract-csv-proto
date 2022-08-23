package main

import (
	"embed"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"text/template"
)

//go:embed abstract.proto.tmpl
var templateFile embed.FS

func main() {
	var protoOut, csvFileName string
	flag.StringVar(&protoOut, "o", "abstract.proto", "The output filename.")
	flag.StringVar(&csvFileName, "f", "", "The input filename.")
	flag.Parse()

	if csvFileName == "" {
		log.Panic("No csv file provided.")
	}

	file, err := os.Open(csvFileName)
	if err != nil {
		log.Panicf("Unable to open csv %s\n", csvFileName)
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Panicf("Error reading records %s", err)
	}

	t := template.Must(template.ParseFS(templateFile, "abstract.proto.tmpl"))
	if err != nil {
		log.Panicf("Error parsing template %s", err)
	}

	p, err := os.Create(fmt.Sprintf("./%s", protoOut))
	if err != nil {
		log.Panic(err)
	}
	t.Execute(p, generateTmplData(protoOut, records))
}
