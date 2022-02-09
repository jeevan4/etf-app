package cli

import (
	"etf-app/models"
	"flag"
	"fmt"
	"net/http"
	"strings"
)

const url string = "https://70a77bonik.execute-api.us-east-1.amazonaws.com/live/similar?ticker="

type etfList []string

func (e *etfList) String() string {
	s := fmt.Sprintf("%s", *e)
	return s
}

func (e *etfList) Set(value string) error {
	*e = append(*e, strings.ToLower(value))
	return nil
}

type QueryCmd struct {
	Holdings bool
	Similar  bool
	etf      etfList
	FlagSet  *flag.FlagSet
}

func NewQueryCmd() *QueryCmd {
	queryCmd := &QueryCmd{}
	queryCmd.FlagSet = flag.NewFlagSet("Query ETF Api", flag.ExitOnError)
	queryCmd.FlagSet.BoolVar(&queryCmd.Holdings, "holdings", false, "Get top holdings for the provided etf")
	queryCmd.FlagSet.BoolVar(&queryCmd.Similar, "similar", false, "Get sililar etf for the provided etf")
	queryCmd.FlagSet.Var(&queryCmd.etf, "etf", "The etf name to query the data for")
	return queryCmd
}

func (q *QueryCmd) Run(args []string) error {
	fmt.Println("running the query command")
	err := q.FlagSet.Parse(args)
	if err != nil {
		return err
	}
	fmt.Printf("Holdings: %v, Similar: %v, Etfs %s \n", q.Holdings, q.Similar, q.etf)
	etfUrl := fmt.Sprintf("%s%s", url, q.etf[0])
	fmt.Println(etfUrl)
	resp, err := http.Get(etfUrl)
	if err != nil {
		return err
	}
	allData := models.AllData{}
	// d, _ := io.ReadAll(resp.Body)
	// fmt.Println(string(d))
	err = allData.FromJson(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Printf("%+v", allData.Topten[0])
	allData.ToJson()
	// hold := models.NewHolding("AAPL", "Apple Inc,", 3.45)
	// hold.ToJson()
	// data := []byte(`{"ticker": "msft", "weightxx":"1.34"}`)
	// hold1 := models.Holding{}
	// err = hold1.FromJson(data)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Printf("%+v\n", hold1)
	return nil
}
