package csvhandler

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	table1        [][]string
	table1Reorder [][]string
	table1Subset  [][]string
	table1Add     [][]string
	table2        [][]string
	table2Float   []interface{}
	table2Subset  []interface{}
	table2Reorder []interface{}
	table2Time    []interface{}
	oneCol        []interface{}
	oneColTime    []interface{}
	invalidTable  [][]string
)

func TestMain(m *testing.M) {
	table1 = [][]string{
		{"col1", "col2", "col3"},
		{"11", "12", "13"},
		{"21", "22", "23"},
		{"31", "32", "33"},
	}
	table1Reorder = [][]string{
		{"col2", "col1", "col3"},
		{"12", "11", "13"},
		{"22", "21", "23"},
		{"32", "31", "33"},
	}
	table1Subset = [][]string{
		{"col2", "col1"},
		{"12", "11"},
		{"22", "21"},
		{"32", "31"},
	}
	table1Add = [][]string{
		{"col1", "col2", "col3", "col4", "col5"},
		{"11", "12", "13", "", ""},
		{"21", "22", "23", "", ""},
		{"31", "32", "33", "", ""},
	}

	table2 = [][]string{
		{"col1", "col2", "col3"},
		{"11.1", "12", "13"},
		{"21.2", "22", "23"},
		{"31", "32", "33.3"},
	}
	table2Float = []interface{}{
		[]string{"col1", "col2", "col3"},
		[]float64{11.1, 12, 13},
		[]float64{21.2, 22, 23},
		[]float64{31, 32, 33.3},
	}
	table2Subset = []interface{}{
		[]string{"col1", "col3"},
		[]float64{11.1, 13},
		[]float64{21.2, 23},
		[]float64{31, 33.3},
	}
	table2Reorder = []interface{}{
		[]string{"col3", "col1"},
		[]float64{13, 11.1},
		[]float64{23, 21.2},
		[]float64{33.3, 31},
	}
	table2Time = []interface{}{
		[]string{"time", "col1", "col2", "col3"},
		[]float64{0, 11.1, 12, 13},
		[]float64{5, 21.2, 22, 23},
		[]float64{10, 31, 32, 33.3},
	}
	oneCol = []interface{}{
		[]string{"col1"},
		[]float64{11.1},
		[]float64{21.2},
		[]float64{31},
	}
	oneColTime = []interface{}{
		[]string{"time", "col1"},
		[]float64{0, 11.1},
		[]float64{5, 21.2},
		[]float64{10, 31},
	}
	invalidTable = [][]string{
		{"pre:col1", "pre:col2", "pre:col3"},
		{"11.1", "1a", "13"},
		{"21.2", "22", "23"},
		{"31", "32", "33.3f"},
	}
	code := m.Run()
	os.Exit(code)
}

func TestArrDifference(t *testing.T) {
	assert.Equal(t, []string{}, ArrDifference([]string{}, []string{}))
	assert.Equal(t, []string{}, ArrDifference([]string{}, []string{"1"}))
	assert.Equal(t, []string{"1"}, ArrDifference([]string{"1"}, []string{}))
	assert.Equal(t, []string{"3", "4"}, ArrDifference([]string{"3", "4"}, []string{"1", "2"}))
	assert.Equal(t, []string{}, ArrDifference([]string{"3", "2"}, []string{"1", "2", "3"}))
	assert.Equal(t, []string{"1", "4"}, ArrDifference([]string{"3", "2", "1", "4"}, []string{"2", "3"}))
	assert.Equal(t, []string{"5"}, ArrDifference([]string{"5", "1", "3"}, []string{"1", "2", "3"}))
}

func TestFindIdxInArr(t *testing.T) {
	assert.Equal(t, -1, findIdxInArr("1", []string{}))
	assert.Equal(t, 0, findIdxInArr("1", []string{"1"}))
	assert.Equal(t, 1, findIdxInArr("1", []string{"3", "1", "1", "2"}))
	assert.Equal(t, -1, findIdxInArr("5", []string{"3", "1", "1", "2"}))
	assert.Equal(t, 3, findIdxInArr("2", []string{"3", "1", "1", "2"}))
}

func TestArrUnionSorted(t *testing.T) {
	assert.Equal(t, []string{}, ArrUnionSorted([]string{}, []string{}))
	assert.Equal(t, []string{"1"}, ArrUnionSorted([]string{}, []string{"1"}))
	assert.Equal(t, []string{"1"}, ArrUnionSorted([]string{"1"}, []string{}))
	assert.Equal(t, []string{"1", "2", "3", "4"}, ArrUnionSorted([]string{"4", "3"}, []string{"1", "2"}))
	assert.Equal(t, []string{"1", "2", "3"}, ArrUnionSorted([]string{"3", "2"}, []string{"1", "2", "3"}))
	assert.Equal(t, []string{"1", "2", "3", "4"}, ArrUnionSorted([]string{"3", "2", "1", "4"}, []string{"2", "3"}))
	assert.Equal(t, []string{"1", "2", "3", "5"}, ArrUnionSorted([]string{"5", "1", "3"}, []string{"1", "2", "3"}))
}

func TestGenerateTableByTitles(t *testing.T) {
	newTable0, err := GenerateTableByTitles(nil, []string{"col2", "col1", "col3"})
	assert.NotNil(t, err)
	assert.Nil(t, newTable0)

	newTable1, err := GenerateTableByTitles(table1, []string{"col2", "col1", "col3"})
	assert.Nil(t, err)
	assert.Equal(t, table1Reorder, newTable1)

	newTable2, err := GenerateTableByTitles(table1, []string{"col2", "col1"})
	assert.Nil(t, err)
	assert.Equal(t, table1Subset, newTable2)

	newTable3, err := GenerateTableByTitles(table1, []string{"col1", "col2", "col3", "col4", "col5"})
	assert.Nil(t, err)
	assert.Equal(t, table1Add, newTable3)
}

func TestSubsetArrByIndicesAndConvertToFloat(t *testing.T) {
	res1, err := subsetArrByIndicesAndConvertToFloat([]string{}, nil)
	assert.Nil(t, err)
	assert.Equal(t, []float64{}, res1)
	res2, err := subsetArrByIndicesAndConvertToFloat(nil, []int{})
	assert.Nil(t, res2)
	assert.NotNil(t, err)
	res3, err := subsetArrByIndicesAndConvertToFloat([]string{}, []int{})
	assert.Nil(t, err)
	assert.Equal(t, []float64{}, res3)
	res4, err := subsetArrByIndicesAndConvertToFloat([]string{"13.1", "13.2", "14"}, []int{0, 1, 2})
	assert.Nil(t, err)
	assert.Equal(t, []float64{13.1, 13.2, 14}, res4)
	res5, err := subsetArrByIndicesAndConvertToFloat([]string{"13.1", "13.2", "14"}, []int{0, 1})
	assert.Nil(t, err)
	assert.Equal(t, []float64{13.1, 13.2}, res5)
	res6, err := subsetArrByIndicesAndConvertToFloat([]string{"13.1", "13.2", "14"}, []int{2, 0})
	assert.Nil(t, err)
	assert.Equal(t, []float64{14, 13.1}, res6)
	res7, err := subsetArrByIndicesAndConvertToFloat([]string{"13.1", "13.2", "14"}, []int{2, 3})
	assert.Nil(t, err)
	assert.Equal(t, []float64{14}, res7)
	res8, err := subsetArrByIndicesAndConvertToFloat([]string{"13.1", "abc", "14"}, []int{1, 2})
	assert.Nil(t, res8)
	assert.NotNil(t, err)
}

func TestGetColumnsInFloat(t *testing.T) {
	resNil, err := GetColumnsInFloat(nil, []string{})
	assert.NotNil(t, err)
	assert.Nil(t, resNil)

	resEmpty, err := GetColumnsInFloat(table2, []string{})
	assert.Nil(t, err)
	assert.Equal(t, []interface{}{}, resEmpty)

	resNotExist, err := GetColumnsInFloat(table2, []string{"col4", "col5"})
	assert.Nil(t, err)
	assert.Equal(t, []interface{}{[]string{}, []float64{}, []float64{}, []float64{}}, resNotExist)

	resInvalid, err := GetColumnsInFloat(invalidTable, []string{"pre:col1", "pre:col2", "pre:col3"})
	assert.NotNil(t, err)
	assert.Nil(t, resInvalid)

	res, err := GetColumnsInFloat(table2, []string{"col1", "col2", "col3"})
	assert.Nil(t, err)
	assert.Equal(t, table2Float, res)

	resSubset, err := GetColumnsInFloat(table2, []string{"col1", "col3"})
	assert.Nil(t, err)
	assert.Equal(t, table2Subset, resSubset)

	resSubsetNotExist, err := GetColumnsInFloat(table2, []string{"col1", "col4", "col3"})
	assert.Nil(t, err)
	assert.Equal(t, table2Subset, resSubsetNotExist)

	resReorder, err := GetColumnsInFloat(table2, []string{"col3", "col1"})
	assert.Nil(t, err)
	assert.Equal(t, table2Reorder, resReorder)
}

/* func TestAppendTimeToDataArr(t *testing.T) {
	resNil, err := AppendTimeToDataArr(nil)
	assert.NotNil(t, err)
	assert.Nil(t, resNil)

	resEmpty, err := AppendTimeToDataArr([]interface{}{})
	assert.NotNil(t, err)
	assert.Nil(t, resEmpty)

	res, err := AppendTimeToDataArr(table2Float)
	assert.Nil(t, err)
	assert.Equal(t, table2Time, res)

	resOneCol, err := AppendTimeToDataArr(oneCol)
	assert.Nil(t, err)
	assert.Equal(t, oneColTime, resOneCol)
} */
