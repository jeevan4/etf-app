package cli

import (
	"encoding/json"
	"etf-app/models"
	"etf-app/schemas"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

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
	etf        etfList
	FlagSet    *flag.FlagSet
	allData    map[string]models.AllData
	etfDetails map[string]etfList
	Holdings   bool
	Similar    bool
	Store      bool
}

func NewQueryCmd() *QueryCmd {
	queryCmd := &QueryCmd{}
	queryCmd.etfDetails = map[string]etfList{}
	queryCmd.allData = map[string]models.AllData{}
	queryCmd.FlagSet = flag.NewFlagSet("Query ETF Api", flag.ExitOnError)
	queryCmd.FlagSet.BoolVar(&queryCmd.Holdings, "holdings", false, "Get top holdings for the provided etf")
	queryCmd.FlagSet.BoolVar(&queryCmd.Similar, "similar", false, "Get sililar etf for the provided etf")
	queryCmd.FlagSet.BoolVar(&queryCmd.Store, "store-to-db", false, "Save details of requsted etfs to local storage")
	queryCmd.FlagSet.Var(&queryCmd.etf, "etf", "The etf name to query the data for")
	return queryCmd
}

func (q *QueryCmd) Run(args []string) error {
	fmt.Println("running the query command")
	err := q.FlagSet.Parse(args)
	if err != nil {
		return err
	}
	fmt.Printf("Holdings: %v, Similar: %v, Etfs %s, Storage: %v \n", q.Holdings, q.Similar, q.etf, q.Store)

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
	if q.Store {
		return q.addToDb()

	}
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

func (q *QueryCmd) addToDb() error {
	db, err := schemas.NewDbConnection()
	if err != nil {
		return err
	}
	insert_quert := "INSERT OR IGNORE INTO stocks (ticker, description, etf, expense_ratio) values"
	holding_quert := "INSERT OR IGNORE INTO holdings (etfname, ticker, weight, refresh_date) values"

	var insert_values []interface{}
	var holding_insert_values []interface{}
	for etf, detail := range q.etfDetails {
		insert_quert += "(?,?,?,?),"
		insert_values = append(insert_values, strings.ToUpper(etf), detail[2], true, detail[1])
	}
	for etf, alldata := range q.allData {
		for _, holdings := range alldata.Topten {
			insert_quert += "(?,?,?,?),"
			holding_quert += "(?,?,?,?),"
			insert_values = append(insert_values, holdings.Ticker, holdings.Name, false, 0)
			topten_date, _ := time.Parse("01/02/2006", alldata.DateOpen)

			holding_insert_values = append(holding_insert_values, strings.ToUpper(etf), holdings.Ticker, holdings.Weight, topten_date.Format("2006/01/02"))
		}
	}
	// fmt.Println(insert_quert[0 : len(insert_quert)-1])
	// fmt.Println(insert_values...)
	insert_query, err := db.Prepare(insert_quert[0 : len(insert_quert)-1])
	holding_query, err := db.Prepare(holding_quert[0 : len(holding_quert)-1])
	if err != nil {
		fmt.Println(err)
	}
	res, err := insert_query.Exec(insert_values...)
	res, err = holding_query.Exec(holding_insert_values...)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res.RowsAffected())
	// insert_stock, err := db.Prepare("INSERT INTO stocks (ticker, description, etf, expense_ratio) ")
	return nil
}
