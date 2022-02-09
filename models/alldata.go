package models

import (
	"encoding/json"
	"io"
	"os"
)

type AllData struct {
	DateOpen   string   `json:"dte_topten"`
	Topten     Holdings `json:"topten"`
	TopSimilar Similars `json:"similar"`
}

func NewAllData(DateOpen string, topten Holdings, topsimilar Similars) *AllData {
	return &AllData{
		DateOpen:   DateOpen,
		Topten:     topten,
		TopSimilar: topsimilar,
	}
}

func (a *AllData) ToJson() error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "\t")
	encoder.Encode(a)
	return nil
}

func (a *AllData) FromJson(b io.Reader) error {
	decoder := json.NewDecoder(b)
	err := decoder.Decode(a)
	if err != nil {
		return err
	}
	// return nil
	// err := json.Unmarshal(b, a)
	// if err != nil {
	// 	return err
	// }
	return nil
}
