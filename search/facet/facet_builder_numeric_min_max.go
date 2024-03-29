package facet

import (
	"reflect"

	"github.com/MuratYMT2/bleve/v2/numeric"
	"github.com/MuratYMT2/bleve/v2/search"
	"github.com/MuratYMT2/bleve/v2/size"
)

var reflectStaticSizeNumericMinMaxAggFacetBuilder int

func init() {
	var nfb NumericMinMaxAggFacetBuilder
	reflectStaticSizeNumericMinMaxAggFacetBuilder = int(reflect.TypeOf(nfb).Size())
}

type NumericMinMaxAggFacetBuilder struct {
	size     int
	field    string
	total    int
	missing  int
	sawValue bool
	min      *float64
	max      *float64
}

func NewNumericMinMaxAggFacetBuilder(field string, size int) *NumericMinMaxAggFacetBuilder {
	return &NumericMinMaxAggFacetBuilder{
		size:  size,
		field: field,
	}
}

func (fb *NumericMinMaxAggFacetBuilder) StartDoc() {
	fb.sawValue = false
}

func (fb *NumericMinMaxAggFacetBuilder) UpdateVisitor(term []byte) {
	fb.sawValue = true
	// only consider the values which are shifted 0
	prefixCoded := numeric.PrefixCoded(term)
	shift, err := prefixCoded.Shift()
	if err == nil && shift == 0 {
		i64, err := prefixCoded.Int64()
		if err == nil {
			f64 := numeric.Int64ToFloat64(i64)

			if fb.min == nil || f64 < *fb.min {
				fb.min = &f64
			}

			if fb.max == nil || f64 > *fb.max {
				fb.max = &f64
			}

			fb.total++
		}
	}
}

func (fb *NumericMinMaxAggFacetBuilder) EndDoc() {
	if !fb.sawValue {
		fb.missing++
	}
}

func (fb *NumericMinMaxAggFacetBuilder) Result() *search.FacetResult {
	rv := search.FacetResult{
		Field:   fb.field,
		Total:   fb.total,
		Missing: fb.missing,
		NumericMinMaxAgg: &search.NumericMinMaxAggFacet{
			Min:   fb.min,
			Max:   fb.max,
			Count: fb.total,
		},
	}

	rv.Other = 0

	return &rv
}

func (fb *NumericMinMaxAggFacetBuilder) Field() string {
	return fb.field
}

func (fb *NumericMinMaxAggFacetBuilder) Size() int {
	return reflectStaticSizeNumericMinMaxAggFacetBuilder + size.SizeOfPtr + len(fb.field)
}
