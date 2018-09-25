package graphql

import (
	"fmt"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
)

func document(doc *ast.Document) {
	fmt.Printf("DOCUMENT ==========>\n")
	fmt.Printf("  Kind: %#v\n", doc.Kind)
	fmt.Printf("  Definitions:\n")
	for _, def := range doc.Definitions {
		definition(def)
	}
}

func definition(node ast.Node) {
	kind := node.GetKind()

	switch kind {
	case "DirectiveDefinition":
		def := node.(*ast.DirectiveDefinition)
		directiveDefinition(def)
	case "EnumDefinition":
		def := node.(*ast.EnumDefinition)
		enumDefinition(def)
	case "EnumValueDefinition":
		def := node.(*ast.EnumValueDefinition)
		enumValueDefinition(def)
	case "FieldDefinition":
		def := node.(*ast.FieldDefinition)
		fieldDefinition(def)
	case "FragmentDefinition":
		def := node.(*ast.FragmentDefinition)
		fragmentDefinition(def)
	case "InputObjectDefinition":
		def := node.(*ast.InputObjectDefinition)
		inputObjectDefinition(def)
	case "InputValueDefinition":
		def := node.(*ast.InputValueDefinition)
		inputValueDefinition(def)
	case "InterfaceDefinition":
		def := node.(*ast.InterfaceDefinition)
		interfaceDefinition(def)
	case "ObjectDefinition":
		def := node.(*ast.ObjectDefinition)
		objectDefinition(def)
	case "OperationDefinition":
		def := node.(*ast.OperationDefinition)
		operationDefinition(def)
	case "OperationTypeDefinition":
		def := node.(*ast.OperationTypeDefinition)
		operationTypeDefinition(def)
	case "ScalarDefinition":
		def := node.(*ast.ScalarDefinition)
		scalarDefinition(def)
	case "SchemaDefinition":
		def := node.(*ast.SchemaDefinition)
		schemaDefinition(def)
	case "TypeDefinition":
		def := node.(ast.TypeDefinition)
		typeDefinition(def)
	case "TypeExtensionDefinition":
		def := node.(*ast.TypeExtensionDefinition)
		typeExtensionDefinition(def)
	case "TypeSystemDefinition":
		def := node.(ast.TypeSystemDefinition)
		typeSystemDefinition(def)
	case "UnionDefinition":
		def := node.(*ast.UnionDefinition)
		unionDefinition(def)
	case "VariableDefinition":
		def := node.(*ast.VariableDefinition)
		variableDefinition(def)
	default:
		fmt.Printf("    Node ---------->\n")
		fmt.Printf("      Kind: %#v\n", kind)
	}
}

func directiveDefinition(def *ast.DirectiveDefinition) {

}

func enumDefinition(def *ast.EnumDefinition) {

}

func enumValueDefinition(def *ast.EnumValueDefinition) {

}

func fieldDefinition(def *ast.FieldDefinition) {
	fmt.Printf("        Field ---------->\n")
	fmt.Printf("          Kind:        %#v\n", def.Kind)
	fmt.Printf("          Name:        %#v\n", def.Name)
	fmt.Printf("          Type:        %#v\n", def.Type)
	fmt.Printf("          Description: %#v\n", def.Description)

	fmt.Printf("          Arguments\n")
	for _, idef := range def.Arguments {
		inputValueDefinition(idef)
	}
}

func fragmentDefinition(def *ast.FragmentDefinition) {

}

func inputObjectDefinition(def *ast.InputObjectDefinition) {

}

func inputValueDefinition(def *ast.InputValueDefinition) {
	fmt.Printf("            InputValue ---------->\n")
	fmt.Printf("              Kind:         %#v\n", def.Kind)
	fmt.Printf("              Name:         %#v\n", def.Name)
	fmt.Printf("              Type:         %#v\n", def.Type)
	fmt.Printf("              DefaultValue: %#v\n", def.DefaultValue)
	fmt.Printf("              Description:  %#v\n", def.Description)
}

func interfaceDefinition(def *ast.InterfaceDefinition) {
	fmt.Printf("    Interface ---------->\n")
	fmt.Printf("      Kind:        %#v\n", def.GetKind())
	fmt.Printf("      Name:        %#v\n", def.GetName())
	fmt.Printf("      Operation:   %#v\n", def.GetOperation())
	fmt.Printf("      Description: %#v\n", def.GetDescription())
}

func objectDefinition(def *ast.ObjectDefinition) {
	fmt.Printf("    Object ---------->\n")
	fmt.Printf("      Kind:        %#v\n", def.GetKind())
	fmt.Printf("      Name:        %#v\n", def.GetName())
	fmt.Printf("      Operation:   %#v\n", def.GetOperation())
	fmt.Printf("      Description: %#v\n", def.GetDescription())

	fmt.Printf("      Fields:\n")
	for _, fdef := range def.Fields {
		fieldDefinition(fdef)
	}
}

func operationDefinition(def *ast.OperationDefinition) {

}

func operationTypeDefinition(def *ast.OperationTypeDefinition) {

}

func scalarDefinition(def *ast.ScalarDefinition) {

}

func schemaDefinition(def *ast.SchemaDefinition) {
	fmt.Printf("      Kind:      %#v\n", def.GetKind())
	fmt.Printf("      Operation: %#v\n", def.GetOperation())
}

func typeDefinition(def ast.TypeDefinition) {

}

func typeExtensionDefinition(def *ast.TypeExtensionDefinition) {

}

func typeSystemDefinition(def ast.TypeSystemDefinition) {

}

func unionDefinition(def *ast.UnionDefinition) {

}

func variableDefinition(def *ast.VariableDefinition) {

}

// Generate ...
func Generate(schema string) (*ast.Document, error) {
	doc, err := parser.Parse(
		parser.ParseParams{
			Source: schema,
			Options: parser.ParseOptions{
				NoLocation: true,
				NoSource:   true,
			},
		},
	)

	if err != nil {
		return nil, err
	}

	document(doc)

	return doc, nil
}
