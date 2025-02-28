package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

// ParseCSVToStruct reads CSV data from a reader and converts it to a slice of structs
// The target parameter should be a pointer to a slice of structs
// The fields parameter maps CSV column names to struct field names
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
			if !field.IsValid() || !field.CanSet() {
				continue
			}

			// Get the value from the CSV row and trim spaces
			value := strings.TrimSpace(row[i])

			// Set the field value based on its type
			switch field.Kind() {
			case reflect.String:
				field.SetString(value)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if value == "" {
					field.SetInt(0)
				} else {
					intValue, err := strconv.ParseInt(value, 10, 64)
					if err != nil {
						// If conversion fails, set to zero
						field.SetInt(0)
					} else {
						field.SetInt(intValue)
					}
				}
			case reflect.Float32, reflect.Float64:
				if value == "" {
					field.SetFloat(0)
				} else {
					floatValue, err := strconv.ParseFloat(value, 64)
					if err != nil {
						// If conversion fails, set to zero
						field.SetFloat(0)
					} else {
						field.SetFloat(floatValue)
					}
				}
			case reflect.Bool:
				boolValue := false
				if strings.ToLower(value) == "true" || value == "1" {
					boolValue = true
				}
				field.SetBool(boolValue)
			case reflect.Slice:
				// Handle slice types (like []string) by splitting on comma
				if field.Type().Elem().Kind() == reflect.String {
					values := strings.Split(value, ",")
					slice := reflect.MakeSlice(field.Type(), len(values), len(values))
					for j, v := range values {
						slice.Index(j).SetString(strings.TrimSpace(v))
					}
					field.Set(slice)
				}
			}
			// For other types, we skip setting the value
		}

		results = reflect.Append(results, newStruct)
	}

	targetValue.Elem().Set(results)
	return nil
}
