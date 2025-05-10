package helpers

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	serviceName string
	folder      = folderDev
	logText     *logrus.Logger
	logJSON     *logrus.Logger
	Logger      *logrus.Logger
)

type LogResponseParam struct {
	ThirdParty     []ThirdPartyLog
	ResponseHeader map[string]string
	ResponseBody   interface{}
	ResponseCode   string
	Timestamp      string
}

type ThirdPartyLog struct {
	MicroName string
	Request   interface{}
	Response  interface{}
	Timestamp string
}

const (
	httpRequest          = "REQUEST"
	httpResponse         = "RESPONSE"
	timeformat           = "2006-01-02T15:04:05-0700"
	nameformat           = "log-2006-01-02.log"
	folderDev            = "storage/logs/"
	folderOcp            = "storage/logs/"
	nameformatthirdparty = "thirdparty-2006-01-02.log"
)

func init() {
	setText()
	setJSON()
	serviceName = "POS Service"
	logText.SetLevel(logrus.InfoLevel)
	logJSON.SetLevel(logrus.InfoLevel)
}

func setJSON() {
	logJSON = logrus.New()
	formatter := new(logrus.JSONFormatter)
	formatter.DisableTimestamp = true
	logJSON.SetFormatter(formatter)
}

func setText() {
	logText = logrus.New()
	formatter := new(logrus.TextFormatter)
	formatter.DisableTimestamp = true
	formatter.DisableQuote = true
	logText.SetFormatter(formatter)
}

func LogRequest(httpMethod, trx_type string, request interface{}, header map[string]string) {
	timestamp := setLogFile()
	logJSON.WithFields(logrus.Fields{
		"service":        serviceName,
		"http_type":      httpRequest,
		"http_method":    httpMethod,
		"request_header": header,
		"request_body":   request,
		"trx_type":       trx_type,
		"timestamp":      timestamp,
	}).Info(httpRequest)

}

func setLogFile() string {
	currentTime := time.Now()
	timestamp := currentTime.Format(timeformat)
	filename := folder + currentTime.Format(nameformat)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	} else {
		logText.SetOutput(file)
		logJSON.SetOutput(file)
	}
	return timestamp
}

func LogResponse(param LogResponseParam) {
	timestamp := setLogFile()
	logJSON.WithFields(logrus.Fields{
		"third_party":     param.ThirdParty,
		"service":         serviceName,
		"http_type":       httpResponse,
		"response_header": param.ResponseHeader,
		"response_body":   MinifyJson(param.ResponseBody),
		"response_code":   param.ResponseCode,
		"timestamp":       timestamp,
	}).Info(httpResponse)
}

func setThirdPartyLogFile() string {
	currentTime := time.Now()
	timestamp := currentTime.Format(timeformat)
	filename := folder + currentTime.Format(nameformatthirdparty)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	} else {
		logText.SetOutput(file)
		logJSON.SetOutput(file)
	}
	return timestamp
}

func LogThirdParty(microname, method, url, fitur string, requestBody, requestHeader, responseBody, responseHeader, responseTime any) {
	timestamp := setThirdPartyLogFile()
	logJSON.WithFields(logrus.Fields{
		"micro_name":      microname,
		"method":          method,
		"path":            url,
		"fitur":           fitur,
		"request_body":    requestBody,
		"request_header":  requestHeader,
		"response_body":   responseBody,
		"response_header": responseHeader,
		"timestamp":       timestamp,
		"response_time":   responseTime,
	}).Info(httpResponse)
}
