// Package types provides type definitions for user signup requests.
package types

// ITransactionXPUB represents the payload from external API response
type ITransactionXPUB struct {
	Found        bool              `json:"found"`
	GapLimit     int               `json:"gap_limit"`
	Transactions []XPUBTransaction `json:"txs"`
}

// XPUBTransaction represents the payload from external API response
type XPUBTransaction struct {
	TxID        string   `json:"txid"`
	BlockHeight int      `json:"block_height"`
	BlockPos    int      `json:"block_pos"`
	Time        int64    `json:"time"`
	BalanceDiff float64  `json:"balance_diff"`
	WalletIDs   []string `json:"wallet_ids"`
	Balance     float64  `json:"balance"`
}
