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
	fmt.Println("This is help!")
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
	}
	return etfCommand
}
