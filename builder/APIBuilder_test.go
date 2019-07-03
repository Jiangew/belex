package builder

import (
	"github.com/jiangew/belex"
	"github.com/stretchr/testify/assert"
	"testing"
)

var builder = NewAPIBuilder()

func TestAPIBuilder_Build(t *testing.T) {
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build(belex.FCOIN).GetExchangeName(), belex.FCOIN)
}
