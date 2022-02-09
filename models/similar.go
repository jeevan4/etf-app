package models

import (
	"encoding/json"
	"os"
)

type Similars []Similar

type Similar struct {
	Ticker string `json:"ticker"`
	Name   string `json:"name"`
}

func NewSimilar(ticker string, name string, weight float32) *Similar {
	return &Similar{
		Ticker: ticker,
		Name:   name,
	}
}

func (s *Similar) ToJson() error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "\t")
	encoder.Encode(s)
	return nil
}

func (s *Similar) FromJson(b []byte) error {
	err := json.Unmarshal(b, s)
	if err != nil {
		return err
	}
	return nil
}
