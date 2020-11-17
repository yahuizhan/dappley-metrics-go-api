package metricsreader

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	configpb "github.com/yahuizhan/dappley-metrics-go-api/config/pb"
	csvhandler "github.com/yahuizhan/dappley-metrics-go-api/metrics_reader/csv_handler"
	rpcpb "github.com/yahuizhan/dappley-metrics-go-api/rpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// RunMetricsReader requests metrics info every 5 seconds and writes into csv
func RunMetricsReader(config *configpb.CliConfig) {

	conn := initRPCClient(int(config.GetPort()))
	defer conn.Close()

	metricsRPCService := rpcpb.NewMetricServiceClient(conn)

	md := metadata.Pairs("password", config.GetPassword())
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	getMetricsInfo(ctx, metricsRPCService)
}

func initRPCClient(port int) *grpc.ClientConn {
	//prepare grpc account
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(fmt.Sprint(":", port), grpc.WithInsecure())
	if err != nil {
		logger.Panic("Error:", err.Error())
	}
	return conn
}

func getMetricsInfo(ctx context.Context, c interface{}) {
	tick := time.NewTicker(time.Duration(5000) * time.Millisecond)
	for {
		select {
		case <-tick.C:
			metricsServiceRequest := rpcpb.MetricsServiceRequest{}
			metricsInfoResponse, err := c.(rpcpb.MetricServiceClient).RpcGetMetricsInfo(ctx, &metricsServiceRequest)
			if err != nil {
				switch status.Code(err) {
				case codes.Unavailable:
					logger.Error("Error: server is not reachable!")
					continue
				default:
					logger.Error("Error:", err.Error())
				}
				return
			}
			logger.Info("Received metricsInfoResponse!")

			m, ok := gjson.Parse(metricsInfoResponse.Data).Value().(map[string]interface{})
			if !ok {
				logger.Warning("parse data is not json")
				continue
			}

			metricsInfoMap := make(map[string]string)
			formCSVRecord(m, "", metricsInfoMap) // add error handling
			//printMap(metricsInfoMap)

			today := time.Now().Format("20060102")
			filepath := "csv/metricsInfo_result" + today + ".csv"
			saveToCSV(filepath, metricsInfoMap, time.Now().Unix())
		}
	}
}

func formCSVRecord(jsonObj map[string]interface{}, titlePrefix string, metricsInfoMap map[string]string) {
	if jsonObj == nil {
		return
	}

	for key, value := range jsonObj {
		prefix := key
		if titlePrefix != "" {
			prefix = titlePrefix + ":" + prefix
		}
		if value != nil {
			switch v := value.(type) {
			case int:
				metricsInfoMap[prefix] = strconv.Itoa(v)
			case float64:
				if isIntegral(v) {
					metricsInfoMap[prefix] = strconv.Itoa(int(v))
				} else {
					metricsInfoMap[prefix] = fmt.Sprintf("%.4f", v)
				}
			default:
				formCSVRecord(v.(map[string]interface{}), prefix, metricsInfoMap)
			}
		}
	}
}

// assume header of csv file is sorted
func saveToCSV(csvFilepath string, metricsInfoMap map[string]string, time int64) {
	var (
		columnTitles []string // will define order of incoming data
		file         *os.File // csv file to be written to
	)

	mapKeys := getMapKeys(metricsInfoMap)
	sort.Strings(mapKeys)
	mapKeys = append([]string{"time"}, mapKeys...)
	metricsInfoMap["time"] = strconv.Itoa(int(time))

	if csvhandler.IsCSVExistAndNonEmpty(csvFilepath) {
		csvData, err := csvhandler.ReadCSV(csvFilepath)
		if err != nil {
			logger.Panic("cannot read csv file for Error: " + err.Error())
		}
		if len(csvData) < 2 { // csv file only contains header line: overwrite it
			columnTitles = mapKeys
			file = csvhandler.CreateNewCSVWithTitles(csvFilepath, columnTitles)
		} else {
			headerline := csvData[0]
			colsToAdd := csvhandler.ArrDifference(mapKeys, headerline)
			if len(colsToAdd) > 0 { // metricsInfoMap has columns that are not found in csv
				columnTitles = csvhandler.ArrUnionSorted(mapKeys, headerline)
				columnTitles = putElementToArrStart(columnTitles, "time")
				newTable, _ := csvhandler.GenerateTableByTitles(csvData, columnTitles)
				file = csvhandler.CreateNewCSVWithTable(csvFilepath, newTable)
			} else { // csv has all columns of metricsInfoMap
				file, err = os.OpenFile(csvFilepath, os.O_WRONLY|os.O_APPEND, 0666)
				columnTitles = headerline
			}
			if err != nil {
				logger.Panic("cannot open csv file for Error: " + err.Error())
			}
		}
	} else {
		columnTitles = mapKeys
		file = csvhandler.CreateNewCSVWithTitles(csvFilepath, columnTitles)
	}

	defer file.Close()
	writer := csv.NewWriter(file)
	writer.Comma = ','
	// write new data to csv
	var metricsInfostr []string
	for _, col := range columnTitles {
		metricsInfostr = append(metricsInfostr, metricsInfoMap[col])
	}
	writer.Write(metricsInfostr)
	writer.Flush()
}

func printMap(m map[string]string) {
	logger.Info("Printing map.......")
	for k, v := range m {
		fmt.Printf("%v: %v\n", k, v)
	}
}

func isIntegral(val float64) bool {
	return val == float64(int(val))
}

func getMapKeys(m map[string]string) []string {
	if m == nil {
		return nil
	}
	var res []string
	for k := range m {
		res = append(res, k)
	}
	return res
}

func putElementToArrStart(arr []string, element string) []string {
	// find first appearance of element in arr and put it to start of arr, with other elements' order unchanged
	// if element cannot be found in arr, return arr itself
	res := []string{element}
	foundElement := false
	for _, e := range arr {
		if e == element && !foundElement {
			foundElement = true
		} else {
			res = append(res, e)
		}
	}
	if foundElement {
		return res
	}
	return arr
}
