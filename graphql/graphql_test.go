package graphql

import (
	"reflect"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
)

type MySchema struct {
	Query
	Mutation
}

type Query struct {
	Product     Product `json:"product" description:"Get a product"`
	productArgs struct {
		ID string `description:""`
	}

	Products []Product `json:"products" description:"Get all products"`
}

func (q Query) resolveProduct(params graphql.ResolveParams) (interface{}, error) {
	return nil, nil
}

func (q Query) resolveProducts(params graphql.ResolveParams) (interface{}, error) {
	return nil, nil
}

type Mutation struct {
	AddProduct     Product `json:"addProduct" description:"Add a product"`
	addProductArgs struct {
		Name string `description:""`
	}

	UpdateProduct     Product `json:"updateProduct" description:"Update a product"`
	updateProductArgs struct {
		ID string `description:""`
	}

	RemoveProduct     bool `json:"removeProduct" description:"Remove a product"`
	removeProductArgs struct {
		ID string `description:""`
	}
}

func (m Mutation) resolveAddProduct(params graphql.ResolveParams) (interface{}, error) {
	return nil, nil
}

func (m Mutation) resolveUpdateProduct(params graphql.ResolveParams) (interface{}, error) {
	return nil, nil
}

func (m Mutation) resolveRemoveProduct(params graphql.ResolveParams) (interface{}, error) {
	return nil, nil
}

type Product struct {
	ID     string   `description:""`
	Name   string   `description:""`
	Vendor *Vendor  `description:""`
	Tags   []string `description:""`
}

type Vendor struct {
	ID   string `description:""`
	Name string `description:""`
}

func TestTypeString(t *testing.T) {
	tests := []struct {
		name               string
		fieldName          string
		fieldValue         interface{}
		expectedTypeString string
	}{
		{"id", "ID", "1111-aaaa", "ID!"},
		{"bool", "", false, "Boolean!"},
		{"bool", "", true, "Boolean!"},
		{"string", "", "cool", "String!"},
		{"float32", "", float32(3.1415), "Float!"},
		{"float64", "", float64(2.7182818284), "Float!"},
		{"int", "", int(-1), "Int!"},
		{"int8", "", int8(-8), "Int!"},
		{"int16", "", int16(-16), "Int!"},
		{"int32", "", int32(-32), "Int!"},
		{"int64", "", int64(-64), "Int!"},
		{"uint", "", uint(1), "Int!"},
		{"uint8", "", uint8(8), "Int!"},
		{"uint16", "", uint16(16), "Int!"},
		{"uint32", "", uint32(32), "Int!"},
		{"uint64", "", uint64(64), "Int!"},
		{"struct", "", Product{}, "Product"},
		{"pointer", "", &Product{}, "Product"},
		{"slice", "", []bool{}, "[Boolean!]"},
		{"slice", "", []int{}, "[Int!]"},
		{"slice", "", []float64{}, "[Float!]"},
		{"slice", "", []string{}, "[String!]"},
		{"slice", "", []Product{}, "[Product]"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fieldType := reflect.ValueOf(tc.fieldValue).Type()
			typeString := TypeString(tc.fieldName, fieldType)
			assert.Equal(t, tc.expectedTypeString, typeString)
		})
	}
}

func TestSchema(t *testing.T) {
	/* s, _ := */ Schema(&MySchema{})
	// PrintObject(s.QueryType())
	// PrintObject(s.MutationType())
}
