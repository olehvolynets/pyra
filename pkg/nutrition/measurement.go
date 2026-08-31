package nutrition

import (
	"fmt"
	"math"
	"strconv"
)

const (
	AllowedDecimalDigits = 2
	MeasurementPrecision = 100 // 10 ** AllowedDecimalDigits

	FloatPrecisionBits   = 64
)

type MeasurementError struct {
	Num string
	Err error
}

func (e *MeasurementError) Error() string {
	return fmt.Sprintf("%q is not a valid measurement", e.Num)
}

type Measurement int32

func ParseMeasurement(s string) (Measurement, error) {
	if s == "" {
		return 0, nil
	}

	floatMeasurement, err := strconv.ParseFloat(s, FloatPrecisionBits)
	if err != nil {
		return 0, &MeasurementError{Num: s, Err: err}
	}

	if math.IsNaN(floatMeasurement) || math.IsInf(floatMeasurement, 0) {
		return 0, &MeasurementError{Num: s}
	}

	return NewMeasurement(floatMeasurement), nil
}

func NewMeasurement(f float64) Measurement {
	f *= MeasurementPrecision

	return Measurement(math.Round(f))
}

func (m Measurement) String() string {
	asFloat := m.Float()

	if asFloat == math.Trunc(asFloat) {
		return strconv.FormatInt(int64(asFloat), 10)
	}

	return strconv.FormatFloat(asFloat, 'f', AllowedDecimalDigits, FloatPrecisionBits)
}

func (m Measurement) Float() float64 {
	return float64(m) / MeasurementPrecision
}

func (m Measurement) Scale(ratio float64) Measurement {
	return NewMeasurement(m.Float() * ratio)
}
