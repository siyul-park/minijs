package evaluator

import (
	"github.com/siyul-park/miniscript/lexer"
	"github.com/siyul-park/miniscript/parser"
	"github.com/siyul-park/miniscript/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEval(t *testing.T) {
	tests := []struct {
		source string
		value  types.Value
	}{
		{
			source: `-1`,
			value:  types.NewInt(-1),
		},
		{
			source: `-1.0`,
			value:  types.NewFloat(-1),
		},
		{
			source: `1 + 2`,
			value:  types.NewInt(3),
		},
		{
			source: `1.0 + 2`,
			value:  types.NewFloat(3),
		},
		{
			source: `1 + 2.0`,
			value:  types.NewFloat(3),
		},
		{
			source: `1.0 + 2.0`,
			value:  types.NewFloat(3),
		},
		{
			source: `1 - 2`,
			value:  types.NewInt(-1),
		},
		{
			source: `1.0 - 2`,
			value:  types.NewFloat(-1),
		},
		{
			source: `1 - 2.0`,
			value:  types.NewFloat(-1),
		},
		{
			source: `1.0 - 2.0`,
			value:  types.NewFloat(-1),
		},
		{
			source: `1 * 2`,
			value:  types.NewInt(2),
		},
		{
			source: `1.0 * 2`,
			value:  types.NewFloat(2),
		},
		{
			source: `1 * 2.0`,
			value:  types.NewFloat(2),
		},
		{
			source: `1.0 * 2.0`,
			value:  types.NewFloat(2),
		},
		{
			source: `1 / 2`,
			value:  types.NewInt(0),
		},
		{
			source: `1.0 / 2`,
			value:  types.NewFloat(0.5),
		},
		{
			source: `1 / 2.0`,
			value:  types.NewFloat(0.5),
		},
		{
			source: `1.0 / 2.0`,
			value:  types.NewFloat(0.5),
		},
	}

	for _, tt := range tests {
		t.Run(tt.source, func(t *testing.T) {
			l := lexer.New(tt.source)
			p := parser.New(l)

			prg, err := p.Parse()
			assert.NoError(t, err)

			env := NewEnvironment()

			val, err := Eval(prg, env)
			assert.NoError(t, err)
			assert.Equal(t, tt.value, val)
		})
	}
}
