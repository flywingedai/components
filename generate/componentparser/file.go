package componentparser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

// Helper wrapper for strings with some additional helper methods
type FileString string

/*
Grab the *ast.File and FileString representation of a file
*/
func readFile(fileName string) (*ast.File, FileString) {

	// Create a new fileset to use with the built-in parser package
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	// file, err := parser.ParseFile(fset, fileName, nil, parser.Trace)
	if err != nil {
		panic(err)
	}

	// Extract the full file from the
	fileData, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	fileString := FileString(fileData)

	return file, fileString

}

/*
Extract the node string from a file. Convert all the types accordingly and
properly handles indexing of the node position.
*/
func (f FileString) Extract(node ast.Node) string {
	return string(f[node.Pos()-1 : node.End()-1])
}

/*
Grab the children nodes of the specified type
*/
func FindChildNodes[T any](node ast.Node) []T {

	nodes := []T{}
	ast.Inspect(node, func(n ast.Node) bool {
		casted, ok := n.(T)
		if ok {
			nodes = append(nodes, casted)
		}
		return true
	})

	return nodes
}

// Find the first child node of the specified type
func FindChildNode[T any](node ast.Node) T {

	nodes := FindChildNodes[T](node)
	if len(nodes) == 0 {
		panic("no nodes of specified type found")
	}
	return nodes[0]
}
