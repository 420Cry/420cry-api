package tests

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	controllers "cry-api/app/controllers/wallet_explorer"
	WalletExplorerTypes "cry-api/app/types/wallet_explorer"
	testmocks "cry-api/tests/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestWalletExplorerController_GetTransactionByXPUB(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockExternalService := new(testmocks.MockExternalService)

	controller := &controllers.WalletExplorerController{
		ExternalService: mockExternalService,
	}

	makeRequest := func(query string) (*gin.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodGet, "/wallet/xpub?"+query, nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		return c, w
	}

	t.Run("Missing xpub parameter", func(t *testing.T) {
		c, w := makeRequest("")
		controller.GetTransactionByXPUB(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"Missing xpub parameter"}`, w.Body.String())
	})

	t.Run("ExternalService returns error", func(t *testing.T) {
		xpub := "testxpub"
		mockExternalService.On("GetTransactionByXPUB", xpub).Return(nil, errors.New("service failure")).Once()

		c, w := makeRequest("xpub=" + xpub)
		controller.GetTransactionByXPUB(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"service failure"}`, w.Body.String())

		mockExternalService.AssertExpectations(t)
	})

	t.Run("Successful call", func(t *testing.T) {
		xpub := "validxpub"
		mockData := &WalletExplorerTypes.ITransactionXPUB{
			Found:        false,
			GapLimit:     0,
			Transactions: nil,
		}
		mockExternalService.On("GetTransactionByXPUB", xpub).Return(mockData, nil).Once()

		c, w := makeRequest("xpub=" + xpub)
		controller.GetTransactionByXPUB(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"xpub":{"found":false,"gap_limit":0,"txs":null}}`, w.Body.String())

		mockExternalService.AssertExpectations(t)
	})
}
