package http

import (
	"encoding/json"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/outputs"
	"net/http"
	"crypto/tls"
	"bytes"
	"log"
	"fmt"
	"errors"
)

var sampleConfig = `
  url = "http://localhost"
  method = "POST"
  authorization_key = ""
`

type Http struct {
	Url              string
	Method           string
	AuthorizationKey string `toml:"authorization_key"`
}

type Metric struct {
	Timestamp   int64                  `json:"timestamp"`
	Measurement string                 `json:"measurement"`
	Tags        map[string]string      `json:"tags"`
	Fields      map[string]interface{} `json:"fields"`
}

type RequestData struct {
	Metrics []Metric `json:"metrics"`
}

type Request struct {
	Data RequestData `json:"data"`
}

func (r *RequestData) Append(metric telegraf.Metric) {
	item := &Metric{
		Timestamp:   metric.Time().UnixNano() / 1000000000, // to seconds
		Measurement: metric.Name(),
		Tags:        metric.Tags(),
		Fields:      metric.Fields(),
	}

	r.Metrics = append(r.Metrics, *item)
}

func (h *Http) Description() string {
	return "Send telegraf metrics to http endpoint"
}

func (h *Http) SampleConfig() string {
	return sampleConfig
}

func (h *Http) Connect() error {
	return nil
}

func (h *Http) Close() error {
	return nil
}

func (h *Http) Write(metrics []telegraf.Metric) error {
	log.Printf("I! Attempt to send data: %s", h.Url)

	requestData := &RequestData{}
	for _, metric := range metrics {
		requestData.Append(metric)
	}
	request := &Request{Data: *requestData}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(h.Method, h.Url, bytes.NewBuffer(jsonRequest))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if len(h.AuthorizationKey) > 0 {
		req.Header.Set("Authorization", fmt.Sprintf("Key %s", h.AuthorizationKey))
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		msg := fmt.Sprintf("Incorrect response code %d", resp.StatusCode)
		return errors.New(msg)
	}

	log.Printf("I! Sent metrics: %d", len(requestData.Metrics))

	return nil
}

func init() {
	outputs.Add("http", func() telegraf.Output { return &Http{} })
}
