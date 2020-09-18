package wave

// Filtering abstract of filter wave
type Filtering interface {
	Filtering(v float64) float64
}

// BatchFiltering filter dataset one time
type BatchFiltering interface {
	BatchFiltering(data []float64) []float64
}

type ComplexFiltering interface {
	Filtering
	BatchFiltering
}
