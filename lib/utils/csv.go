package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"strings"
)

func ParseCSVToStruct(reader io.Reader, target interface{}, fields map[string]string) error {
	targetValue := reflect.ValueOf(target)

	if targetValue.Kind() != reflect.Pointer || targetValue.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("Target must be a pointer to slice")
	}

	sliceElemType := targetValue.Elem().Type().Elem()

	csvReader := csv.NewReader(reader)

	headers, err := csvReader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV headers: %v", err)
	}

	for i := range headers {
		headers[i] = strings.TrimSpace(strings.ToLower(headers[i]))
	}

	results := reflect.MakeSlice(targetValue.Elem().Type(), 0, 0)
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read CSV row: %v", err)
		}

		newStruct := reflect.New(sliceElemType).Elem()
		for i, header := range headers {
			if i >= len(row) {
				continue
			}

			fieldName, exists := fields[header]
			if !exists {
				continue
			}

			field := newStruct.FieldByName(fieldName)
			if !field.IsValid() {
				continue
			}

			field.SetString(strings.TrimSpace(row[i]))
		}

		results = reflect.Append(results, newStruct)
	}

	targetValue.Elem().Set(results)
	return nil
}
