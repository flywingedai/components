package generate

import (
	"os"

	"github.com/flywingedai/components/generate/componentparser"
	"github.com/spf13/cobra"
)

func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func NewRootCmd() *cobra.Command {

	// Create the new base cobra command
	baseCommand := &cobra.Command{}

	// Set the context
	baseCommand.Use = "components"
	baseCommand.Short = "Generate mock objects for your Golang interfaces using mockery, and then support component based testing and structure"
	baseCommand.Example = "components $DIRECTORY"

	// Create the main run command
	baseCommand.Run = func(cmd *cobra.Command, args []string) {

		// Create the parser according to the directory specified
		p := componentparser.New(cmd)

		if len(args) != 1 {
			panic("Invalid usage - " + cmd.Example)
		}
		p.Args.Directory = args[0]

		// Parse all files in the path specified
		p.Parse()

		// Create the interface files for each of the structs that were found
		for _, structData := range p.Structs {
			generateInterface(structData)
		}

		// Run the mockery command for each of the structs that were found
		for _, structData := range p.Structs {
			callMockery(structData)
		}

		/*
			Extend each of the mock files with additional functionality for
			components. Simply adds some extra data to the end of each generated
			mock file
		*/
		for _, structData := range p.Structs {
			extendMocks(structData)
		}

		// Create test files for each of the structs
		for _, structData := range p.Structs {
			generateTest(structData)
		}

	}

	return baseCommand
}
