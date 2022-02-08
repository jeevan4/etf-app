package cli

import (
	"flag"
	"fmt"
	"strings"
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
	return nil
}
