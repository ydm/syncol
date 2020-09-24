package syncol

type Collection interface {
	Init()
	Put(item interface{})
	Get() (interface{}, bool)
}
