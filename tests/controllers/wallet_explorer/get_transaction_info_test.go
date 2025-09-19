package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	controllers "cry-api/app/controllers/wallet_explorer"
	WalletExplorerTypes "cry-api/app/types/wallet_explorer"
	testmocks "cry-api/tests/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestWalletExplorerController_GetTransactionInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockTransactionService := new(testmocks.MockTransactionService)

	controller := &controllers.WalletExplorerController{
		TransactionService: mockTransactionService,
	}

	makeRequest := func(query string) (*gin.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodGet, "/transaction?"+query, nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		return c, w
	}

	t.Run("Missing txid parameter", func(t *testing.T) {
		c, w := makeRequest("")
		controller.GetTransactionInfo(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"Missing txid parameter"}`, w.Body.String())
	})

	t.Run("External service error", func(t *testing.T) {
		txid := "txid123"
		mockTransactionService.On("GetTransactionByTxID", txid).
			Return(nil, assert.AnError).
			Once()

		c, w := makeRequest("txid=" + txid)
		controller.GetTransactionInfo(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockTransactionService.AssertExpectations(t)
	})

	t.Run("Successful call", func(t *testing.T) {
		txid := "txid123"
		mockData := &WalletExplorerTypes.ITransactionData{
			Hash:        "txid123",
			Ver:         1,
			VinSz:       1,
			VoutSz:      1,
			LockTime:    0,
			Size:        200,
			RelayedBy:   "node1",
			BlockHeight: 1000,
			Time:        1650000000,
			Inputs: []WalletExplorerTypes.Input{
				{
					Script: "input-script",
				},
			},
			Out: []WalletExplorerTypes.Output{
				{
					Script: "output-script",
				},
			},
		}

		mockTransactionService.On("GetTransactionByTxID", txid).
			Return(mockData, nil).
			Once()

		c, w := makeRequest("txid=" + txid)
		controller.GetTransactionInfo(c)

		assert.Equal(t, http.StatusOK, w.Code)

		expectedJSON := `{
			"transaction_data": {
				"hash": "txid123",
				"ver": 1,
				"vin_sz": 1,
				"vout_sz": 1,
				"lock_time": 0,
				"size": 200,
				"relayed_by": "node1",
				"block_height": 1000,
				"time": 1650000000,
				"tx_index": null,
				"inputs": [{"script": "input-script"}],
				"out": [{"script": "output-script", "tx_index": null, "value": null}]
			}
		}`
		assert.JSONEq(t, expectedJSON, w.Body.String())

		mockTransactionService.AssertExpectations(t)
	})
}
