package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	input := struct {
		Whatever string `json:"whatever"`
	}{
		Whatever: "whatever",
	}

	{
		// Fail case, non-pointer should fail.
		assert.NotNil(t, Decode(input, ""))
		assert.NotNil(t, Decode(input, input))
	}

	{
		// Success case
		out := struct {
			OutWhatever string `json:"whatever"`
		}{}

		require.Nil(t, Decode(input, &out))
	}

}
