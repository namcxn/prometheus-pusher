package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/cxnam/prometheus-pusher/pkg/logger"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/ghodss/yaml"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"

	"github.com/cxnam/prometheus-pusher/pkg/config"
)

type queryConfig map[string]string

var (
	queryConfigFile = flag.String("c", "queries.yaml", "Query config file")
	metricInterval  = flag.Duration("i", 10*time.Second, "Metric push interval")

	// logger     = log.With(log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr)), "caller", log.DefaultCaller)
	log        = logger.GetLogger("push-status-page")
	httpClient = &http.Client{}
)

func main() {
	flag.Parse()
	qConfig := queryConfig{}
	qcd, err := ioutil.ReadFile(*queryConfigFile)
	if err != nil {
		log.Infof("msg", "Couldn't read config file", "error", err.Error())
	}
	if err := yaml.Unmarshal(qcd, &qConfig); err != nil {
		log.Infof("msg", "Couldn't parse config file", "error", err.Error())
	}

	// prometheusURL := fmt.Sprintf(config.Config.Systemmetric.Prometheusurl)

	// client, err := api.NewClient(api.Config{Address: *prometheusURL})
	client, err := api.NewClient(api.Config{Address: config.Config.Systemmetric.Prometheusurl})

	if err != nil {
		log.Infof("msg", "Couldn't create Prometheus client", "error", err.Error())
	}

	api := v1.NewAPI(client)

	for {
		for metricID, query := range qConfig {
			ts := time.Now()

			resp, warnings, err := api.Query(context.Background(), query, ts)
			if err != nil {
				log.Infof("msg", "Couldn't query Prometheus", "error", err.Error())
				continue
			}
			if len(warnings) > 0 {
				fmt.Printf("Warnings: %v\n", warnings)
			}

			vec := resp.(model.Vector)
			if l := vec.Len(); l != 1 {
				log.Infof("msg", "Expected query to return single value", "samples", l)
				continue
			}

			value := vec[0].Value
			if "NaN" == value.String() {
				log.Infof("msg", "Expected query to return", value)
				value = 0
				// continue
			}

			log.Infof(metricID, value)
			if err := sendStatusPage(ts, metricID, float64(value)); err != nil {
				log.Infof("msg", "Couldn't send metric to Statuspage", "error", err.Error())
				continue
			}
		}
		time.Sleep(*metricInterval)
	}
}

func sendStatusPage(ts time.Time, metricID string, value float64) error {
	values := url.Values{
		"data[timestamp]": []string{strconv.FormatInt(ts.Unix(), 10)},
		"data[value]":     []string{strconv.FormatFloat(value, 'f', -1, 64)},
	}
	url := config.Config.Systemmetric.Statuspageurl + path.Join("/v1", "pages", config.Config.Systemmetric.Statuspageid, "metrics", metricID, "data.json")
	req, err := http.NewRequest("POST", url, strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "OAuth "+config.Config.Systemmetric.Statuspagetoken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respStr, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Empty API Error")
		}
		return errors.New("API Error: " + string(respStr))
	}
	return nil
}
