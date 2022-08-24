package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/iancoleman/strcase"
)

type abstractTmplData struct {
	ProtoMessageName string
	ProtoFields      []protoField
	Enums            []enum
}

type protoField struct {
	FieldType string
	Name      string
	Number    int
}

type enum struct {
	Name    string
	Options []string
}

type recordField struct {
	name       string
	recordType string
	units      string
	startIdx   int
	endIdx     int
}

var colFieldName, colFieldType, colFieldUnits, colFieldOptionValues int

func generateTmplData(protoOut string, records [][]string) abstractTmplData {
	// Remove header
	header, records := records[0], records[1:]

	// Assign cols
	for idx, col := range header {
		switch col {
		case "Name":
			colFieldName = idx
		case "Type":
			colFieldType = idx
		case "Units":
			colFieldUnits = idx
		case "Option values":
			colFieldOptionValues = idx
		}
	}

	// Loop through records to figure out field start and end indices
	recordFields := getRecordFields(records)

	var protoFields []protoField
	var enums []enum

	// Loop through recordFields and generate data
	for idx, recordField := range recordFields {
		protoField := protoField{
			Name:   strcase.ToSnake(recordField.name),
			Number: idx + 1,
		}

		switch recordField.recordType {
		case "<string>":
			protoField.FieldType = "string"
		case "<double>":
			protoField.Name = fmt.Sprintf("%s_%s", protoField.Name, recordField.units)
			protoField.FieldType = "double"
		case "<select>":
			enum := getEnum(records, recordField)
			enums = append(enums, enum)
			protoField.FieldType = enum.Name
		case "<multiselect>":
			enum := getEnum(records, recordField)
			enums = append(enums, enum)
			protoField.FieldType = fmt.Sprintf("repeated %s", enum.Name)
		default:
			log.Fatal("Unrecognized recordType fro recordField")
		}

		protoFields = append(protoFields, protoField)
	}

	return abstractTmplData{
		ProtoMessageName: strcase.ToCamel((strings.Split(protoOut, ".")[0])),
		ProtoFields:      protoFields,
		Enums:            enums,
	}
}

func getRecordFields(records [][]string) []recordField {
	var recordFields []recordField

	for idx, row := range records {
		// New field
		if row[colFieldName] != "" {
			// Add endIndex to previous record field
			if len(recordFields) > 0 {
				recordFields[len(recordFields)-1].endIdx = idx - 1
			}

			recordFields = append(recordFields, recordField{
				name:       row[colFieldName],
				recordType: row[colFieldType],
				units:      row[colFieldUnits],
				startIdx:   idx,
			})

		}
	}

	// Add last endIndex
	recordFields[len(recordFields)-1].endIdx = len(records)

	return recordFields
}

func getEnum(records [][]string, recordField recordField) (enum enum) {
	enum.Name = strcase.ToCamel(recordField.name)

	for i := recordField.startIdx; i < recordField.endIdx; i++ {
		option := fmt.Sprintf("%s_%s",
			strcase.ToScreamingSnake(recordField.name),
			strcase.ToScreamingSnake(records[i][colFieldOptionValues]))
		enum.Options = append(enum.Options, option)
	}

	return
}
