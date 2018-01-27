package graphviz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEdge(t *testing.T) {
	tests := []struct {
		name            string
		from            string
		to              string
		edgeType        string
		edgeDir         string
		label           string
		color           string
		style           string
		arrowhead       string
		expectedDotCode string
	}{
		{
			"SimpleEdge",
			"root", "left", EdgeTypeDirected, "",
			"", "", "", "",
			`root -> left [];`,
		},
		{
			"EdgeWithType",
			"root", "right", EdgeTypeUndirected, "",
			"normal", "", "", "",
			`root -- right [label="normal"];`,
		},
		{
			"EdgeWithProperties",
			"parent", "child", EdgeTypeDirected, EdgeDirNone,
			"red", ColorGold, StyleDashed, ArrowheadBox,
			`parent -> child [dirType=none, label="red", color=gold, style=dashed, arrowhead=box];`,
		},
		{
			"EdgeWithProperties",
			"parent", "child", EdgeTypeUndirected, EdgeDirBoth,
			"black", ColorOrchid, StyleDotted, ArrowheadOpen,
			`parent -- child [dirType=both, label="black", color=orchid, style=dotted, arrowhead=open];`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := NewEdge(tc.from, tc.to, tc.edgeType, tc.edgeDir, tc.label, tc.color, tc.style, tc.arrowhead)
			assert.Equal(t, tc.expectedDotCode, e.DotCode())
		})
	}
}
