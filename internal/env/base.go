package env

import "os"

type envGetter interface {
	GetEnv(string) string
}

type defaultEnvGetter struct{}

func (g *defaultEnvGetter) GetEnv(key string) string {
	return os.Getenv(key)
}

type envOption = func(*env)

// env defines the internal data structure for environment variables.
type env struct {
	key, value string
	isRead     bool
	getter     envGetter
}

// GetKey returns the key of the environment variable.
func (e *env) GetKey() string {
	return e.key
}

// GetValue returns the value of the environment variable. If the environment variable is not
// set, an empty string is returned. Furthermore, the read operation is performed lazily. The
// result is cached for further use.
func (e *env) GetValue() string {
	if !e.isRead {
		e.value = e.getter.GetEnv(e.key)
		e.isRead = true
	}
	return e.value
}

// newEnv creates a new env available to use.
func newEnv(key string, opts ...envOption) env {
	e := env{key: key, getter: &defaultEnvGetter{}}
	for _, opt := range opts {
		opt(&e)
	}
	return e
}

// withGetter creates an option for newEnv. This option replaces the default env getter
// with the one passed in.
func withGetter(getter envGetter) envOption {
	return func(e *env) {
		e.getter = getter
	}
}
