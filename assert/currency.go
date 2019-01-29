package assert

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/currency"
	"testing"
)

func CurrencyEquals(t *testing.T, expected string, actual currency.Amount) {
	formatted := fmt.Sprintf("%s", actual)
	assert.Equal(t, expected, formatted)
}