package componentparser

import (
	"fmt"
	"go/ast"
	"strings"
)

/*
A field is a name and a type. If in a struct, it can also have a tag. If a tag
is in this format, `component:"$TAG"`, it will be in the TAG field.
*/
type Field struct {
	Name string
	Type string

	MockPkg  string
	MockNew  string
	MockType string
}

type Fields []Field

/*
Name: The name of the method.
Recv: The reciever for the method.
Args: List of all the args for the method.
Returns: List of al lthe returns
*/
type MethodData struct {
	Name    string
	Recv    Field
	Args    Fields
	Returns Fields
}

/////////////////////
// REPRESENTATIONS //
/////////////////////

// Return args in format $Arg1_Name $Arg1_Type, ..., $ArgN_Name $ArgN_Type
func (fields Fields) AsArgs(includeParenthesis bool) string {
	returnString := ""
	for i, field := range fields {
		returnString += field.Name + " " + field.Type
		if i != len(fields)-1 {
			returnString += ", "
		}
	}
	if includeParenthesis && len(fields) > 1 {
		returnString = "(" + returnString + ")"
	}
	return returnString
}

// Return args in format $Arg1_Name, ..., $ArgN_Name
func (fields Fields) AsParams() string {
	returnString := ""
	for i, field := range fields {
		returnString += field.Name
		if i != len(fields)-1 {
			returnString += ", "
		}
	}
	return returnString
}

// Return args in format $Arg1_Type, ..., $ArgN_Type
func (fields Fields) AsTypes(includeParenthesis bool) string {
	returnString := ""
	for i, field := range fields {
		returnString += field.Type
		if i != len(fields)-1 {
			returnString += ", "
		}
	}
	if includeParenthesis && len(fields) > 1 {
		returnString = "(" + returnString + ")"
	}
	return returnString
}

// Return args in format $Arg1_Name interface{}, ..., $ArgN_Name interface{}
func (fields Fields) AsInterface(includeParenthesis bool) string {
	returnString := ""
	for i, field := range fields {
		returnString += field.Name + " interface{}"
		if i != len(fields)-1 {
			returnString += ", "
		}
	}
	if includeParenthesis && len(fields) > 1 {
		returnString = "(" + returnString + ")"
	}
	return returnString
}

// Get generic representations
func (fields Fields) Generic(append bool) (string, string) {
	generic := ""
	genericLong := ""
	if len(fields) > 0 {
		if append {
			generic += fmt.Sprintf(", %s", fields.AsParams())
			genericLong += fmt.Sprintf(", %s", fields.AsArgs(false))
		} else {
			generic += fmt.Sprintf("[%s]", fields.AsParams())
			genericLong += fmt.Sprintf("[%s]", fields.AsArgs(false))
		}

	}
	return generic, genericLong
}

func ReplaceScopedNames(s, packageName string, scopedNames map[string]bool) string {
	for name := range scopedNames {
		s = strings.ReplaceAll(s, " "+name, " "+packageName+"."+name)
	}
	return s
}

func (fields Fields) AddPointers(s string, asInterface bool) string {
	if s == "" {
		return s
	}

	prefix := "*"
	suffix := ""
	if asInterface {
		prefix = "tests.RemoveInterfacePointer[{{REPLACE}}]("
		suffix = ")"
	}

	splits := strings.Split(s, ", ")
	for i, split := range splits {

		fieldPrefix := prefix
		if asInterface {
			fieldPrefix = strings.ReplaceAll(prefix, "{{REPLACE}}", fields[i].Type)
		}

		if strings.Contains(split, " ") {
			splits[i] = strings.ReplaceAll(split, " ", " "+fieldPrefix) + suffix
		} else {
			splits[i] = fieldPrefix + split + suffix
		}
	}
	return strings.Join(splits, ", ")
}

/////////////////
// CONSTRUCTOR //
/////////////////

// Convert a []*ast.Field -> Fields
func ConvertASTFieldList(fileString FileString, node *ast.FieldList) Fields {

	// Check for definitions that don't inherently have a type
	if len(node.List) != node.NumFields() {

		// Extract the string, accounting for parenthesis
		typeString := fileString.Extract(node)
		if typeString[0] == '(' {
			typeString = typeString[1 : len(typeString)-1]
			return ConvertTypeString(typeString)
		}

	}

	fields := Fields{}

	for i, fieldNode := range node.List {

		field := Field{}

		/*
			First, extract any tags that may be present on this field. The
			components package cares about a "pkg" and "new" tag. These
			correspond to the mock package name and the new function for
			that package.
		*/
		fieldString := fileString.Extract(fieldNode)

		for _, tag := range []string{"pkg", "new", "type"} {
			tagID := fmt.Sprintf("%s:\"", tag)
			index := strings.Index(fieldString, tagID)
			if index >= 0 {
				index += len(tagID)
				endIndex := strings.Index(fieldString[index:], `"`)
				if endIndex >= 0 {

					if tag == "pkg" {
						field.MockPkg = fieldString[index : index+endIndex]
					} else if tag == "new" {
						field.MockNew = fieldString[index : index+endIndex]
					} else if tag == "type" {
						field.MockType = fieldString[index : index+endIndex]
					}

				}
			}
		}

		// Extract the name and type from the
		if len(fieldNode.Names) > 0 {
			field.Name = fieldNode.Names[0].Name
		} else {
			field.Name = fmt.Sprintf("_a%d", i)
		}
		field.Type = fileString.Extract(fieldNode.Type)

		// Fix the mock tags if exists
		if field.MockPkg == "-" {
			split := strings.Split(field.Type, ".")
			if len(split) < 2 {
				panic("bad type for auto mock inference " + field.Name + " ")
			}
			field.MockPkg = split[0] + "_mocks"
			field.MockType = strings.Join(split[1:], ".")
			field.MockNew = "New" + field.MockType
		}

		fields = append(fields, field)
	}

	return fields
}

func ConvertTypeString(typeString string) Fields {
	fields := Fields{}

	// Split by the commas
	fieldStrings := strings.Split(strings.ReplaceAll(typeString, "\n", ""), ",")

	isReturn := true
	for _, fieldString := range fieldStrings {
		for fieldString[0] == ' ' || fieldString[0] == '\t' {
			fieldString = fieldString[1:]
		}

		split := strings.Split(fieldString, " ")
		if len(split) == 2 {
			isReturn = false
		}
	}

	// Check that each field exists in the fields
	for i := len(fieldStrings) - 1; i >= 0; i-- {
		fieldString := fieldStrings[i]
		for fieldString[0] == ' ' || fieldString[0] == '\t' {
			fieldString = fieldString[1:]
		}

		split := strings.Split(fieldString, " ")
		fieldName := split[0]

		if isReturn {
			fields = append(Fields{Field{Name: split[0], Type: split[0]}}, fields...)
		} else {
			if len(split) == 1 {
				fields = append(Fields{Field{Name: fieldName, Type: fields[0].Type}}, fields...)
			} else {
				fields = append(Fields{Field{Name: fieldName, Type: split[1]}}, fields...)
			}
		}
	}

	return fields
}

// Get type without point
func CleanType(t string) string {

	index := strings.Index(t, "[")
	if index != -1 {
		t = t[:index]
	}

	if len(t) <= 1 {
		return t
	}

	if t[0:1] == "*" {
		return t[1:]
	} else {
		return t
	}
}
