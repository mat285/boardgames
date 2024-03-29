package v1alpha1

type Deserializer interface {
	Deserialize([]byte, interface{}) error
}

type TypeProvider interface {
	EmptyType() interface{}
}
