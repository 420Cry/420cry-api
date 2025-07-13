// Package types provides type definitions for user signup requests.
package types

// ITransactionData represents the payload from external API response
type ITransactionData struct {
	Found          bool   `json:"found"`
	Label          string `json:"label"`
	Txid           string `json:"txid"`
	IsCoinbase     bool   `json:"is_coinbase"`
	WalletID       string `json:"wallet_id"`
	BlockHeight    int    `json:"block_height"`
	BlockPos       int    `json:"block_pos"`
	Time           int64  `json:"time"`
	Size           int    `json:"size"`
	In             any    `json:"in"`
	Out            any    `json:"out"`
	UpdatedToBlock int    `json:"updated_to_block"`
}
