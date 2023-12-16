package generate

import (
	"path"

	"github.com/flywingedai/components/generate/componentparser"
	"github.com/flywingedai/components/generate/helpers"
	"github.com/flywingedai/components/generate/templates"
)

func extendMocks(structData *componentparser.StructData) {

	// Start with the Expecterchain definition for this mock
	dataString := templates.BulkReplace(templates.ExpecterChain, map[string]string{"InterfaceName": structData.Options.InterfaceName})

	/*
		Loop through each of the methods and add each of their chain definitions
		onto the mock file.
	*/
	for _, method := range structData.Methods {

		/*
			Format the response types. As long as there is some response, we
			format correctly with an extra space at the beginning.
		*/
		responseTypes := method.Returns.AsTypes(true)
		if responseTypes != "" {
			responseTypes = " " + responseTypes
		}

		// Form all the replacement pairs for this method chain
		pairs := map[string]string{
			"InterfaceName": structData.Options.InterfaceName,
			"Method":        method.Name,
			"ArgsInterface": method.Args.AsInterface(false),
			"Args":          method.Args.AsArgs(false),
			"ArgsShort":     method.Args.AsParams(),

			"ReturnsArgs":  method.Returns.AsArgs(false),
			"ReturnsTypes": responseTypes,
			"ReturnsShort": method.Returns.AsParams(),
		}

		// Add this chain to the data string
		dataString += templates.BulkReplace(templates.Chain, pairs)
	}

	helpers.WriteToFile(path.Join(structData.Options.MockFolder, structData.Options.MockFile), dataString, structData.Imports, structData.Options.MockPackage)

}
