package syncol

type Collection[T any] interface {
	Init()
	Put(item T)
	Get() (T, bool)
}
