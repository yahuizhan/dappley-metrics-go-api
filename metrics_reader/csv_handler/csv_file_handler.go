package csvhandler

import (
	"encoding/csv"
	"os"

	logger "github.com/sirupsen/logrus"
)

func ReadCSV(filepath string) ([][]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	return csvReader.ReadAll()
}

func IsCSVExistAndNonEmpty(filepath string) bool {
	stat, err := os.Stat(filepath)
	if os.IsNotExist(err) || stat.Size() == 0 {
		return false
	}
	return true
}

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

func CreateNewCSVWithTable(filepath string, newTable [][]string) *os.File {
	file, err := os.Create(filepath)
	if err != nil {
		logger.Panic("cannot create a csv file for Error: " + err.Error())
	}
	writer := csv.NewWriter(file)
	writer.Comma = ','
	writer.WriteAll(newTable)
	return file
}

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
