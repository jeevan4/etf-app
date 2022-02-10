package cli

import (
	"encoding/json"
	"etf-app/models"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const url string = "https://70a77bonik.execute-api.us-east-1.amazonaws.com/live/similar?ticker="
const etfDataUrl string = "https://70a77bonik.execute-api.us-east-1.amazonaws.com/live/other?ticker="

type etfList []string

func (e *etfList) String() string {
	s := fmt.Sprintf("%s", *e)
	return s
}

func (e *etfList) Set(value string) error {
	*e = append(*e, strings.ToLower(value))
	return nil
}

func (e *etfList) FromJson(b io.Reader) error {
	decoder := json.NewDecoder(b)
	err := decoder.Decode(e)
	if err != nil {
		return err
	}
	return nil
}

type QueryCmd struct {
	Holdings   bool
	Similar    bool
	etf        etfList
	FlagSet    *flag.FlagSet
	allData    map[string]models.AllData
	etfDetails map[string]etfList
}

func NewQueryCmd() *QueryCmd {
	queryCmd := &QueryCmd{}
	queryCmd.etfDetails = map[string]etfList{}
	queryCmd.allData = map[string]models.AllData{}
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

	// create in and out channel to fetch positions/similars
	in := make(chan string)
	out := make(chan map[string]models.AllData)

	// create in and out channel to fetch data about ETF itself
	inEtf := make(chan string)
	outEtf := make(chan map[string]etfList)

	// invoke goroutines for the number of etfs given
	for i := 0; i < len(q.etf); i++ {
		go q.fetchData(in, out, i)
		go q.fetchEtfData(inEtf, outEtf, i)
	}

	for _, val := range q.etf {
		in <- val
		inEtf <- val
	}
	close(in)
	close(inEtf)

	// get etf data from channels
	for i := 0; i < len(q.etf); i++ {
		data := <-outEtf
		// q.etfDetails = append(q.etfDetails, data)
		for key, val := range data {
			q.etfDetails[key] = val
		}
	}

	// get holdings/similars data from channels
	for i := 0; i < len(q.etf); i++ {
		data := <-out
		// q.allData = append(q.allData, data)
		for key, val := range data {
			q.allData[key] = val
		}

	}
	// for key, val := range q.allData {
	// 	fmt.Println("Data for :", key, q.etfDetails[key][2])
	// 	for id, holding := range val.Topten {
	// 		fmt.Printf("%d Name: %s, Tikr: %s, Weight %s\n", id, holding.Name, holding.Ticker, holding.Weight)
	// 	}

	// }
	// for key, val := range q.etfDetails {
	// 	fmt.Println(key, val)
	// }
	// fmt.Println(q.allData, q.etfDetails)
	return nil
}

func (q *QueryCmd) fetchData(in <-chan string, out chan<- map[string]models.AllData, id int) error {
	for tickr := range in {
		start := time.Now()
		etfUrl := fmt.Sprint(url, tickr)
		resp, err := http.Get(etfUrl)
		if err != nil {
			return err
		}
		allData := models.AllData{}
		allData.FromJson(resp.Body)
		fmt.Println("data ready for:", tickr, "Id:", id, "elapsed:", time.Now().Sub(start))
		out <- map[string]models.AllData{tickr: allData}
	}
	return nil
}

func (q *QueryCmd) fetchEtfData(in <-chan string, out chan<- map[string]etfList, id int) error {
	for tickr := range in {
		start := time.Now()
		etfUrl := fmt.Sprint(etfDataUrl, tickr)
		resp, err := http.Get(etfUrl)
		if err != nil {
			return err
		}
		etfData := etfList{}
		etfData.FromJson(resp.Body)
		fmt.Println(etfData, "elapsed", time.Now().Sub(start))
		out <- map[string]etfList{tickr: etfData}
	}
	return nil
}

func (q *QueryCmd) ToJson(v interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(v)
	if err != nil {
		return err
	}
	return nil
}
