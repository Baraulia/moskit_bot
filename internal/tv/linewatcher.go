package tv

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/shopspring/decimal"
	"io"
	"moskitbot/internal/repository"
	"moskitbot/pkg/logging"
	"net/http"
	"strings"
	"time"
)

type LineWatcher struct {
	httpClient *http.Client
	uri        string
	ticker     *time.Ticker
	logger     *logging.Logger
	db         repository.Repository
	cache      *redis.Client
}

type Line struct {
	ID          int
	Pair        string
	Description string
	Val         float32
	Typ         LineType
	Timeframe   Interval
}

var client = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func NewLineWatcher(updateFrequency time.Duration, logger *logging.Logger, uri string, db repository.Repository, cache *redis.Client) *LineWatcher {
	httpClient := &http.Client{}
	ticker := time.NewTicker(updateFrequency)

	return &LineWatcher{
		httpClient: httpClient,
		ticker:     ticker,
		uri:        uri,
		logger:     logger,
		db:         db,
		cache:      cache,
	}
}

func (r *LineWatcher) StartWatching(ctx context.Context, errorChan chan error, responseChan chan Alarm) {
	for {
		select {
		case <-r.ticker.C:
			response, err := r.do(context.Background())
			if err != nil {
				errorChan <- err
				continue
			}

			for _, data := range response.Data {
				var lines []Line
				alarm := Alarm{}

				value := make(map[string]decimal.Decimal)

				rowLines := r.cache.

				for i, column := range data.ColumnsValue {
					if column.GreaterThanOrEqual(r.maxValue) {
						if _, ok := r.alarmed[data.Instrument][r.columns[i]]; !ok {
							value[r.columns[i]] = column
						}
					} else if column.LessThanOrEqual(r.minValue) {
						if _, ok := r.alarmed[data.Instrument][r.columns[i]]; !ok {
							value[r.columns[i]] = column
						}
					} else {
						delete(r.alarmed[data.Instrument], r.columns[i])
					}
				}

				if len(value) != 0 {
					actualValue := make(map[string]decimal.Decimal)
					alarm.instrument = data.Instrument
					for key, val := range value {
						if _, ok := r.alarmed[data.Instrument][key]; !ok {
							r.alarmed[data.Instrument][key] = struct{}{}
							listKey := strings.Split(key, "|")
							listKey[1] = ForResponse(listKey[1])
							actualValue[strings.Join(listKey, "")] = val
						}
					}
					alarm.value = actualValue
					responseChan <- alarm
					continue
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func (r *LineWatcher) do(ctx context.Context) (Response, error) {
	response := Response{}
	payload, err := r.payload()
	if err != nil {
		return response, err
	}

	req, err := http.NewRequest("POST", r.uri, payload)
	if err != nil {
		return response, err
	}

	req = req.WithContext(ctx)

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "go-tradingview")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return response, err
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("HTTP error %d: returned %s", resp.StatusCode, raw)
	}

	err = json.Unmarshal(raw, &response)
	return response, err
}

func (r *LineWatcher) payload() (*bytes.Reader, error) {
	data := Request{}
	var stringPairs string

	pairs := r.cache.Get(context.Background(), "pairs")
	if err := pairs.Scan(stringPairs); err != nil {
		return nil, err
	}

	data.Symbols.Tickers = strings.Split(stringPairs, ",")
	data.Symbols.Query.Types = []interface{}{}
	data.Columns = []string{"ask"}

	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(payload), nil
}
