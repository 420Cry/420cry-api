package mocks

import (
	WalletExplorer "cry-api/app/types/wallet_explorer"

	"github.com/stretchr/testify/mock"
)

// MockTransactionService mocks the MockTransactionService for testing purposes.
type MockTransactionService struct {
	mock.Mock
}

// GetTransactionByXPUB mocks the GetTransactionByXPUB method of the MockTransactionService.
func (m *MockTransactionService) GetTransactionByXPUB(xpub string) (*WalletExplorer.ITransactionXPUB, error) {
	args := m.Called(xpub)
	if result := args.Get(0); result != nil {
		return result.(*WalletExplorer.ITransactionXPUB), args.Error(1)
	}
	return nil, args.Error(1)
}

// GetTransactionByTxID mocks the GetTransactionByTxID method of the MockTransactionService.
func (m *MockTransactionService) GetTransactionByTxID(txid string) (*WalletExplorer.ITransactionData, error) {
	args := m.Called(txid)
	if result := args.Get(0); result != nil {
		return result.(*WalletExplorer.ITransactionData), args.Error(1)
	}
	return nil, args.Error(1)
}
