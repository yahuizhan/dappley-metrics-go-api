package csvhandler

import (
	"encoding/csv"
	"os"

	logger "github.com/sirupsen/logrus"
)

// ReadCSV reads the entire csv file at filepath
func ReadCSV(filepath string) ([][]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	return csvReader.ReadAll()
}

// IsCSVExistAndNonEmpty checks if file at filepath exists and not empty
func IsCSVExistAndNonEmpty(filepath string) bool {
	stat, err := os.Stat(filepath)
	if os.IsNotExist(err) || stat.Size() == 0 {
		return false
	}
	return true
}

// CreateNewCSVWithTitles creates new csv file at filepath and writes header to it
func CreateNewCSVWithTitles(filepath string, colTitles []string) *os.File {
	file, err := os.Create(filepath)
	if err != nil {
		logger.Panic("cannot create a csv file for Error: " + err.Error())
	}
	writer := csv.NewWriter(file)
	writer.Comma = ','
	writer.Write(colTitles)
	writer.Flush()
	return file
}

/* func CreateNewCSVWithTable(filepath string, newTable [][]string) *os.File {
	file, err := os.Create(filepath)
	if err != nil {
		logger.Panic("cannot create a csv file for Error: " + err.Error())
	}
	writer := csv.NewWriter(file)
	writer.Comma = ','
	writer.WriteAll(newTable)
	return file
} */

/* func GetCSVHeader(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	return csvReader.Read()
} */

/* func AddNewColumnsToCSV(filepath string, colsToAdd []string) (file *os.File, err error) {
	//
} */
