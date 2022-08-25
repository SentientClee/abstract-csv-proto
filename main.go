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

//go:embed path.proto.tmpl
var templateFile embed.FS

func main() {
	var protoOut, csvFileName string
	flag.StringVar(&protoOut, "o", "abstract.proto", "The output filename.")
	flag.StringVar(&csvFileName, "f", "", "The input filename.")
	flag.Parse()

	if csvFileName == "" {
		log.Fatal("No csv file provided.")
	}

	file, err := os.Open(csvFileName)
	if err != nil {
		log.Fatalf("Unable to open csv %s\n", csvFileName)
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error reading records %s", err)
	}

	t, err := template.New("path.proto.tmpl").
		Funcs(template.FuncMap{"inc": inc}).
		ParseFS(templateFile, "path.proto.tmpl")
	if err != nil {
		log.Fatalf("Error parsing template %s", err)
	}

	p, err := os.Create(fmt.Sprintf("./%s", protoOut))
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(p, genPathAbstract(protoOut, records))
}

func inc(num int) int {
	return num + 1
}
