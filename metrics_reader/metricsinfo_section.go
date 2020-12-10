package metricsreader

// MetricsInfoSection formulates data in a way that is easy to plot
type MetricsInfoSection struct {
	Consts []*Constant `json:"consts"`
	Plots  []*PlotData `json:"plots"`
}

// NewMetricsInfoSection initializes a new object of MetricsInfoSection
func NewMetricsInfoSection() *MetricsInfoSection {
	return &MetricsInfoSection{
		Consts: nil,
		Plots:  nil,
	}
}

// SetConsts sets value to MetricsInfoSection.Consts
func (miSection *MetricsInfoSection) SetConsts(constants []*Constant) {
	miSection.Consts = constants
}

// SetPlotData sets value to MetricsInfoSection.PlotData
func (miSection *MetricsInfoSection) SetPlotData(plotData []*PlotData) {
	miSection.Plots = plotData
}

// Constant contains unit and value of a constant metrics
type Constant struct {
	Title    string  `json:"title"`
	UnitType string  `json:"unitType"`
	Value    float64 `json:"value"`
}

// NewConstant returns a pointer to a new Constant
func NewConstant(title string, unitType string, value float64) *Constant {
	return &Constant{
		Title:    title,
		UnitType: unitType,
		Value:    value,
	}
}

// PlotData contains title, data unit and data of one plot
type PlotData struct {
	Title    string        `json:"title"`
	UnitType string        `json:"unitType"`
	Data     []interface{} `json:"data"`
}

// NewPlotData returns a pointer to a new PlotData
func NewPlotData(title string, unitType string, data []interface{}) *PlotData {
	return &PlotData{
		Title:    title,
		UnitType: unitType,
		Data:     data,
	}
}
