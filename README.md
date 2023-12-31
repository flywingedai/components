# component
Golang code testing and genreation package. 

## Test

## Generation

### Installation

For generation tool:
```sh
go install github.com/flywingedai/components@latest
```

For use of the tests package
```sh
go get github.com/flywingedai/components/tests
```

### Struct Registration
```golang
type component struct {
    /* 
        generate::components
        interfaceName::$STRING_VALUE
        interfaceFolder::$STRING_VALUE
        interfaceFile::$STRING_VALUE
        mockFolder::$STRING_VALUE
        mockFile::$STRING_VALUE
        skipTestFile::$BOOL_VALUE
        blackbox::$BOOL_VALUE
        expecters::$$STRING_VALUE
        config::$STRING_VALUE
    */
    subComponent interfaces.SubComponent `pkg:"$mockPackage" new:"$newMock" type:"mockType"`
}
```

Params values that have the "pkg" tag must be the same as the component values
with the only difference being whether or not they are exposed variables.

#### Component Options
To define options for component generation, just them in the above format in a
comment.

In order for the components command to register a struct as something that
should be generated, you need to add the `generate::components` string somewhere
withing that structs body. The above shows the standard way of doing it, but as
long as that string appears somewhere, it will be recognized.

The components command recognizes each of the above option before the "::" as
valid. The comman will look for all instances of "::" in the struct body, and
extract the option key and the option value. If any non-recognized options are
found, the command will fail. Below is a brief description of all the options:

- **generate:** Requires value to be "components". Registers the struct as a
component that should be generated by the command.
- **interfaceName:** [Optional] Name of the generated interface. If ignored,
will set the generated interface name to the name of the component with a
capital letter to export it.
- **interfaceFolder:** [Optional] The folder (and package) the interface will
live in. It is assumed the package name is the base of the provided directory
path. If ignored, will create the interface in the same folder and package the
struct is defined in.
- **interfaceFile:** [Optional] The base name of the file to place the generated
interface into. Will default to `$interfaceName.go` if nothing is provided AND the
interface folder is different than the package folder. If the interface folder
is the package folder (which is what `interfaceFolder` defaults to), the
`interfaceFile` will be set to the file the struct was defined in. (Will append
the interface to the end of the file along with the generated "New" function.)
- **mockFolder:** [Optional] The folder (and package) the mocks generated by
the mockery command will live in. It is assumed the package name is the base of
the provided directory path. Defaults to
`interfaceFolder/{{interfaceFolderBase}}_mocks`. If you pass in a value of
`__package__`, it will default to `packageFolder/{{packageFolderBase}}_mocks`
instead.
- **mockFile:** [Optional] The name of the generated mock file for this struct.
Defaults to `{{interfaceName}}.go` with interfaceName having a lowercase first
letter.
- **skipTestFile:** [Optional] Whether or not to generate a test file. This test
file will have mock definitions to use for creating tests for this specific
component. Defaults to `false`. Set to true by `skipTestFile::true`. If true,
blackbox and expecters options don't have any effect as those are options
specific to the test file.
- **blackbox:** [Optional] Whether or not the generated test file will be placed
in the same package as the struct or not. If enabled, this facilitates 
"blackbox" testing where the test files are all part of a new `{{package}}_test`
package which does not have access to private values within the package. To
enable, `blackbox::true`
- **expecters:** [Optional] Whether or not the generate test file will have
expecter bindings automatically generated for the given mock fields. Each mock
that should be included should be separated by a ",". To ignore all values, set
expecters = "-".
- **config:** [Optional] The mockery config file to use for this component
generation. The path should be relative to the place you execute the components
command or be absolute. Some options do not work because the components package
needs them set a specific way. `with-expecter` will always be true, and
`filename` is automatically inherited based on the `mockFile` option.

#### Tagging
To specify that a field of the component is a mock and should be treated as
such, you will need to set the `pkg`, `new`, and `type` tags.

- **pkg:** The name of the mock package
- **new:** The name of the function which creates a new mocked version of this type
- **type:** The name of the type as referred to by the mocks

IThe above tags can be automatically inferred by using the `pkg:"-"` tag. This
will tell the components command that all the values are standard. The command
will automatically generated standard mocks, so if you are mocking something
that was made with the components command, this will work. Standard values are
below:

- **pkg:** The name of the existing package + "_mocks"
- **new:** "New" + existing type
- **type:** The same as the defined type

example:
```golang
type component struct {
    subComponent  subcomponent.SubComponent `pkg:"-"`
}
```

is cconverted to

```golang
type mocks struct {
    subComponent  *subcomponent_mocks.SubComponent
}

func buildMocks(t *testing.T) (maincomponent.MainComponent, *mocks) {
	params := initParams()

	params.SubComponent = subcomponent_mocks.NewSubComponent(t)

	return maincomponent.New(params), convert(params)
}
```

### Usage

Simply run the command below on a specified directory, and
`components $DIRECTORY`
