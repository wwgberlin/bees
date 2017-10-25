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
	return distance(r, g, b, left) < distance(r, g, b, right)
}
func distance(r, g, b uint8, prod Product) float64 {
	return math.Sqrt(float64((prod.b-b)*(prod.b-b) + (prod.r-r)*(prod.r-r) + (prod.r-g)*(prod.r-g)))
}
