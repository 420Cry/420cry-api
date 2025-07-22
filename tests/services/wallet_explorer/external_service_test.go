package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	services "cry-api/app/services/wallet_explorer"
	EnvTypes "cry-api/app/types/env"

	"github.com/stretchr/testify/assert"
)

// helper to create EnvConfig with test URLs
func makeTestEnvConfig(apiURL, walletExplorerURL string) *EnvTypes.EnvConfig {
	return &EnvTypes.EnvConfig{
		BlockchainConfig: EnvTypes.BlockchainConfig{
			API: apiURL,
		},
		WalletExplorerConfig: EnvTypes.WalletExplorerConfig{
			API: walletExplorerURL,
		},
	}
}

func TestGetTransactionByTxID_Success(t *testing.T) {
	// Mock server to simulate Blockchain API
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/rawtx/testtxid", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{
			"hash": "testtxid",
			"ver": 1,
			"vin_sz": 1,
			"vout_sz": 1,
			"lock_time": 0,
			"size": 200,
			"relayed_by": "node1",
			"block_height": 123,
			"time": 1650000000,
			"inputs": [{"script": "input-script"}],
			"out": [{"script": "output-script"}]
		}`)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	cfg := makeTestEnvConfig(server.URL, "")
	svc := services.NewExternalService(cfg)

	data, err := svc.GetTransactionByTxID("testtxid")
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, "testtxid", data.Hash)
	assert.Equal(t, 123, data.BlockHeight)
	assert.Len(t, data.Inputs, 1)
	assert.Equal(t, "input-script", data.Inputs[0].Script)
}

func TestGetTransactionByTxID_Non200Status(t *testing.T) {
	handler := func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`Bad request`))
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	cfg := makeTestEnvConfig(server.URL, "")
	svc := services.NewExternalService(cfg)

	data, err := svc.GetTransactionByTxID("testtxid")
	assert.Error(t, err)
	assert.Nil(t, data)
	assert.Contains(t, err.Error(), "external API returned status 400")
}

func TestGetTransactionByTxID_InvalidJSON(t *testing.T) {
	handler := func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`invalid json`))
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	cfg := makeTestEnvConfig(server.URL, "")
	svc := services.NewExternalService(cfg)

	data, err := svc.GetTransactionByTxID("testtxid")
	assert.Error(t, err)
	assert.Nil(t, data)
	assert.Contains(t, err.Error(), "failed to parse JSON")
}

func TestGetTransactionByXPUB_Success(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/xpub-txs", r.URL.Path)
		assert.Equal(t, "pub=testxpub&gap_limit=5", r.URL.RawQuery)
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{
			"found": true,
			"gap_limit": 5,
			"txs": []
		}`)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	cfg := makeTestEnvConfig("", server.URL)
	svc := services.NewExternalService(cfg)

	data, err := svc.GetTransactionByXPUB("testxpub")
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.True(t, data.Found)
	assert.Equal(t, 5, data.GapLimit)
	assert.Empty(t, data.Transactions)
}

func TestGetTransactionByXPUB_Timeout(t *testing.T) {
	// Simulate server that delays response beyond client timeout (30s)
	handler := func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(35 * time.Second)
		w.WriteHeader(http.StatusOK)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	cfg := makeTestEnvConfig("", server.URL)
	svc := services.NewExternalService(cfg)

	_, err := svc.GetTransactionByXPUB("testxpub")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to fetch data")
}

func TestGetTransactionByXPUB_Non200Status(t *testing.T) {
	handler := func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(`Server error`)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	cfg := makeTestEnvConfig("", server.URL)
	svc := services.NewExternalService(cfg)

	data, err := svc.GetTransactionByXPUB("testxpub")
	assert.Error(t, err)
	assert.Nil(t, data)
	assert.Contains(t, err.Error(), "external API returned status 500")
}

func TestGetTransactionByXPUB_InvalidJSON(t *testing.T) {
	handler := func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`not valid json`)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	cfg := makeTestEnvConfig("", server.URL)
	svc := services.NewExternalService(cfg)

	data, err := svc.GetTransactionByXPUB("testxpub")
	assert.Error(t, err)
	assert.Nil(t, data)
	assert.Contains(t, err.Error(), "failed to parse JSON")
}
