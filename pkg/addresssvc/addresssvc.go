package addresssvc

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readCSV(filename string) ([][]string, error) {
	//Open CSV file
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	// Read File into a Variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}

	return lines[1:][:], nil
}

func readPostcodes(filename string) (map[string]struct{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	postcode := make(map[string]struct{})

	for scanner.Scan() {
		postcode[scanner.Text()] = struct{}{}
	}

	fmt.Println("loaded all postcodes")

	return postcode, nil

}

//func init() {
//	p, err := readPostcodes("pkg/addresssvc/data/london_postcodes_ver1.txt")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	postCodes = p
//}

func CheckPostCode(code string) bool {
	code = strings.TrimSpace(strings.Replace(strings.ToLower(code), " ", "", -1))

	_, ok := postCodes[code]
	return ok
}

func GetPostcodeCoordinates(code string) (string, string) {
	coordinates, ok := postCodesCoordinates[code]
	if !ok {
		return "", ""
	}
	splittedCoordinates := strings.SplitAfter(coordinates, ",")
	if len(splittedCoordinates) == 2 {
		return strings.ReplaceAll(splittedCoordinates[0], ",", ""), strings.ReplaceAll(splittedCoordinates[1], ",", "")
	}
	return "", ""

}

func CreateCustomCodes() {
	csv, err := readCSV("pkg/addresssvc/data/London_postcodes.csv")

	if err != nil {
		fmt.Println("something happened: ", err)
		return
	}

	var postcodeMap map[string]string

	postcodeMap = make(map[string]string, 0)
	for x, row := range csv {
		for y, _ := range row {
			londonZone, err := strconv.Atoi(row[27])
			if err != nil {
				fmt.Println(err)
				return
			}

			if y == 0 && row[1] == "Yes" && londonZone < 5 {
				postcode := csv[x][y]
				postcode = strings.Replace(strings.ToLower(postcode), " ", "", 1)
				postcodeMap[postcode] = row[2] + "," + row[3]

				//postcodeList = append(postcodeList, postcode)
			}
		}
	}

	f, err := os.Create("postcodes_with_latitude.go")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer f.Close()

	b := &bytes.Buffer{}

	for code, cordinates := range postcodeMap {
		b.WriteString("\"" + code + "\"" + ":" + "\"" + cordinates + "\"" + ",\n")
	}

	f.Write(b.Bytes())
}
