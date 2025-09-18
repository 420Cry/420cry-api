package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	services "cry-api/app/services/coin_market_cap"
	CoinMarketCap "cry-api/app/types/coin_market_cap"
	EnvTypes "cry-api/app/types/env"

	"github.com/stretchr/testify/assert"
)

// helper to create EnvConfig with test URLs and API key
func makeTestEnvConfig(apiURL, apiKey string) *EnvTypes.EnvConfig {
	return &EnvTypes.EnvConfig{
		CoinMarketCapConfig: EnvTypes.CoinMarketCapConfig{
			API:    apiURL,
			APIKey: apiKey,
		},
	}
}

func TestGetFearAndGreedLastest_Success(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v3/fear-and-greed/latest", r.URL.Path)
		assert.Equal(t, "test-api-key", r.Header.Get("X-CMC_PRO_API_KEY"))

		w.WriteHeader(http.StatusOK)
		resp := CoinMarketCap.FearGreedData{
			Data: CoinMarketCap.FearGreedEntry{
				Value:               70,
				ValueClassification: "Greed",
				UpdateTime:          time.Unix(1695000000, 0),
			},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	cfg := makeTestEnvConfig(server.URL, "test-api-key")
	svc := services.NewCoinMarketCapServiceService(cfg)

	data, err := svc.GetFearAndGreedLastest()
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, 70, data.Data.Value)
	assert.Equal(t, "Greed", data.Data.ValueClassification)

	// Compare using Unix timestamp to avoid timezone issues
	assert.Equal(t, int64(1695000000), data.Data.UpdateTime.Unix())
}

func TestGetFearAndGreedLastest_Non200Status(t *testing.T) {
	handler := func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`Bad request`))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	cfg := makeTestEnvConfig(server.URL, "test-api-key")
	svc := services.NewCoinMarketCapServiceService(cfg)

	data, err := svc.GetFearAndGreedLastest()
	assert.Error(t, err)
	assert.Nil(t, data)
	assert.Contains(t, err.Error(), "API request failed with status 400")
}

func TestGetFearAndGreedLastest_InvalidJSON(t *testing.T) {
	handler := func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`invalid json`))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	cfg := makeTestEnvConfig(server.URL, "test-api-key")
	svc := services.NewCoinMarketCapServiceService(cfg)

	data, err := svc.GetFearAndGreedLastest()
	assert.Error(t, err)
	assert.Nil(t, data)
	assert.Contains(t, err.Error(), "failed to decode response body")
}

func TestGetFearAndGreedHistorical_Success(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v3/fear-and-greed/historical", r.URL.Path)
		assert.Equal(t, "test-api-key", r.Header.Get("X-CMC_PRO_API_KEY"))

		w.WriteHeader(http.StatusOK)
		resp := CoinMarketCap.FearGreedHistorical{
			Data: []CoinMarketCap.FearGreedDataPoint{
				{
					Timestamp:           "1695000000",
					Value:               50,
					ValueClassification: "Neutral",
				},
			},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	cfg := makeTestEnvConfig(server.URL, "test-api-key")
	svc := services.NewCoinMarketCapServiceService(cfg)

	data, err := svc.GetFearAndGreedHistorical()
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Len(t, data.Data, 1)
	assert.Equal(t, 50, data.Data[0].Value)
	assert.Equal(t, "Neutral", data.Data[0].ValueClassification)
	assert.Equal(t, "1695000000", data.Data[0].Timestamp)
}

func TestGetFearAndGreedHistorical_Non200Status(t *testing.T) {
	handler := func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`Server error`))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	cfg := makeTestEnvConfig(server.URL, "test-api-key")
	svc := services.NewCoinMarketCapServiceService(cfg)

	data, err := svc.GetFearAndGreedHistorical()
	assert.Error(t, err)
	assert.Nil(t, data)
	assert.Contains(t, err.Error(), "API request failed with status 500")
}

func TestGetFearAndGreedHistorical_InvalidJSON(t *testing.T) {
	handler := func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`invalid json`))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	cfg := makeTestEnvConfig(server.URL, "test-api-key")
	svc := services.NewCoinMarketCapServiceService(cfg)

	data, err := svc.GetFearAndGreedHistorical()
	assert.Error(t, err)
	assert.Nil(t, data)
	assert.Contains(t, err.Error(), "failed to decode response body")
}
