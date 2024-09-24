package utils

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"testing"
	"text/template"

	"github.com/dongri/phonenumber"
)

// Contains find is string in slices
func Contains(slices []string, comparizon string) bool {
	for _, a := range slices {
		if a == comparizon {
			return true
		}
	}

	return false
}

func Bint(biSet bool) int {
	bitSet := true
	bitSetVar := 0
	if bitSet {
		bitSetVar = 1
	}

	return bitSetVar
}

// IsEmail check email the real email string
func IsEmail(email string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	return re.MatchString(email)
}

func IsUsernameTag(name string) bool {
	re := regexp.MustCompile(`^\@[a-zA-Z0-9]+`)

	return re.MatchString(name)
}

func IsCodeTag(code string) bool {
	re := regexp.MustCompile(`([a-zA-Z0-9,-]\-?[a-zA-Z0-9,-]\-?)`)

	return re.MatchString(code)
}

// IsNumeric ...
func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func IsPhone(s string) (phonenumber.ISO3166, string, bool) {
	d := phonenumber.GetISO3166ByNumber(s, true)
	if d.CountryName == "" {
		return d, "", false
	}

	number := phonenumber.ParseWithLandLine(s, d.Alpha3)
	if number == "" {
		return d, "", false
	}

	return d, number, true
}

// ToInt ...
func ToInt(s string) (int, error) {
	if len(s) == 0 {
		return 0, nil
	}

	if !IsNumeric(s) {
		return 0, fmt.Errorf("%s", "must be a number")
	}

	return strconv.Atoi(s)
}

func ToJSON(v interface{}) ([]byte, error) {
	str, err := json.Marshal(v)
	return str, err
}

func EncodeBase64(v interface{}) (string, error) {
	str, err := ToJSON(v)
	if err != nil {
		return "", err
	}

	return b64.StdEncoding.EncodeToString(str), nil
}

func DecodeBase64(v string) string {
	uDec, _ := b64.StdEncoding.DecodeString(v)
	return string(uDec)
}

func ParseTpl(fileName string, data interface{}) (string, error) {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return "", err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	log.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	log.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	log.Printf("\tSys = %v MiB", bToMb(m.Sys))
	log.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func LoadFixture(t *testing.T, n string) string {
	fixtureDir := "./test-fixtures"
	p := filepath.Join(fixtureDir, n)
	b, err := ioutil.ReadFile(p)
	if err != nil {
		t.Fatalf("Error while trying to read %s: %v\n", n, err)
	}

	return string(b)
}

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
