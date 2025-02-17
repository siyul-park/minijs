package minijs_test

import (
	"bytes"
	"testing"

	"github.com/siyul-park/minijs"

	"github.com/stretchr/testify/assert"
)

func TestREPL_Start(t *testing.T) {
	var output bytes.Buffer
	input := bytes.NewReader([]byte(`"hello, " + "world"`))

	r := minijs.NewREPL("")

	err := r.Start(input, &output)
	assert.NoError(t, err)
	assert.Equal(t, "\"hello, world\"\n", output.String())
}
