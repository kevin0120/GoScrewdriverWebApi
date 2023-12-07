package typeDef

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_PsetUnmarshalJSON(t *testing.T) {
	data := []byte("11")
	pset := Pset{}
	if err := json.Unmarshal(data, &pset); err != nil {
		assert.Equal(t, 11, pset.Value)
	}
}
