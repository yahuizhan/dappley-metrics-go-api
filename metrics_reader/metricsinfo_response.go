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
	title      string   // plot title
	unitType   string   // unit type of plot; can be one of {"number", "bytes", "percentage"}
	csvColumns []string // data columns from csv to be plotted, ordered as x, y1, y2, ...
}

var plotDataMap = map[string][]plotInfo{
	"block": {
		plotInfo{"Cost for Adding Transaction To Block", "number", []string{"block:height", "block:txAddToBlockCost"}},
		plotInfo{"Size of Transaction Pool", "number", []string{"block:height", "block:txPoolSize"}},
	},
	"cpu": {
		plotInfo{"CPU Percentage of Current Process", "percentage", []string{"time", "cpu:currentProcessCpuPercent"}},
		plotInfo{"Total CPU Percentage of All Processes", "percentage", []string{"time", "cpu:totalProcessCpuPercent"}},
	},
	"disk": {
		plotInfo{"Bytes Read from Disk", "bytes", []string{"time", "disk:readBytes"}},
		plotInfo{"Bytes Written to Disk", "bytes", []string{"time", "disk:writeBytes"}},
		plotInfo{"Disk Used", "bytes", []string{"time", "disk:used"}},
		plotInfo{"Percentage Of Disk Used", "percentage", []string{"time", "disk:usedPercent"}},
		plotInfo{"Change In Disk Used", "number", []string{"time", "disk:UsedChange"}},
	},
	"memory": {
		plotInfo{"Memory In Use By Current Process", "bytes", []string{"time", "memory:currentProcessMemInUse"}},
		plotInfo{"Memory Percentage Used By Current Process", "percentage", []string{"time", "memory:currentProcessMemPercent"}},
		plotInfo{"Total Memory In Use By All Processes", "bytes", []string{"time", "memory:totalProcessMemInUse"}},
		plotInfo{"Total Memory Percentage Used By All Processes", "percentage", []string{"time", "memory:totalProcessMemPercent"}},
	},
	"network": {
		plotInfo{"Bytes Sent through Network", "bytes", []string{"time", "network:bytesSent"}},
		plotInfo{"Bytes Received through Network", "bytes", []string{"time", "network:bytesRecv"}},
		plotInfo{"Packets Transferred through Network", "number", []string{"time", "network:packetsSent", "network:packetsRecv"}},
		plotInfo{"Number of Network Connections", "number", []string{"time", "network:connectionTypeInNum", "network:connectionTypeOutNum"}},
	},
	"txRequest": {
		plotInfo{"Number of Concurrent Transaction Requests", "number", []string{"time", "txRequest:txRequestSend:concurrent", "txRequest:txRequestSendFromMiner:concurrent"}},
		plotInfo{"Cost Time of Transaction Requests", "number", []string{"time", "txRequest:txRequestSend:costTime", "txRequest:txRequestSendFromMiner:costTime"}},
		//plotInfo{"Queries Per Second(QPS) of Transaction Requests", "number", []string{"time", "txRequest:txRequestSend:qps", "txRequest:txRequestSendFromMiner:qps"}},
	},
}

// MetricsInfoResponse forms JSON structure used at the frontend to generate plots
type MetricsInfoResponse struct {
	Success    bool                           `json:"success"`
	ErrMessage string                         `json:"error"`
	Content    map[string]*MetricsInfoSection `json:"content"`
}

// NewFailMetricsInfoResponse initializes a failure response with given errMessage
func NewFailMetricsInfoResponse(errMessage string) *MetricsInfoResponse {
	return &MetricsInfoResponse{
		Success:    false,
		ErrMessage: errMessage,
		Content:    nil,
	}
}

// FormMetricsInfoResponse initializes and returns a MetricsInfoResponse
func FormMetricsInfoResponse(filepath string, fromtime int) *MetricsInfoResponse {
	response := &MetricsInfoResponse{
		Success:    false,
		ErrMessage: "",
		Content:    make(map[string]*MetricsInfoSection),
	}

	allRecords, err := csvhandler.ReadCSV(filepath)
	if err != nil {
		response.Content = nil
		response.ErrMessage = "Unable to parse file as CSV for Error: " + err.Error()
		return response
	}

	// save only the last "limit" number of records
	var records [][]string
	if fromtime > 0 {
		records, err = csvhandler.SubsetDataArrByTime(allRecords, fromtime)
		if err != nil {
			response.Content = nil
			response.ErrMessage = "Unable to subset CSV for Error: " + err.Error()
			return response
		}
	} else {
		records = allRecords
	}

	for _, sec := range allSections {
		response.Content[sec] = NewMetricsInfoSection()
	}

	if err = response.formResponseData(records); err != nil {
		response.Content = nil
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
			response.Content[section].SetConsts(constants)
		}

		allPlotData := make([]*PlotData, len(plotInfoArr))
		for idx, plotinfo := range plotInfoArr {
			onePlotData, err := csvhandler.GetColumnsInFloat(csvArr, plotinfo.csvColumns)
			if err != nil {
				return err
			}
			allPlotData[idx] = NewPlotData(plotinfo.title, plotinfo.unitType, onePlotData)
		}
		response.Content[section].SetPlotData(allPlotData)
	}

	return nil
}
