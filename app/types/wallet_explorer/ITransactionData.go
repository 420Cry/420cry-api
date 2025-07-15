// Package types provides type definitions for user signup requests.
package types

// ITransactionData represents the payload from external API response
type ITransactionData struct {
	Hash        string   `json:"hash"`
	Ver         int      `json:"ver"`
	VinSz       int      `json:"vin_sz"`
	VoutSz      int      `json:"vout_sz"`
	LockTime    any      `json:"lock_time"` // can be int or string
	Size        int      `json:"size"`
	Weight      *int     `json:"weight,omitempty"` // optional
	Fee         *int64   `json:"fee,omitempty"`    // optional
	RelayedBy   string   `json:"relayed_by"`
	BlockHeight int      `json:"block_height"`
	BlockIndex  *int64   `json:"block_index,omitempty"`
	TxIndex     any      `json:"tx_index"` // can be string or number
	DoubleSpend *bool    `json:"double_spend,omitempty"`
	Time        int64    `json:"time"`
	Inputs      []Input  `json:"inputs"`
	Out         []Output `json:"out"`
}

// Input represents the payload from external API response
type Input struct {
	Sequence *int64   `json:"sequence,omitempty"`
	Witness  *string  `json:"witness,omitempty"`
	Script   string   `json:"script"`
	Index    *int     `json:"index,omitempty"`
	PrevOut  *PrevOut `json:"prev_out,omitempty"`
}

// PrevOut represents the payload from external API response
type PrevOut struct {
	Type              *int       `json:"type,omitempty"`
	Spent             *bool      `json:"spent,omitempty"`
	Value             any        `json:"value"` // sometimes int, sometimes string
	SpendingOutpoints []Outpoint `json:"spending_outpoints"`
	N                 any        `json:"n"`        // sometimes string or int
	TxIndex           any        `json:"tx_index"` // can vary
	Script            string     `json:"script"`
}

// Output represents the payload from external API response
type Output struct {
	Type              *int       `json:"type,omitempty"`
	Spent             *bool      `json:"spent,omitempty"`
	Value             any        `json:"value"` // sometimes int, sometimes string
	SpendingOutpoints []Outpoint `json:"spending_outpoints,omitempty"`
	N                 *int       `json:"n,omitempty"`
	TxIndex           any        `json:"tx_index"` // can vary
	Script            string     `json:"script"`
	Addr              *string    `json:"addr,omitempty"`
}

// Outpoint represents the payload from external API response
type Outpoint struct {
	TxIndex any `json:"tx_index"`
	N       int `json:"n"`
}
