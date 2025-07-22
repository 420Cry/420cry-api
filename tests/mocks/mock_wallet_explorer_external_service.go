package mocks

import (
	WalletExplorer "cry-api/app/types/wallet_explorer"

	"github.com/stretchr/testify/mock"
)

// MockExternalService mocks the ExternalServiceInterface for testing purposes.
type MockExternalService struct {
	mock.Mock
}

// GetTransactionByXPUB mocks the GetTransactionByXPUB method of the ExternalServiceInterface.
func (m *MockExternalService) GetTransactionByXPUB(xpub string) (*WalletExplorer.ITransactionXPUB, error) {
	args := m.Called(xpub)
	if result := args.Get(0); result != nil {
		return result.(*WalletExplorer.ITransactionXPUB), args.Error(1)
	}
	return nil, args.Error(1)
}

// GetTransactionByTxID mocks the GetTransactionByTxID method of the ExternalServiceInterface.
func (m *MockExternalService) GetTransactionByTxID(txid string) (*WalletExplorer.ITransactionData, error) {
	args := m.Called(txid)
	if result := args.Get(0); result != nil {
		return result.(*WalletExplorer.ITransactionData), args.Error(1)
	}
	return nil, args.Error(1)
}
