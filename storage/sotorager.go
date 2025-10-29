package storage

type Storager interface {
	Store(filename string, data []byte) error
	Load() ([]byte, error)
}
