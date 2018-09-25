package graphql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var schema = `
schema {
  query: Query
  mutation: Mutation
}

type Query {
  asset(id: ID!): Asset
  assets: [Asset!]!

  symbol(id: ID!): Symbol
  symbols: [Symbol!]!

  font(id: ID!): Font
  fonts(name: String): [Font!]!
}

type Mutation {
  createSymbols(inputs: [InputSymbol!]!): Symbol!
  updateSymbols(inputs: [Symbol!]!): Symbol
  deleteSymbols(ids: [ID!]!): Boolean

  createFonts(inputs: [InputFont!]!): Font!
  updateFonts(inputs: [Font!]!): Font
  deleteFonts(ids: [ID!]!): Boolean
}

interface Asset {
  id: ID!
}

type Symbol implements Asset {
  id: ID!
  svg: String!
  createdAt: String!
  updatedAt: String!
}

type InputSymbol {
  svg: String!
}

type Font implements Asset {
  id: ID!
  svg: String!
  name: String!
  createdAt: String!
  updatedAt: String!
}

type InputFont {
  svg: String!
  name: String!
}
`

func TestGenerate(t *testing.T) {
	doc, err := Generate(schema)
	assert.NoError(t, err)
	assert.NotNil(t, doc)
}
