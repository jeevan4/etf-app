package cli

import (
	"etf-app/schemas"
	"fmt"
)

type EtfCommands struct {
	Name     string
	Commands map[string]func(args []string) error
}

func (e *EtfCommands) Help() {
	for key, _ := range e.Commands {
		fmt.Println("etf-app ", key, "--help")
	}
}

func (e *EtfCommands) Query(args []string) error {
	fmt.Println("Command from Query")
	queryCmd := NewQueryCmd()
	queryCmd.Run(args)
	// print the fetched results
	// queryCmd.ToJson(queryCmd.allData)
	// queryCmd.ToJson(queryCmd.etfDetails)
	return nil
}

func (e *EtfCommands) List(args []string) error {
	fmt.Println("Command from List")
	listCmd := NewListCommand()
	listCmd.Run(args)
	return nil
}

func NewEtfCommand(name string) *EtfCommands {
	// prepare database
	err := schemas.InitDatabase()
	if err != nil {
		panic(err)
	}
	etfCommand := &EtfCommands{
		Name: name,
	}
	etfCommand.Commands = map[string]func(args []string) error{
		"query": etfCommand.Query,
		"list":  etfCommand.List,
	}
	return etfCommand
}
