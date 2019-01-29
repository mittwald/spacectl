package costestimator

import "fmt"

type Unit interface {
	Format(n int) string
}

type SimpleUnit struct {
	fn func(n int) string
}

func (u *SimpleUnit) Format(n int) string {
	return u.fn(n)
}

func NewUnit(name string) Unit {
	return &SimpleUnit{
		fn: func(n int) string {
			return fmt.Sprintf("%d %s", n, name)
		},
	}
}
