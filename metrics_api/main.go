package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"strconv"

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

	metricsConfig := &configpb.MetricsConfig{}
	config.LoadConfig(filePath, metricsConfig)

	go metricsreader.RunMetricsReader(metricsConfig)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", returnConnSuccess).Methods("GET", "OPTIONS")
	router.HandleFunc("/getLatest/{fromtime}", returnLatestData).Methods("GET", "OPTIONS")
	router.HandleFunc("/getHistory/{filename}", returnHistoryData).Methods("GET", "OPTIONS")
	router.HandleFunc("/getListOfDataFiles", returnListDataFiles).Methods("GET", "OPTIONS")

	logger.Fatalln(http.ListenAndServe(":9000", router))
}

func returnConnSuccess(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode("API is up!")
}

func returnHistoryData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	filename := vars["filename"]
	filepath := dir + filename
	miResponse := metricsreader.FormMetricsInfoResponse(filepath, 0)
	json.NewEncoder(w).Encode(miResponse)
}

func returnLatestData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	fromtime := vars["fromtime"]
	filepath := dir + findLatestDataFilename(dir)
	from, err := strconv.Atoi(fromtime)
	if err != nil {
		failresponse := metricsreader.NewFailMetricsInfoResponse(err.Error())
		json.NewEncoder(w).Encode(failresponse)
	}
	miResponse := metricsreader.FormMetricsInfoResponse(filepath, from)
	json.NewEncoder(w).Encode(miResponse)
}

func returnListDataFiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
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
