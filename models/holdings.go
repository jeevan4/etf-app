package models

import (
	"encoding/json"
	"os"
)

type Holdings []Holding

type Holding struct {
	Ticker string `json:"ticker"`
	Name   string `json:"name"`
	Weight string `json:"weight"`
}

func NewHolding(ticker string, name string, weight string) *Holding {
	return &Holding{
		Ticker: ticker,
		Name:   name,
		Weight: weight,
	}
}

func (h *Holding) ToJson() error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "\t")
	encoder.Encode(h)
	return nil
}

func (h *Holding) FromJson(b []byte) error {
	err := json.Unmarshal(b, h)
	if err != nil {
		return err
	}
	return nil
}
