package metricsreader

import (
	"encoding/csv"
	"errors"
	"os"
	"strconv"
	"strings"
)

var (
	allSections              []string = []string{"block", "cpu", "disk", "memory", "network", "txRequest"}
	blockConstantsFields     []string = []string{}
	cpuConstantsFields       []string = []string{"totalCoreNum"}
	diskConstantsFields      []string = []string{}
	memoryConstantsFields    []string = []string{"systemMem"}
	networkConstantsFields   []string = []string{}
	txRequestConstantsFields []string = []string{}
)

var (
	ErrNotEnoughDataInCSV error = errors.New("CSV File Only Contains One Rows; Not Enought For Plotting")
)

// MetricsInfoResponse forms JSON structure used at the frontend to generate plots
type MetricsInfoResponse struct {
	Success    bool                           `json:"success"`
	ErrMessage string                         `json:"error"`
	Data       map[string]*MetricsInfoSection `json:"data"`
}

// FormMetricsInfoResponse initializes and returns a MetricsInfoResponse
func FormMetricsInfoResponse(filepath string, limit int) *MetricsInfoResponse {
	response := &MetricsInfoResponse{
		Success:    false,
		ErrMessage: "",
		Data:       make(map[string]*MetricsInfoSection),
	}

	if allRecords, err := readCSV(filepath); err != nil {
		response.Data = nil
		response.ErrMessage = "Unable to parse file as CSV for Error: " + err.Error()
		return response
	}
	//logger.Infoln(allRecords)

	// save only the last "limit" number of records
	var records [][]string
	if limit > 0 && len(allRecords) > (limit+1) {
		titleLine := allRecords[0]
		// erase rows #0 ~ #(len(records) - limit - 1)
		records := allRecords[(len(allRecords) - limit):]
		records = append([][]string{titleLine}, records...)
	} else {
		records = allRecords
	}

	for _, sec := range allSections {
		response.Data[sec] = NewMetricsInfoSection()
	}

	if err = response.formResponseData(records); err != nil {
		response.Data = nil
		response.ErrMessage = "Unable to form response for Error: " + err.Error()
		return response
	}

	response.Success = true
	return response
}

func readCSV(filepath string) ([][]string, error) {
	if file, err := os.Open(filepath); err != nil {
		return nil, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	return csvReader.ReadAll()
}

func (response *MetricsInfoResponse) formResponseData(csvArr [][]string) error {
	if len(csvArr) < 2 {
		return ErrNotEnoughDataInCSV
	}
	// block
	if blkPlotsData, err := formPlotDataArr(csvArr, []string{"block:height", "block:txAddToBlockCost", "block:txPoolSize"}); err != nil {
		return err
	}
	response.Data["block"].PopulateData(nil, blkPlotsData)

	// cpu
	if cpuConstants, err := getColumnsInFloat(csvArr, []string{"cpu:totalCoreNum"}); err != nil {
		return err
	}
	if cpuPlotData, err := formPlotDataArrWithTime(csvArr, []string{"cpu:totalProcessCpuPercent", "cpu:currentProcessCpuPercent"}); err != nil {
		return err
	}
	response.Data["cpu"].PopulateData(cpuConstants, cpuPlotData)

	// disk
	if diskPlotData, err := formPlotDataArrWithTime(csvArr, []string{"disk:used", "disk:UsedChange", "disk:usedPercent", "disk:readBytes", "disk:writeBytes"}); err != nil {
		return err
	}
	response.Data["disk"].PopulateData(nil, diskPlotData)

	// memory
	if memoryConstants, err := getColumnsInFloat(csvArr, []string{"memory:systemMem"}); err != nil {
		return err
	}
	if memoryPlotData, err := formPlotDataArrWithTime(csvArr, []string{"memory:currentProcessMemInUse", "memory:totalProcessMemInUse", "memory:currentProcessMemPercent", "memory:totalProcessMemPercent"}); err != nil {
		return err
	}
	response.Data["memory"].PopulateData(memoryConstants, memoryPlotData)

	// network
	if networkConnCSVData, err := formOnePlotDataWithTime(csvArr, []string{"network:connectionTypeInNum", "network:connectionTypeOutNum"}); err != nil {
		return err
	}
	if networkBytesCSVData, err := formOnePlotDataWithTime(csvArr, []string{"network:bytesRecv", "network:bytesSent"}); err != nil {
		return err
	}
	if networkPacketsCSVData, err := formOnePlotDataWithTime(csvArr, []string{"network:packetsRecv", "network:packetsSent"}); err != nil {
		return err
	}
	networkPlotData := make(map[string][]interface{})
	networkPlotData["connection"] = networkConnCSVData
	networkPlotData["bytes"] = networkBytesCSVData
	networkPlotData["packets"] = networkPacketsCSVData
	response.Data["network"].PopulateData(nil, networkPlotData)

	// txRequest
	if txRequestConcurrCSVData, err := formOnePlotDataWithTime(csvArr, []string{"txRequest:txRequestSend:concurrent", "txRequest:txRequestSendFromMiner:concurrent"}); err != nil {
		return err
	}
	if txRequestCostTimeCSVData, err := formOnePlotDataWithTime(csvArr, []string{"txRequest:txRequestSend:costTime", "txRequest:txRequestSendFromMiner:costTime"}); err != nil {
		return err
	}
	if txRequestQPSCSVData, err := formOnePlotDataWithTime(csvArr, []string{"txRequest:txRequestSend:qps", "txRequest:txRequestSendFromMiner:qps"}); err != nil {
		return err
	}
	txRequestPlotData := make(map[string][]interface{})
	txRequestPlotData["concurrent"] = txRequestConcurrCSVData
	txRequestPlotData["costTime"] = txRequestCostTimeCSVData
	txRequestPlotData["qps"] = txRequestQPSCSVData
	response.Data["txRequest"].PopulateData(nil, txRequestPlotData)

	return nil
}

func formPlotDataArr(csvArr [][]string, columns []string) (map[string][]interface{}, error) {
	// assume column #1 is x-axis
	if len(columns) < 2 {
		return nil, ErrNotEnoughDataInCSV
	}

	allPlotData := make(map[string][]interface{})
	for i := 1; i < len(columns); i++ {
		col := columns[i]
		if onePlotData, err := getColumnsInFloat(csvArr, []string{columns[0], col}); err != nil {
			return nil, err
		}
		newKey := strings.Join(strings.Split(col, ":")[1:], ":")
		allPlotData[newKey] = onePlotData
	}
	return allPlotData, nil
}

func formPlotDataArrWithTime(csvArr [][]string, columns []string) (map[string][]interface{}, error) {
	// add time as x-axis

	allPlotData := make(map[string][]interface{})
	for i := 0; i < len(columns); i++ {
		col := columns[i]
		if onePlotData, err := formOnePlotDataWithTime(csvArr, []string{col}); err != nil {
			return nil, err
		}
		newKey := strings.Join(strings.Split(col, ":")[1:], ":")
		allPlotData[newKey] = onePlotData
	}
	return allPlotData, nil
}

func formOnePlotDataWithTime(csvArr [][]string, columns []string) ([]interface{}, error) {
	if plotArr, err := getColumnsInFloat(csvArr, columns); err != nil {
		return nil, err
	}
	return appendTimeToDataArr(plotArr)
}

func getColumnsInFloat(csvArr [][]string, columns []string) ([]interface{}, error) {
	// subset csvArr to keep only columns
	var res []interface{}
	if len(csvArr) < 2 {
		return nil, ErrNotEnoughDataInCSV
	}

	var columnsIdx []int
	var newTitleLine []string
	titleLine := csvArr[0]
	for _, col := range columns {
		idx := findIndex(titleLine, col)
		if idx >= 0 {
			columnsIdx = append(columnsIdx, idx)
			prefixRemovedCol := strings.Join(strings.Split(col, ":")[1:], ":")
			newTitleLine = append(newTitleLine, prefixRemovedCol)
		}
	}
	res = append(res, newTitleLine)
	for i := 1; i < len(csvArr); i++ {
		line := csvArr[i]
		if newLine, err := subsetArrByIndicesAndConvertToFloat(line, columnsIdx); err != nil {
			return nil, errors.New("ErrGetColumnsInFloat: " + err.Error())
		}
		res = append(res, newLine)
	}
	return res, nil
}

func findIndex(arr []string, element string) int {
	for idx, item := range arr {
		if item == element {
			return idx
		}
	}
	return -1
}

func subsetArrByIndicesAndConvertToFloat(arr []string, indices []int) ([]float64, error) {
	var res []float64
	for _, i := range indices {
		if converted, err := strconv.ParseFloat(arr[i], 64); err != nil {
			return nil, err
		}
		res = append(res, converted)
	}
	return res, nil
}

func appendTimeToDataArr(dataRows []interface{}) ([]interface{}, error) {
	if len(dataRows) < 2 {
		return nil, ErrNotEnoughDataInCSV
	}

	var updated []interface{}
	titleRow := dataRows[0].([]string)
	titleRow = append([]string{"time"}, titleRow...)
	updated = append(updated, titleRow)

	for i := 1; i < len(dataRows); i++ {
		row := dataRows[i].([]float64)
		row = append([]float64{float64(5 * (i - 1))}, row...)
		updated = append(updated, row)
	}

	return updated, nil
}
