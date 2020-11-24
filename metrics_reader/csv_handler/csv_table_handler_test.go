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
	table3        [][]string
	table3Subset1 [][]string
	table3Subset2 [][]string
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
	table3 = [][]string{
		{"time", "col1"},
		{"1", "11.1"},
		{"2", "21.2"},
		{"3", "31"},
	}
	table3Subset1 = [][]string{
		{"time", "col1"},
		{"2", "21.2"},
		{"3", "31"},
	}
	table3Subset2 = [][]string{
		{"time", "col1"},
	}
	invalidTable = [][]string{
		{"time", "pre:col2", "pre:col3"},
		{"11", "1a", "13.1"},
		{"21", "22.2", "23"},
		{"31bb", "32", "33.3f"},
	}
	code := m.Run()
	os.Exit(code)
}

func TestFindIdxInArr(t *testing.T) {
	assert.Equal(t, -1, findIdxInArr("1", []string{}))
	assert.Equal(t, 0, findIdxInArr("1", []string{"1"}))
	assert.Equal(t, 1, findIdxInArr("1", []string{"3", "1", "1", "2"}))
	assert.Equal(t, -1, findIdxInArr("5", []string{"3", "1", "1", "2"}))
	assert.Equal(t, 3, findIdxInArr("2", []string{"3", "1", "1", "2"}))
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
	assert.Equal(t, []*float64{}, res1)
	res2, err := subsetArrByIndicesAndConvertToFloat(nil, []int{})
	assert.Nil(t, res2)
	assert.NotNil(t, err)
	res3, err := subsetArrByIndicesAndConvertToFloat([]string{}, []int{})
	assert.Nil(t, err)
	assert.Equal(t, []*float64{}, res3)
	res4, err := subsetArrByIndicesAndConvertToFloat([]string{"13.1", "13.2", "14"}, []int{0, 1, 2})
	assert.Nil(t, err)
	//assert.Equal(t, []float64{13.1, 13.2, 14}, res4)
	assert.Equal(t, 3, len(res4))
	assert.Equal(t, 13.1, *(res4[0]))
	assert.Equal(t, 13.2, *(res4[1]))
	assert.Equal(t, float64(14), *(res4[2]))
	res5, err := subsetArrByIndicesAndConvertToFloat([]string{"13.1", "13.2", "14"}, []int{0, 1})
	assert.Nil(t, err)
	//assert.Equal(t, []float64{13.1, 13.2}, res5)
	assert.Equal(t, 2, len(res5))
	assert.Equal(t, 13.1, *(res5[0]))
	assert.Equal(t, 13.2, *(res5[1]))
	res6, err := subsetArrByIndicesAndConvertToFloat([]string{"13.1", "13.2", "14"}, []int{2, 0})
	assert.Nil(t, err)
	//assert.Equal(t, []float64{14, 13.1}, res6)
	assert.Equal(t, 2, len(res6))
	assert.Equal(t, float64(14), *(res6[0]))
	assert.Equal(t, 13.1, *(res6[1]))
	res7, err := subsetArrByIndicesAndConvertToFloat([]string{"13.1", "13.2", "14"}, []int{2, 3})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(res7))
	assert.Equal(t, float64(14), *(res7[0]))
	res8, err := subsetArrByIndicesAndConvertToFloat([]string{"13.1", "abc", "14"}, []int{1, 2})
	assert.Nil(t, res8)
	assert.NotNil(t, err)
	res9, err := subsetArrByIndicesAndConvertToFloat([]string{"", ""}, []int{0})
	assert.Nil(t, err)
	assert.Equal(t, []*float64{(*float64)(nil)}, res9)
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
	assert.Equal(t, []interface{}{[]string{}, []*float64{}, []*float64{}, []*float64{}}, resNotExist)

	resInvalid, err := GetColumnsInFloat(invalidTable, []string{"pre:col1", "pre:col2", "pre:col3"})
	assert.NotNil(t, err)
	assert.Nil(t, resInvalid)

	res, err := GetColumnsInFloat(table2, []string{"col1", "col2", "col3"})
	assert.Nil(t, err)
	//assert.Equal(t, table2Float, res)
	compareTables(t, table2Float, res)

	resSubset, err := GetColumnsInFloat(table2, []string{"col1", "col3"})
	assert.Nil(t, err)
	//assert.Equal(t, table2Subset, resSubset)
	compareTables(t, table2Subset, resSubset)

	resSubsetNotExist, err := GetColumnsInFloat(table2, []string{"col1", "col4", "col3"})
	assert.Nil(t, err)
	//assert.Equal(t, table2Subset, resSubsetNotExist)
	compareTables(t, table2Subset, resSubsetNotExist)

	resReorder, err := GetColumnsInFloat(table2, []string{"col3", "col1"})
	assert.Nil(t, err)
	//assert.Equal(t, table2Reorder, resReorder)
	compareTables(t, table2Reorder, resReorder)
}

func compareTables(t *testing.T, expected []interface{}, actual []interface{}) {
	assert.Equal(t, len(expected), len(actual))
	assert.True(t, len(actual) > 1)
	assert.Equal(t, expected[0], actual[0])

	for i := 1; i < len(expected); i++ {
		var rowExpected []float64 = expected[i].([]float64)
		var rowActual []*float64 = actual[i].([]*float64)
		for j, v := range rowActual {
			assert.Equal(t, rowExpected[j], *v)
		}
	}
}

func TestSubsetDataArrByTime(t *testing.T) {
	subset1, err := SubsetDataArrByTime(table3, 2)
	assert.Nil(t, err)
	assert.Equal(t, table3Subset1, subset1)

	subset2, err := SubsetDataArrByTime(table3, 5)
	assert.Nil(t, err)
	assert.Equal(t, table3Subset2, subset2)

	nilData, err := SubsetDataArrByTime(nil, 2)
	assert.Nil(t, nilData)
	assert.NotNil(t, err)

	insuffData, err := SubsetDataArrByTime(table3Subset2, 2)
	assert.Nil(t, insuffData)
	assert.NotNil(t, err)

	invalidData, err := SubsetDataArrByTime(invalidTable, 2)
	assert.Nil(t, invalidData)
	assert.NotNil(t, err)
}
