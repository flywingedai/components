package componentparser

import (
	"go/ast"
	"io/fs"
	"path"
	"path/filepath"
	"strings"

	"github.com/flywingedai/components/generate/helpers"
	"github.com/spf13/cobra"
	"golang.org/x/tools/go/packages"
)

/*
A parser looks through a directory and all it's children. It creates a group of
StructData objects representing all the structs that should be generated.
*/
type Parser struct {
	Args ParserArgs

	Structs        map[string]*StructData
	PackageImports map[string]map[string]bool // List of all the imports needed for a package

	// Values that get updated as the parser is walking
	File          string     // Which file is currently being parsed
	FileString    FileString // The extracted file string corresponding to .File
	PackageFolder string     // Which package fodler is currently being parsed
	PackageName   string     // Which package is currently being parsed
	ModulePath    string     // Go Package Path

}

type ParserArgs struct {
	/*
		The directory to run the parser in. All child folders will be walked
		through automatically.
	*/
	Directory string

	/*
		The default config to use for mockery generation commands. Will be
		automatically passed into all child struct generate commands.
	*/
	Config string

	// Regex to match against struct tags
	Match string
}

func New(cmd *cobra.Command) *Parser {
	return &Parser{
		Args:           ParserArgs{},
		Structs:        map[string]*StructData{},
		PackageImports: map[string]map[string]bool{},
	}
}

/*
Parse everything in the specified directory according to the args.
*/
func (p *Parser) Parse() {
	var err error

	// Recursively call ParseDir on each of the dirs
	filepath.WalkDir(p.Args.Directory, func(dir string, d fs.DirEntry, err error) error {
		if err != nil {
			panic(err)
		}

		if d.IsDir() {
			p.ParseDir(dir)
		}
		return nil
	})

	// Clean up the data.
	for key, structData := range p.Structs {
		if !structData.Options.Generate {
			delete(p.Structs, key)
			continue
		}

		if structData.ConvertFunction == "" {
			panic("struct " + structData.Name + " in " + structData.StructFile + " does not have a *Params.Convert() *" + structData.Name + " function!")
		}

		/*
			Update the parameters that had default values. Panic if any invalid
			args passed in by the user.
		*/

		// Interface management
		if structData.Options.InterfaceName == "" {
			structData.Options.InterfaceName = helpers.ToTitle(structData.Name)
		}

		if structData.Options.InterfaceFolder == "" {
			structData.Options.InterfaceFolder = structData.PackageFolder
			structData.Options.InterfacePackage = structData.PackageName
		} else {
			structData.Options.InterfaceFolder, err = filepath.Abs(structData.Options.InterfaceFolder)
			if err != nil {
				panic(err)
			}
			structData.Options.InterfacePackage = path.Base(structData.Options.InterfaceFolder)
		}
		structData.Options.InterfaceFolder, err = filepath.Abs(structData.Options.InterfaceFolder)
		if err != nil {
			panic(err)
		}

		if structData.Options.InterfaceFile == "" {
			if structData.Options.InterfaceFolder != structData.PackageFolder {
				structData.Options.InterfaceFile = path.Join(structData.Options.InterfaceFolder, "interface.go")
			} else {
				structData.Options.InterfaceFile = structData.StructFile
			}
		} else {
			structData.Options.InterfaceFile = path.Join(structData.Options.InterfaceFolder, structData.Options.InterfaceFile)
		}

		// Mock management
		if structData.Options.MockFolder == "" {
			structData.Options.MockFolder = path.Join(structData.Options.InterfaceFolder, path.Base(structData.Options.InterfaceFolder)+"_mocks")
		} else {
			structData.Options.MockFolder, err = filepath.Abs(structData.Options.MockFolder)
			if err != nil {
				panic(err)
			}
		}
		structData.Options.MockPackage = path.Base(structData.Options.MockFolder)

		if structData.Options.MockFile == "" {
			structData.Options.MockFile = helpers.ToCamel(structData.Options.InterfaceName) + ".go"
		}

		// The struct must not be exported if the interface name is also the same
		if structData.Options.InterfaceName == structData.Name {
			panic("struct " + structData.Name + " in " + structData.StructFile + " has the same interface name!")
		}

	}
}

/*
Parse method that populates all the needed data for a specific package/dir.
*/
func (p *Parser) ParseDir(dir string) {

	/*
		List of packages present in this directory. By default, we don't want to
		include any _test packages. There should only be one in each directory.
	*/
	pkgs, err := packages.Load(&packages.Config{
		Dir:  dir,
		Mode: packages.NeedFiles + packages.NeedImports + packages.NeedName,
	})
	if err != nil {
		panic(err)
	}

	for _, pkg := range pkgs {
		if strings.Contains(pkg.Name, "_test") {
			continue
		}

		p.PackageName = pkg.Name
		p.PackageFolder, err = filepath.Abs(dir)
		if err != nil {
			panic(err)
		}

		/*
			Once we've found the root package, loop through all the files that
			are a part of that package and parse them individually.
		*/
		for _, fileName := range pkg.GoFiles {
			p.File, err = filepath.Abs(fileName)
			if err != nil {
				panic(err)
			}

			p.ParseFile()
		}

	}

}

/*
Parser method that reads an individual file and adjusts the data saved in
the Structs field as data is read in.
*/
func (p *Parser) ParseFile() {

	// Extract the *ast.File for the file and the full file contents
	file, fileString := readFile(p.File)
	p.FileString = fileString

	/*
		For every *ast.Node, we parse and accumulate relevant information into
		the p.Structs map. Any extra values will be removed if they don't match
		later. For now we just want to get all the data into memory.
	*/
	ast.Inspect(file, func(n ast.Node) bool {

		switch node := n.(type) {

		// Used to capture type declarations
		case *ast.GenDecl:
			p.ParseGenDecl(node)

		case *ast.FuncDecl:
			p.ParseFuncDecl(node)
		}

		return true

	})

}
