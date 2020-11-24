package metricsreader

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
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
func RunMetricsReader(config *configpb.MetricsConfig) {

	conn := initRPCClient(config.GetHost(), int(config.GetPort()))
	defer conn.Close()

	metricsRPCService := rpcpb.NewMetricServiceClient(conn)

	md := metadata.Pairs("password", config.GetPassword())
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	csvCols := config.GetMetricsinfoCsvColumns()
	getMetricsInfo(ctx, metricsRPCService, csvCols)
}

func initRPCClient(host string, port int) *grpc.ClientConn {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(fmt.Sprint(host+":", port), grpc.WithInsecure())
	if err != nil {
		logger.Panic("Error:", err.Error())
	}
	logger.Infof("Connected to %v:%v", host, port)
	return conn
}

func getMetricsInfo(ctx context.Context, c rpcpb.MetricServiceClient, csvCols []string) {
	// set up csv file to store metricsInfo
	today := time.Now().Format("20060102")
	filepath := "csv/metricsInfo_result" + today + ".csv"
	var file *os.File
	var err error
	if csvhandler.IsCSVExistAndNonEmpty(filepath) {
		file, err = os.OpenFile(filepath, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			logger.Panic("Error:", err.Error())
		}
		writer := csv.NewWriter(file)
		emptyLine := make([]string, len(csvCols))
		for idx, col := range csvCols {
			if col == "time" {
				emptyLine[idx] = strconv.Itoa(int(time.Now().Unix()))
			} else {
				emptyLine[idx] = ""
			}
		}
		writer.Write(emptyLine)
		writer.Flush()
	} else {
		file = csvhandler.CreateNewCSVWithTitles(filepath, csvCols)
	}
	defer file.Close()

	// requesting metricsInfo every 5s and save to csv
	tick := time.NewTicker(time.Duration(5000) * time.Millisecond)
	for {
		select {
		case <-tick.C:
			metricsServiceRequest := rpcpb.MetricsServiceRequest{}
			metricsInfoResponse, err := c.RpcGetMetricsInfo(ctx, &metricsServiceRequest)
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
			timeNow := time.Now().Unix()
			logger.Info("Received metricsInfoResponse!")

			m, ok := gjson.Parse(metricsInfoResponse.Data).Value().(map[string]interface{})
			if !ok {
				logger.Warning("parse data is not json")
				continue
			}

			metricsInfoMap := make(map[string]string)
			metricsInfoMap["time"] = strconv.Itoa(int(timeNow))
			formCSVRecord(m, "", metricsInfoMap) // add error handling
			//printMap(metricsInfoMap)

			saveToCSV(file, metricsInfoMap, csvCols)
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
					metricsInfoMap[prefix] = fmt.Sprintf("%.2f", v)
				}
			default:
				formCSVRecord(v.(map[string]interface{}), prefix, metricsInfoMap)
			}
		}
	}
}

// assume header of csv file is sorted
func saveToCSV(file *os.File, metricsInfoMap map[string]string, columnTitles []string) {
	writer := csv.NewWriter(file)
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
