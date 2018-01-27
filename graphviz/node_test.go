package graphviz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNode(t *testing.T) {
	tests := []struct {
		name            string
		nodeName        string
		group           string
		label           string
		color           string
		style           string
		shape           string
		fontcolor       string
		fontname        string
		expectedDotCode string
	}{
		{
			"SimpleNode",
			"root", "",
			"", "", "", "", "", "",
			`root [];`,
		},
		{
			"NodeWithLabel",
			"root", "",
			"root", "", "", "", "", "",
			`root [label="root"];`,
		},
		{
			"NodeWithGroup",
			"struct0", "group0",
			"<f0> left|<f1> middle|<f2> right",
			ColorBlue, StyleBold, ShapeBox, ColorGray, "",
			`struct0 [group=group0, label="<f0> left|<f1> middle|<f2> right", color=blue, style=bold, shape=box, fontcolor=gray];`,
		},
		{
			"ComplexNode",
			"struct1", "group1",
			"a | { b | { c | <here> d | e } | f } | g | h",
			ColorNavy, StyleDashed, ShapeOval, ColorBlack, "Arial",
			`struct1 [group=group1, label="a | { b | { c | <here> d | e } | f } | g | h", color=navy, style=dashed, shape=oval, fontcolor=black, fontname="Arial"];`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			n := NewNode(tc.nodeName, tc.group, tc.label, tc.color, tc.style, tc.shape, tc.fontcolor, tc.fontname)
			assert.Equal(t, tc.expectedDotCode, n.DotCode())
		})
	}
}
