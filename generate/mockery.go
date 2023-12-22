package generate

import (
	"os"
	"os/exec"

	"github.com/flywingedai/components/generate/componentparser"
)

func callMockery(structData *componentparser.StructData) {
	/*
		Switch up the base directory for the interfaces to match the package
		folder for the struct in question.
	*/
	originalDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	err = os.Chdir(structData.Options.InterfaceFolder)
	if err != nil {
		panic(err)
	}

	// Run the tailored mockery command for that struct
	mockeryCommand := exec.Command("mockery",
		"--name", structData.Options.InterfaceName,
		"--filename", structData.Options.MockFile,
		"--output", structData.Options.MockFolder,
		"--outpkg", structData.Options.MockPackage,
		"--config", structData.Options.Config,
		"--with-expecter",
	)

	// Set output so the mockery output is viewable
	mockeryCommand.Stderr = os.Stderr
	mockeryCommand.Stdout = os.Stdout
	err = mockeryCommand.Run()
	if err != nil {
		panic(err)
	}

	// Switch back to the original directory after the mockery call is made
	err = os.Chdir(originalDir)
	if err != nil {
		panic(err)
	}
}
