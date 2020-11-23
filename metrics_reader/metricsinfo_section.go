package metricsreader

// MetricsInfoSection formulates data in a way that is easy to plot
type MetricsInfoSection struct {
	Consts   map[string]float64   `json:"consts"`
	PlotInfo map[string]*PlotData `json:"plotInfo"`
}

// NewMetricsInfoSection initializes a new object of MetricsInfoSection
func NewMetricsInfoSection() *MetricsInfoSection {
	return &MetricsInfoSection{
		Consts:   make(map[string]float64),
		PlotInfo: nil,
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
func (miSection *MetricsInfoSection) SetPlotData(plotData map[string]*PlotData) {
	miSection.PlotInfo = plotData
}

// PlotData contains title and data of one plot
type PlotData struct {
	Title string        `json:"title"`
	Data  []interface{} `json:"data"`
}

// NewPlotData returns a pointer to a new PlotData
func NewPlotData(title string, data []interface{}) *PlotData {
	return &PlotData{
		Title: title,
		Data:  data,
	}
}
