package csvhandler

import (
	"errors"
	"strconv"
	"strings"
)

var (
	errNotEnoughDataInCSV error = errors.New("CSV File Only Contains One Rows; Not Enought For Plotting")
	errArrNil             error = errors.New("subsetArrByIndicesAndConvertToFloat: cannot subset nil array")
)

// ArrDifference computes A - B
/* func ArrDifference(a, b []string) []string {
	diff := []string{}

	m := make(map[string]bool)
	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if !(m[item]) {
			diff = append(diff, item)
		}
	}
	return diff
} */

// ArrUnionSorted computes A + B
/* func ArrUnionSorted(a, b []string) (union []string) {
	diff := ArrDifference(a, b)
	union = append(b, diff...)
	sort.Strings(union)
	return union
} */

// GenerateTableByTitles generates a new table based on columnTitles and data in oldTable
// if columnTitles is a subset of oldTable's titles, return the subset of oldTable based on columnTitles
// if columnTitles contains extra columns not found in oldTable, add the columns and assign "" as values
func GenerateTableByTitles(oldTable [][]string, columnTitles []string) (newTable [][]string, err error) {
	if oldTable == nil || len(oldTable) < 2 {
		return nil, errNotEnoughDataInCSV
	}
	newTable = append(newTable, columnTitles)

	oldTitles := oldTable[0]
	oldTitleIndics := make(map[string]int)
	for _, col := range columnTitles {
		idx := findIdxInArr(col, oldTitles)
		oldTitleIndics[col] = idx
	}

	for i := 1; i < len(oldTable); i++ {
		row := oldTable[i]
		var newRow []string
		for _, col := range columnTitles {
			if oldTitleIndics[col] < 0 {
				newRow = append(newRow, "")
			} else {
				newRow = append(newRow, row[oldTitleIndics[col]])
			}
		}
		newTable = append(newTable, newRow)
	}
	return
}

func findIdxInArr(elem string, arr []string) int {
	for k, v := range arr {
		if v == elem {
			return k
		}
	}

	return -1
}

// GetColumnsInFloat subsets csvTable by columnTitles and convert all data values to *float64
// if some element in columnTitles does not exist in csvTable, it is ignored
func GetColumnsInFloat(csvTable [][]string, columnTitles []string) ([]interface{}, error) {
	// subset csvTable to keep only columnTitles and convert data values to float64
	if csvTable == nil || len(csvTable) < 2 {
		return nil, errNotEnoughDataInCSV
	}

	res := []interface{}{}
	if columnTitles == nil || len(columnTitles) == 0 {
		return res, nil
	}

	columnsIdx := []int{}
	newTitleLine := []string{}
	titleLine := csvTable[0]
	for _, col := range columnTitles {
		idx := findIdxInArr(col, titleLine)
		if idx >= 0 {
			columnsIdx = append(columnsIdx, idx)
			prefixRemovedCol := removeFirstPrefix(col, ":")
			newTitleLine = append(newTitleLine, prefixRemovedCol)
		}
	}
	res = append(res, newTitleLine)
	for i := 1; i < len(csvTable); i++ {
		row := csvTable[i]
		newRow, err := subsetArrByIndicesAndConvertToFloat(row, columnsIdx)
		if err != nil {
			return nil, err
		}
		res = append(res, newRow)
	}
	return res, nil
}

func removeFirstPrefix(title string, delim string) string {
	secs := strings.Split(title, delim)
	if len(secs) < 2 { // title does not contain delim
		return title
	}
	return strings.Join(secs[1:], delim)
}

// subset 'arr' by indices given in 'indices'; if index in 'indices' is out of range of 'arr', it gets ignored
func subsetArrByIndicesAndConvertToFloat(arr []string, indices []int) ([]*float64, error) {
	if arr == nil {
		return nil, errArrNil
	}
	res := []*float64{}
	for _, i := range indices {
		if i >= 0 && i < len(arr) {
			if strings.TrimSpace(arr[i]) == "" {
				res = append(res, nil)
			} else {
				converted, err := strconv.ParseFloat(arr[i], 64)
				if err != nil {
					return nil, err
				}
				res = append(res, &converted)
			}
		}
	}
	return res, nil
}

// SubsetDataArrByTime keeps only rows that have time >= fromtime
func SubsetDataArrByTime(csvArr [][]string, fromtime int) ([][]string, error) {
	// assume first row is header and first column is time
	if csvArr == nil || len(csvArr) < 2 {
		return nil, errNotEnoughDataInCSV
	}
	res := [][]string{csvArr[0]}
	for i := 1; i < len(csvArr); i++ {
		row := csvArr[i]
		time, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, err
		}
		if time >= fromtime {
			res = append(res, row)
		}
	}
	return res, nil
}
