package middleware

import (
	"checkout-service/internal/helpers"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
)

const (
	nameformat = "log-2006-01-02.log"
	timeformat = "2006-01-02T15:04:05-0700"
	folder     = "storage/logs/"
)

func CustomLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		req := c.Request()

		var reqBody map[string]any
		json.Unmarshal(req.Body(), &reqBody)
		helpers.LogRequest(c.Method(), c.Path(), reqBody, c.GetReqHeaders())
		// Log informasi khusus sebelum permintaan diproses
		log.Printf("[%s] %s %s\n", c.Method(), c.Path(), c.IP())

		// Next akan memanggil middleware atau handler berikutnya
		err := c.Next()

		// Log informasi khusus setelah permintaan selesai
		log.Printf("[%s] %s %s - %v\n", c.Method(), c.Path(), c.IP(), time.Since(start))

		return err
	}
}

func LogRequest() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// buf := new(bytes.Buffer)
		l := logrus.New()
		l.SetFormatter(&logrus.TextFormatter{
			ForceColors: true,
		})
		req := c.Request()
		res := c.Response()

		logger := logrus.NewEntry(l)
		currentTime := time.Now()
		filename := folder + currentTime.Format(nameformat)
		file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println(err)
		} else {
			logrus.SetOutput(file)
		}

		formatter := new(logrus.JSONFormatter)
		formatter.DisableTimestamp = true
		logrus.SetFormatter(formatter)

		var reqBody map[string]any
		var resBody map[string]any
		json.Unmarshal(req.Body(), &reqBody)
		json.Unmarshal(res.Body(), &resBody)

		logrus.WithFields(logrus.Fields{
			"http_type":      "REQUEST",
			"method":         string(req.Header.Method()),
			"uri":            req.URI(),
			"ip":             c.IP(),
			"host":           string(req.Host()),
			"request_body":   helpers.MinifyJson(reqBody),
			"request_header": c.GetReqHeaders(),
		}).Info()

		logger = logger.WithFields(logrus.Fields{
			"method":   string(req.Header.Method()),
			"uri":      req.URI(),
			"ip":       c.IP(),
			"host":     string(req.Host()),
			"request":  helpers.MinifyJson(req.Body()),
			"response": res.Body(),
		})

		logger.Infoln("Incoming request")

		if err := c.Next(); err != nil {
			logger = logger.WithFields(logrus.Fields{
				"error": err,
			})
			logger.Errorln("Request failed")
			return err
		}

		logger = logger.WithFields(logrus.Fields{
			"status": res.StatusCode(),
			"length": res.Header.ContentLength(),
		})

		logger.Infoln("Request completed")

		return nil
	}
}
