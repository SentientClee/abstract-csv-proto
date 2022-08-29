package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/iancoleman/strcase"
)

type pathAbstract struct {
	ProtoMessageName string
	BiopsyFields     []protoField
	ResectionFields  []protoField
	Enums            []enum
}

type protoField struct {
	FieldType string
	Name      string
}

type enum struct {
	Name    string
	Options []string
}

type recordField struct {
	name          string
	recordType    string
	units         string
	isBiopsyField bool
	startIdx      int
	endIdx        int
}

var colFieldName, colFieldType, colFieldUnits, colFieldOptionValues, colBiopsyField int

func genPathAbstract(protoOut string, records [][]string) pathAbstract {
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
		case "Biopsy fields":
			colBiopsyField = idx
		}
	}

	// Loop through records to figure out field start and end indices
	recordFields := getRecordFields(records)

	var resectionFields []protoField
	var biopsyFields []protoField
	var enums []enum

	// Loop through recordFields and generate data
	for _, recordField := range recordFields {
		protoField := protoField{
			Name: strcase.ToSnake(recordField.name),
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
			log.Fatalf("Unrecognized recordType %s for recordField %s", recordField.recordType, recordField.name)
		}

		if recordField.isBiopsyField {
			biopsyFields = append(biopsyFields, protoField)
		}
		resectionFields = append(resectionFields, protoField)
	}

	return pathAbstract{
		ProtoMessageName: strcase.ToCamel((strings.Split(protoOut, ".")[0])),
		BiopsyFields:     biopsyFields,
		ResectionFields:  resectionFields,
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
				name:          row[colFieldName],
				recordType:    row[colFieldType],
				units:         row[colFieldUnits],
				isBiopsyField: row[colBiopsyField] == "TRUE",
				startIdx:      idx,
			})

		}
	}

	// Add last endIndex
	recordFields[len(recordFields)-1].endIdx = len(records)

	return recordFields
}

func getEnum(records [][]string, recordField recordField) (enum enum) {
	lowerRecordName := strings.ToLower(recordField.name)
	enum.Name = strcase.ToCamel(lowerRecordName)

	for i := recordField.startIdx; i < recordField.endIdx; i++ {
		option := fmt.Sprintf("%s_%s",
			strcase.ToScreamingSnake(lowerRecordName),
			strcase.ToScreamingSnake(records[i][colFieldOptionValues]))
		enum.Options = append(enum.Options, option)
	}

	return
}
