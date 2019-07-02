package builder

import (
	"github.com/jiangew/belex"
	"github.com/stretchr/testify/assert"
	"testing"
)

var builder = NewAPIBuilder()

func TestAPIBuilder_Build(t *testing.T) {
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build(belex.HUOBI_PRO).GetExchangeName(), belex.HUOBI_PRO)
	assert.Equal(t, builder.APIKey("").APISecretkey("").BuildFuture(belex.HBDM).GetExchangeName(), belex.HBDM)
}
