package cli

import (
	"etf-app/schemas"
	"flag"
	"fmt"
)

// type etfList []string

// func (e *etfList) String() string {
// 	s := fmt.Sprintf("%s", *e)
// 	return s
// }

// func (e *etfList) Set(value string) error {
// 	*e = append(*e, strings.ToLower(value))
// 	return nil
// }

// func (e *etfList) FromJson(b io.Reader) error {
// 	decoder := json.NewDecoder(b)
// 	err := decoder.Decode(e)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

type ListCmd struct {
	FlagSet *flag.FlagSet
	etf     bool
	stocks  bool
}

func NewListCommand() *ListCmd {
	ListCmd := &ListCmd{}
	ListCmd.FlagSet = flag.NewFlagSet("List ETF/Stocks Currently Stored", flag.ExitOnError)
	ListCmd.FlagSet.BoolVar(&ListCmd.etf, "etf", false, "List all etfs currently stored")
	ListCmd.FlagSet.BoolVar(&ListCmd.stocks, "stock", false, "List all stocks currently stored")
	return ListCmd
}

func (l *ListCmd) Run(args []string) error {
	fmt.Println("running the List command")
	err := l.FlagSet.Parse(args)
	if err != nil {
		return err
	}
	fmt.Printf("Etf: %v, Stock: %v \n", l.etf, l.stocks)
	if l.etf {
		return l.getEtfsFromDb()
	}

	if l.stocks {
		return l.getStocksFromDb()
	}
	return nil
}

func (l *ListCmd) getEtfsFromDb() error {
	db, err := schemas.NewDbConnection()
	if err != nil {
		return err
	}
	etfQuery, err := db.Prepare("Select ticker, description, expense_ratio from stocks where etf = True;")
	if err != nil {
		return err
	}
	rows, _ := etfQuery.Query()
	columns, _ := rows.Columns()
	var result [][]interface{}
	for rows.Next() {
		rowdata := make([]interface{}, len(columns))
		for i := 0; i < len(columns); i++ {
			rowdata[i] = new(interface{})

		}
		rows.Scan(rowdata...)
		result = append(result, rowdata)
	}
	for _, val := range result {
		for _, field := range val {
			fmt.Printf("%v ", *field.(*interface{}))
		}
		fmt.Println()
	}
	return nil
}

func (l *ListCmd) getStocksFromDb() error {
	db, err := schemas.NewDbConnection()
	if err != nil {
		return err
	}
	stockQuery, err := db.Prepare("Select ticker, description from stocks where etf = False;")
	if err != nil {
		return err
	}
	rows, _ := stockQuery.Query()
	columns, _ := rows.Columns()
	var result [][]interface{}
	for rows.Next() {
		rowdata := make([]interface{}, len(columns))
		for i := 0; i < len(columns); i++ {
			rowdata[i] = new(interface{})

		}
		rows.Scan(rowdata...)
		result = append(result, rowdata)
	}
	for _, val := range result {
		for _, field := range val {
			fmt.Printf("%v ", *field.(*interface{}))
		}
		fmt.Println()
	}
	return nil
}
