package types

type APIRequestable interface {
    ToBytes() ([]byte, error) 
}
