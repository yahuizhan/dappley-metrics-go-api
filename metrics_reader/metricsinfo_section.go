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

func (miSection *MetricsInfoSection) PopulateData(consts []interface{}, plotData map[string][]interface{}) {
	if consts != nil {
		miSection.setConsts(consts)
	}
	miSection.PlotData = plotData
}

func (miSection *MetricsInfoSection) setConsts(csvData []interface{}) {
	if len(csvData) < 2 {
		return
	}
	titles := csvData[0].([]string)
	data := csvData[1].([]float64)
	for idx, title := range titles {
		miSection.Consts[title] = data[idx]
	}
}
