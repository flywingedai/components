package generate

import (
	"fmt"

	"github.com/flywingedai/components/generate/componentparser"
	"github.com/flywingedai/components/generate/helpers"
	"github.com/flywingedai/components/generate/templates"
)

// Generate an interface based on the struct passed int
func generateInterface(structData *componentparser.StructData) {

	/*
		// 	Create each of the methods captured during parsing. Aggregate them all
		// 	into a single string for readability.
	*/
	methodString := ""
	for j, m := range structData.Methods {
		methodString += fmt.Sprintf(templates.Method, m.Name, m.Args.AsArgs(false), m.Returns.AsTypes(true))
		if j != len(structData.Methods)-1 {
			methodString += "\n"
		}

	}

	interfaceString := fmt.Sprintf(templates.Interface, structData.Options.InterfaceName, methodString)
	helpers.WriteToFile(structData.Options.InterfaceFile, interfaceString, structData.Imports, structData.Options.InterfacePackage)

	interfaceName := structData.Options.InterfaceName
	if structData.Options.InterfacePackage != structData.PackageName {
		interfaceName = structData.Options.InterfacePackage + "." + interfaceName
	}
	newString := fmt.Sprintf(templates.New, interfaceName)
	helpers.WriteToFile(structData.StructFile, newString, structData.Imports, structData.PackageName)

}
