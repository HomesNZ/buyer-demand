package util

import (
	"github.com/pkg/errors"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

type Number interface {
	constraints.Float | constraints.Integer
}

func Median[T Number](data []T) float64 {
	dataCopy := make([]T, len(data))
	copy(dataCopy, data)

	slices.Sort(dataCopy)

	var median float64
	l := len(dataCopy)
	if l == 0 {
		return 0
	} else if l%2 == 0 {
		median = float64((dataCopy[l/2-1] + dataCopy[l/2]) / 2.0)
	} else {
		median = float64(dataCopy[l/2])
	}

	return median
}

func IncreasedPercent[T Number](a T, b T) (float64, error) {
	if a == 0 {
		return 0, errors.New("The old data can not be zero")
	}

	return float64((b-a)/a) * 100, nil
}
