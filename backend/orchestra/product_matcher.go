package orchestra

import (
	"math"
	"reflect"
)

type ProductMatcher struct {
	products []Product
	r        uint8
	g        uint8
	b        uint8
}

func newProductMatcher(products []Product, r, g, b uint8) ProductMatcher {
	dest := make([]Product, len(products))
	reflect.Copy(reflect.ValueOf(dest), reflect.ValueOf(products))
	return ProductMatcher{products: dest, r: r, g: g, b: b}
}
func (a ProductMatcher) Len() int      { return len(a.products) }
func (a ProductMatcher) Swap(i, j int) { a.products[i], a.products[j] = a.products[j], a.products[i] }
func (a ProductMatcher) Less(i, j int) bool {
	left := a.products[i]
	right := a.products[j]
	r, g, b := a.r, a.g, a.b
	return left.distance(r, g, b) < right.distance(r, g, b)
}
func (p *Product) distance(r, g, b uint8) float64 {
	return math.Sqrt(float64((p.r-r)*(p.r-r) + (p.g-g)*(p.g-g) + (p.b-b)*(p.b-b)))
}
