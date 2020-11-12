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
	rpcpb "github.com/yahuizhan/dappley-metrics-go-api/rpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var CurrCSVFilename string = ""

// RunMetricsReader requests metrics info every 5 seconds and writes into csv
func RunMetricsReader(cliConfig *configpb.CliConfig) {

	conn := initRPCClient(int(cliConfig.GetPort()))
	defer conn.Close()

	metricsRPCService := rpcpb.NewMetricServiceClient(conn)

	md := metadata.Pairs("password", cliConfig.GetPassword())
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
	timeNow := time.Now().Format("2006Jan02150405")
	fmt.Println(timeNow)
	filepath := "csv/metricsInfo_result" + timeNow + ".csv"
	// write to csv
	file, err := os.Create(filepath)
	if err != nil {
		logger.Panic("cannot create a csv file due to Error -- " + err.Error())
	}
	file.Close()
	CurrCSVFilename = filepath

	var columnTitles []string
	isColumnTitlesWritten := false

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
				default:
					logger.Error("Error:", err.Error())
				}
				return
			}
			logger.Info("metricsInfo:", metricsInfoResponse.Data)

			m, ok := gjson.Parse(metricsInfoResponse.Data).Value().(map[string]interface{})
			if !ok {
				logger.Warning("parse data is not json")
				continue
			}

			metricsInfoMap := make(map[string]string)
			formCSVRecord(m, "", metricsInfoMap) // add error handling
			printMap(metricsInfoMap)

			// open file
			file, err = os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			writer := csv.NewWriter(file)
			writer.Comma = ','
			// start writing
			var metricsInfostr []string
			if !isColumnTitlesWritten {
				for k := range metricsInfoMap {
					columnTitles = append(columnTitles, k)
				}
				sort.Strings(columnTitles)
				writer.Write(columnTitles)
				isColumnTitlesWritten = true
			}
			for _, col := range columnTitles {
				metricsInfostr = append(metricsInfostr, metricsInfoMap[col])
			}
			writer.Write(metricsInfostr)
			writer.Flush()

			file.Close()
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

func printMap(m map[string]string) {
	logger.Info("Printing map.......")
	for k, v := range m {
		fmt.Printf("key is %v, value is %v\n", k, v)
	}
}

func isIntegral(val float64) bool {
	return val == float64(int(val))
}
