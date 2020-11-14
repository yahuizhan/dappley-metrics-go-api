package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	logger "github.com/sirupsen/logrus"
	"github.com/yahuizhan/dappley-metrics-go-api/config"
	configpb "github.com/yahuizhan/dappley-metrics-go-api/config/pb"
	metricsreader "github.com/yahuizhan/dappley-metrics-go-api/metrics_reader"
)

var (
	dir string = "./csv/"
)

func main() {

	var filePath string
	flag.StringVar(&filePath, "f", "default.conf", "CLI config file path")
	flag.Parse()

	cliConfig := &configpb.CliConfig{}
	config.LoadConfig(filePath, cliConfig)

	go metricsreader.RunMetricsReader(cliConfig)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/getLatest", returnLatestData)
	router.HandleFunc("/getHistory/{filename}", returnHistoryData)
	router.HandleFunc("/getListOfDataFiles", returnListDataFiles)

	logger.Fatalln(http.ListenAndServe(":9000", router))
}

func returnHistoryData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]
	filepath := dir + filename
	miResponse := metricsreader.FormMetricsInfoResponse(filepath, 0)
	json.NewEncoder(w).Encode(miResponse)
}

func returnLatestData(w http.ResponseWriter, r *http.Request) {
	//filepath := dir + "metricsInfo_result2020Nov09110626.csv"
	filepath := dir + findLatestDataFilename(dir)
	logger.Info(filepath)
	miResponse := metricsreader.FormMetricsInfoResponse(filepath, 10)
	//logger.Info(miResponse)
	//b, _ := json.Marshal(miResponse)
	//logger.Info(string(b))
	json.NewEncoder(w).Encode(miResponse)
}

func returnListDataFiles(w http.ResponseWriter, r *http.Request) {
	filenames, _ := findCSVFiles(dir)
	json.NewEncoder(w).Encode(filenames)
}

func findCSVFiles(dir string) (allFilenames []string, newestFile string) {
	var newestTime int64 = 0
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		filename := f.Name()
		allFilenames = append(allFilenames, filename)
		if f.ModTime().Unix() > newestTime {
			newestFile = filename
			newestTime = f.ModTime().Unix()
		}
	}
	return allFilenames, newestFile
}

func findLatestDataFilename(dir string) string {
	_, latestFilename := findCSVFiles(dir)
	return latestFilename
}
