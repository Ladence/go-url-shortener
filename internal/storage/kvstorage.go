package storage

type KvStorage interface {
	Push(key string, value any) error
	Get(key string) (any, error)
}
