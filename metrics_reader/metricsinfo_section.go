package metricsreader

// MetricsInfoSection formulates data in a way that is easy to plot using Google Charts
type MetricsInfoSection struct {
	Consts   map[string]float64       `json:"consts"`
	PlotData map[string][]interface{} `json:"plotData"`
}

// NewMetricsInfoSection initializes a new object of MetricsInfoSection
func NewMetricsInfoSection() *MetricsInfoSection {
	return &MetricsInfoSection{
		Consts:   make(map[string]float64),
		PlotData: nil,
	}
}

// SetConsts sets value to MetricsInfoSection.Consts
func (miSection *MetricsInfoSection) SetConsts(csvData []interface{}) {
	if len(csvData) < 2 {
		return
	}
	titles := csvData[0].([]string)
	for i := 1; i < len(csvData); i++ {
		dataRow := csvData[i].([]*float64)
		for idx, title := range titles {
			if dataRow[idx] != nil {
				miSection.Consts[title] = *(dataRow[idx])
			}
		}
		if len(miSection.Consts) == len(titles) {
			break
		}
	}
}

// SetPlotData sets value to MetricsInfoSection.PlotData
func (miSection *MetricsInfoSection) SetPlotData(plotData map[string][]interface{}) {
	miSection.PlotData = plotData
}
