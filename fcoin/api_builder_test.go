package fcoin

import (
	"github.com/jiangew/belex/exchange"
	"github.com/stretchr/testify/assert"
	"testing"
)

var builder = NewAPIBuilder()

func TestAPIBuilder_Build(t *testing.T) {
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build(exchange.FCOIN).GetExchangeName(), exchange.FCOIN)
}
