package lb

type Selector[T any] interface {
	Get() T
	Count() int
}
