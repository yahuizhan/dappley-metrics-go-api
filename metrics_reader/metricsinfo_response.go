package metricsreader

import (
	csvhandler "github.com/yahuizhan/dappley-metrics-go-api/metrics_reader/csv_handler"
)

var allSections []string = []string{"block", "cpu", "disk", "memory", "network", "txRequest"}

var constantFieldsMap = map[string][]string{
	"cpu":    {"cpu:totalCoreNum"},
	"memory": {"memory:systemMem"},
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
		plotInfo{"concurrent", []string{"time", "txRequest:txRequestSend:concurrent", "txRequest:txRequestSendFromMiner:concurrent"}},
		plotInfo{"costTime", []string{"time", "txRequest:txRequestSend:costTime", "txRequest:txRequestSendFromMiner:costTime"}},
		plotInfo{"qps", []string{"time", "txRequest:txRequestSend:qps", "txRequest:txRequestSendFromMiner:qps"}},
	},
}

// MetricsInfoResponse forms JSON structure used at the frontend to generate plots
type MetricsInfoResponse struct {
	Success    bool                           `json:"success"`
	ErrMessage string                         `json:"error"`
	Data       map[string]*MetricsInfoSection `json:"data"`
	TimeRange  []string                       `json:"timeRange"`
}

// NewFailMetricsInfoResponse initializes a failure response with given errMessage
func NewFailMetricsInfoResponse(errMessage string) *MetricsInfoResponse {
	return &MetricsInfoResponse{
		Success:    false,
		ErrMessage: errMessage,
		Data:       nil,
		TimeRange:  nil,
	}
}

// FormMetricsInfoResponse initializes and returns a MetricsInfoResponse
func FormMetricsInfoResponse(filepath string, limit int, fromtime int) *MetricsInfoResponse {
	response := &MetricsInfoResponse{
		Success:    false,
		ErrMessage: "",
		Data:       make(map[string]*MetricsInfoSection),
		TimeRange:  nil,
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
	} else if fromtime > 0 {
		records, err = csvhandler.SubsetDataArrByTime(allRecords, fromtime)
		if err != nil {
			response.Data = nil
			response.ErrMessage = "Unable to subset CSV for Error: " + err.Error()
			return response
		}
	} else {
		records = allRecords
	}

	response.TimeRange = findTimeRange(records)

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
	for section, plotInfoArr := range plotDataMap {
		if len(constantFieldsMap[section]) > 0 {
			constants, err := csvhandler.GetColumnsInFloat(csvArr, constantFieldsMap[section])
			if err != nil {
				return err
			}
			response.Data[section].SetConsts(constants)
		}

		allPlotData := make(map[string][]interface{})
		for _, plotinfo := range plotInfoArr {
			onePlotData, err := csvhandler.GetColumnsInFloat(csvArr, plotinfo.inputTitles)
			if err != nil {
				return err
			}
			allPlotData[plotinfo.title] = onePlotData
		}
		response.Data[section].SetPlotData(allPlotData)
	}

	return nil
}

func findTimeRange(csvArr [][]string) []string {
	// get upper and lower limit from first column; assume csvArr is ordered and 1st row is header
	res := []string{csvArr[1][0]}
	res = append(res, csvArr[len(csvArr)-1][0])
	return res
}
