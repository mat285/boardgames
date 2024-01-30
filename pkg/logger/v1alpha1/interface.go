package v1alpha1

type Interface interface {
	Printf(string, ...interface{})

	Infof(string, ...interface{})
	Debugf(string, ...interface{})
}
