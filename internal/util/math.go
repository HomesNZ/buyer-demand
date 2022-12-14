package util

import (
	"github.com/pkg/errors"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
	"math"
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

func IncreasedPercent[T Number](current T, prev T, decimal int) (float64, error) {
	if prev == current {
		return 0, nil
	}

	if prev == 0 {
		return 0, errors.New("The previous data can not be zero")
	}

	v := ((float64(current) - float64(prev)) / float64(prev)) * 100
	base := math.Pow10(decimal)
	return math.Round(v*base) / base, nil
}
