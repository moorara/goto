package graphql

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/graphql-go/graphql"
)

type (
	types map[string]graphql.Output
)

const (
	idName         = "ID"
	queryName      = "Query"
	mutationName   = "Mutation"
	descriptionTag = "description"
)

// TypeString returns the GraphQL type for a struct field as string
func TypeString(fieldName string, fieldType reflect.Type) string {
	if fieldName == idName {
		return "ID!"
	}

	switch fieldType.Kind() {
	case reflect.Bool:
		return "Boolean!"
	case reflect.String:
		return "String!"
	case reflect.Float32, reflect.Float64:
		return "Float!"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "Int!"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "Int!"
	case reflect.Struct:
		return fieldType.Name()
	case reflect.Ptr:
		return TypeString(fieldName, fieldType.Elem())
	case reflect.Slice:
		return "[" + TypeString(fieldName, fieldType.Elem()) + "]"
	default:
		return ""
	}
}

// Output returns the graphql.Output for a GraphQL type
func Output(cache types, graphqlType string) graphql.Output {
	// Look up the cache
	output := cache[graphqlType]
	if output != nil {
		return output
	}

	l := len(graphqlType)
	if graphqlType[l-1:] == "!" {
		output = Output(cache, graphqlType[:l-1])
		output = graphql.NewNonNull(output)
	} else if graphqlType[0:1] == "[" && graphqlType[l-1:] == "]" {
		output = Output(cache, graphqlType[1:l-1])
		output = graphql.NewList(output)
	}

	// Update the cache
	cache[graphqlType] = output

	return output
}

// Object returns the graphql.Object for a struct
func Object(cache types, value reflect.Value) *graphql.Object {
	graphqlObject := new(graphql.Object)
	graphqlName := value.Type().Name()

	// Look up the cache
	graphqlObject, ok := cache[graphqlName].(*graphql.Object)
	if ok {
		return graphqlObject
	}

	graphqlFields := graphql.Fields{}

	// Extract GraphQL types
	for i := 0; i < value.NumField(); i++ {
		f := value.Type().Field(i)
		v := value.Field(i)
		t := v.Type()

		// Skip unexported fields
		if !v.CanSet() {
			continue
		}

		switch t.Kind() {
		case reflect.Struct:
			Object(cache, v)
		case reflect.Slice, reflect.Ptr:
			tt := t.Elem()               // slice/pointer type
			vv := reflect.New(tt).Elem() // slice/pointer value
			if tt.Kind() == reflect.Struct {
				Object(cache, vv)
			}
		}

		graphqlType := TypeString(f.Name, v.Type())
		graphqlOutput := Output(cache, graphqlType)
		graphqlDescription := f.Tag.Get(descriptionTag)

		graphqlFields[f.Name] = &graphql.Field{
			Type:        graphqlOutput,
			Description: graphqlDescription,
		}

		// TODO
		fieldTokens := parseFieldName(f.Name)
		argsFieldName := strings.ToLower(fieldTokens[0]) + strings.Join(fieldTokens[1:], "") + "Args"
		argsV := value.FieldByName(argsFieldName)

		// If a field for arguments with anonymouse struct type is defined ...
		if argsV.IsValid() {
			argsT := argsV.Type()
			if argsT.Kind() == reflect.Struct && argsT.Name() == "" {
				// TODO
				for j := 0; j < argsV.NumField(); j++ {
					f := argsT.Field(j)
					v := argsV.Field(j)
					t := v.Type()

					fmt.Printf("  ----> %+v:  %+v %+v \n", argsFieldName, f.Name, t)
				}
			}
		}

		// TODO
		resolverFuncName := "resolve" + f.Name
		resolverV := value.MethodByName(resolverFuncName)

		// TODO
		if resolverV.IsValid() {
			resolverT := resolverV.Type()
			fmt.Printf("  ----> %+v:  %+v \n", resolverFuncName, resolverT)
		}
	}

	// Create a new graphql.Object for the struct
	graphqlObject = graphql.NewObject(
		graphql.ObjectConfig{
			Name:   graphqlName,
			Fields: graphqlFields,
		},
	)

	// Update the cache
	cache[graphqlName] = graphqlObject

	return graphqlObject
}

// Schema returns the graphql.Schema for a schema struct
func Schema(schema interface{}) (*graphql.Schema, error) {
	schemaValue := reflect.ValueOf(schema)
	schemaType := reflect.TypeOf(schema)

	// If a pointer is passed, navigate to the value
	if schemaType.Kind() == reflect.Ptr {
		schemaValue = schemaValue.Elem()
		schemaType = schemaType.Elem()
	}

	// Make sure schema is a struct
	if schemaType.Kind() != reflect.Struct {
		return nil, errors.New("schema should be a struct")
	}

	// Make sure schema has Query
	queryValue := schemaValue.FieldByName(queryName)
	if !queryValue.IsValid() || queryValue.Type().Name() != queryName {
		return nil, errors.New("schema should have a field Query with type Query")
	}

	// Make sure schema has Mutation
	mutationValue := schemaValue.FieldByName(mutationName)
	if !mutationValue.IsValid() || mutationValue.Type().Name() != mutationName {
		return nil, errors.New("schema should have a field Mutation with type Mutation")
	}

	// Define a cache for GraphQL types
	cache := types{
		"ID":      graphql.ID,
		"Int":     graphql.Int,
		"Float":   graphql.Float,
		"String":  graphql.String,
		"Boolean": graphql.Boolean,
	}

	queryObject := Object(cache, queryValue)
	mutationObject := Object(cache, mutationValue)

	graphqlSchema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    queryObject,
			Mutation: mutationObject,
		},
	)

	if err != nil {
		return nil, err
	}

	return &graphqlSchema, nil
}

// PrintOutput prints a graphql.Output
func PrintOutput(out graphql.Output) {
	_, list := out.(*graphql.List)
	_, nonNull := out.(*graphql.NonNull)

	fmt.Printf("    TYPE ..............................................\n")
	fmt.Printf("      Name:        %v\n", out.Name())
	fmt.Printf("      List:        %v\n", list)
	fmt.Printf("      Not Null:    %v\n", nonNull)
	fmt.Printf("      Description: %v\n", out.Description())
}

// PrintArgument prints a graphql.Argument
func PrintArgument(arg *graphql.Argument) {
	fmt.Printf("    ARG ..............................................\n")
	fmt.Printf("      Name:        %v\n", arg.Name())
	fmt.Printf("      Description: %v\n", arg.Description())
}

// PrintFieldDefinition prints a graphql.FieldDefinition
func PrintFieldDefinition(fd *graphql.FieldDefinition) {
	fmt.Printf("  FIELD -----------------------------------------------\n")
	fmt.Printf("    Name:        %v\n", fd.Name)
	fmt.Printf("    Description: %v\n", fd.Description)

	PrintOutput(fd.Type)
	for _, arg := range fd.Args {
		PrintArgument(arg)
	}

	fmt.Printf("    Resolve:     %v\n", fd.Resolve)
}

// PrintObject prints a graphql.Object
func PrintObject(obj *graphql.Object) {
	fmt.Printf("OBJECT ================================================\n")
	fmt.Printf("  Name:        %v\n", obj.Name())
	fmt.Printf("  Description: %v\n", obj.Description())

	for _, fd := range obj.Fields() {
		PrintFieldDefinition(fd)
	}
}
