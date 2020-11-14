package metricsreader

import (
	"strings"

	csvhandler "github.com/yahuizhan/dappley-metrics-go-api/metrics_reader/csv_handler"
)

var allSections []string = []string{"block", "cpu", "disk", "memory", "network", "txRequest"}

/* var constantFieldsMap = map[string]string{
	"cpu":    "totalCoreNum",
	"memory": "systemMem",
}

type plotInfo struct {
	title       string
	inputTitles []string // x, y1, y2, ...
}

var plotDataMap = map[string][]plotInfo{
	"block": {
		plotInfo{"txAddToBlockCost", []string{"block:height", "block:txAddToBlockCost"}},
		plotInfo{"txPoolSize", []string{"block:height", "block:txPoolSize"}},
	},
	"cpu": {
		plotInfo{"currentProcessCpuPercent", []string{"time", "cpu:currentProcessCpuPercent"}},
		plotInfo{"totalProcessCpuPercent", []string{"time", "cpu:totalProcessCpuPercent"}},
	},
	"disk": {
		plotInfo{"used", []string{"time", "disk:used"}},
		plotInfo{"UsedChange", []string{"time", "disk:UsedChange"}},
		plotInfo{"usedPercent", []string{"time", "disk:usedPercent"}},
		plotInfo{"readBytes", []string{"time", "disk:readBytes"}},
		plotInfo{"writeBytes", []string{"time", "disk:writeBytes"}},
	},
	"memory": {
		plotInfo{"currentProcessMemInUse", []string{"time", "memory:currentProcessMemInUse"}},
		plotInfo{"currentProcessMemPercent", []string{"time", "memory:currentProcessMemPercent"}},
		plotInfo{"totalProcessMemInUse", []string{"time", "memory:totalProcessMemInUse"}},
		plotInfo{"totalProcessMemPercent", []string{"time", "memory:totalProcessMemPercent"}},
	},
	"network": {
		plotInfo{"connection", []string{"time", "network:connectionTypeInNum", "network:connectionTypeOutNum"}},
		plotInfo{"bytes", []string{"time", "network:bytesSent", "network:bytesRecv"}},
		plotInfo{"packets", []string{"time", "network:packetsSent", "network:packetsRecv"}},
	},
	"txRequest": {
		plotInfo{"concurrent", []string{"time", "txRequest:txRequestSend:concurren", "txRequest:txRequestSendFromMiner:concurrent"}},
		plotInfo{"costTime", []string{"time", "txRequest:txRequestSend:costTime", "txRequest:txRequestSendFromMiner:costTime"}},
		plotInfo{"qps", []string{"time", "txRequest:txRequestSend:qps", "txRequest:txRequestSendFromMiner:qps"}},
	},
} */

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

	allRecords, err := csvhandler.ReadCSV(filepath)
	if err != nil {
		response.Data = nil
		response.ErrMessage = "Unable to parse file as CSV for Error: " + err.Error()
		return response
	}

	// save only the last "limit" number of records
	var records [][]string
	if limit > 0 && len(allRecords) > (limit+1) {
		titleLine := allRecords[0]
		// erase rows #0 ~ #(len(records) - limit - 1)
		records = allRecords[(len(allRecords) - limit):]
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

func (response *MetricsInfoResponse) formResponseData(csvArr [][]string) error {
	// block
	blkPlotsData, err := formPlotDataArr(csvArr, []string{"block:height", "block:txAddToBlockCost", "block:txPoolSize"})
	if err != nil {
		return err
	}
	response.Data["block"].PopulateData(nil, blkPlotsData)

	// cpu
	cpuConstants, err := csvhandler.GetColumnsInFloat(csvArr, []string{"cpu:totalCoreNum"})
	if err != nil {
		return err
	}
	cpuPlotData, err := formPlotDataArrWithTime(csvArr, []string{"cpu:totalProcessCpuPercent", "cpu:currentProcessCpuPercent"})
	if err != nil {
		return err
	}
	response.Data["cpu"].PopulateData(cpuConstants, cpuPlotData)

	// disk
	diskPlotData, err := formPlotDataArrWithTime(csvArr, []string{"disk:used", "disk:UsedChange", "disk:usedPercent", "disk:readBytes", "disk:writeBytes"})
	if err != nil {
		return err
	}
	response.Data["disk"].PopulateData(nil, diskPlotData)

	// memory
	memoryConstants, err := csvhandler.GetColumnsInFloat(csvArr, []string{"memory:systemMem"})
	if err != nil {
		return err
	}
	memoryPlotData, err := formPlotDataArrWithTime(csvArr, []string{"memory:currentProcessMemInUse", "memory:totalProcessMemInUse", "memory:currentProcessMemPercent", "memory:totalProcessMemPercent"})
	if err != nil {
		return err
	}
	response.Data["memory"].PopulateData(memoryConstants, memoryPlotData)

	// network
	networkConnCSVData, err := formOnePlotDataWithTime(csvArr, []string{"network:connectionTypeInNum", "network:connectionTypeOutNum"})
	if err != nil {
		return err
	}
	networkBytesCSVData, err := formOnePlotDataWithTime(csvArr, []string{"network:bytesRecv", "network:bytesSent"})
	if err != nil {
		return err
	}
	networkPacketsCSVData, err := formOnePlotDataWithTime(csvArr, []string{"network:packetsRecv", "network:packetsSent"})
	if err != nil {
		return err
	}
	networkPlotData := make(map[string][]interface{})
	networkPlotData["connection"] = networkConnCSVData
	networkPlotData["bytes"] = networkBytesCSVData
	networkPlotData["packets"] = networkPacketsCSVData
	response.Data["network"].PopulateData(nil, networkPlotData)

	// txRequest
	txRequestConcurrCSVData, err := formOnePlotDataWithTime(csvArr, []string{"txRequest:txRequestSend:concurrent", "txRequest:txRequestSendFromMiner:concurrent"})
	if err != nil {
		return err
	}
	txRequestCostTimeCSVData, err := formOnePlotDataWithTime(csvArr, []string{"txRequest:txRequestSend:costTime", "txRequest:txRequestSendFromMiner:costTime"})
	if err != nil {
		return err
	}
	txRequestQPSCSVData, err := formOnePlotDataWithTime(csvArr, []string{"txRequest:txRequestSend:qps", "txRequest:txRequestSendFromMiner:qps"})
	if err != nil {
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
	allPlotData := make(map[string][]interface{})
	for i := 1; i < len(columns); i++ {
		col := columns[i]
		onePlotData, err := csvhandler.GetColumnsInFloat(csvArr, []string{columns[0], col})
		if err != nil {
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
		onePlotData, err := formOnePlotDataWithTime(csvArr, []string{col})
		if err != nil {
			return nil, err
		}
		newKey := strings.Join(strings.Split(col, ":")[1:], ":")
		allPlotData[newKey] = onePlotData
	}
	return allPlotData, nil
}

func formOnePlotDataWithTime(csvArr [][]string, columns []string) ([]interface{}, error) {
	plotArr, err := csvhandler.GetColumnsInFloat(csvArr, columns)
	if err != nil {
		return nil, err
	}
	return csvhandler.AppendTimeToDataArr(plotArr)
}
