package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/iancoleman/strcase"
)

type abstractTmplData struct {
	ProtoMessageName string
}

type protoField struct {
	fieldType string
	name      string
	number    int
}

type enum struct {
	name    string
	options []string
}

type recordField struct {
	name       string
	recordType string
	startIdx   int
	endIdx     int
}

func generateTmplData(protoOut string, records [][]string) abstractTmplData {
	// Remove header
	records = records[1:]

	// Loop through records to figure out field start and end indices
	recordFields := getRecordFields(records)
	log.Println(records)
	log.Println(recordFields)

	var protoFields []protoField
	var enums []enum

	// Loop through recordFields and generate data
	for idx, recordField := range recordFields {
		protoField := protoField{
			name:   strcase.ToSnake(recordField.name),
			number: idx + 1,
		}

		switch recordField.recordType {
		case "<select>", "<multiselect>":
			enumName := strcase.ToCamel(recordField.name)
			enumOptions := []string{}

			for i := recordField.startIdx; i < recordField.endIdx; i++ {
				option := fmt.Sprintf("%s_%s",
					strcase.ToScreamingSnake(recordField.name),
					strcase.ToScreamingSnake(records[i][2]))
				enumOptions = append(enumOptions, option)
			}

			enums = append(enums, enum{name: enumName, options: enumOptions})
			protoField.fieldType = enumName
		default:
			log.Panic("Unrecognized recordType fro recordField")
		}

		protoFields = append(protoFields, protoField)
	}

	log.Println(protoFields)
	log.Println(enums)

	return abstractTmplData{
		ProtoMessageName: strcase.ToCamel((strings.Split(protoOut, ".")[0])),
	}
}

func getRecordFields(records [][]string) []recordField {
	var recordFields []recordField

	for idx, row := range records {
		// New field
		if row[0] != "" {
			// Add endIndex to previous record field
			if len(recordFields) > 0 {
				recordFields[len(recordFields)-1].endIdx = idx - 1
			}

			recordFields = append(recordFields, recordField{
				name:       row[0],
				recordType: row[1],
				startIdx:   idx,
			})

		}
	}

	// Add last endIndex
	recordFields[len(recordFields)-1].endIdx = len(records)

	return recordFields
}
