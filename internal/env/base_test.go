package env

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"testing"
)

const (
	envVarKey = "TESTING_ENV_VAR"
)

func TestEnv_HasValue(t *testing.T) {
	envVarValue := "TEST_CASE_TEST_ENV_VALUE"

	// set up
	setEnvVar(t, envVarValue)

	// act
	e := newEnv(envVarKey)

	// assert
	assert.Equal(t, envVarKey, e.GetKey())
	assert.Equal(t, envVarValue, e.GetValue())

	// tear down
	unsetEnvVar(t)
}

func TestEnv_NoValue(t *testing.T) {
	// act
	e := newEnv(envVarKey)

	// assert
	assert.Equal(t, envVarKey, e.GetKey())
	assert.Equal(t, "", e.GetValue())
}

type mockEnvGetter struct {
	*mock.Mock
}

func (g *mockEnvGetter) GetEnv(key string) string {
	args := g.Mock.Called(key)
	return args.String(0)
}

func TestEnv_OnlyCallGetOnce(t *testing.T) {
	envVarValue := "TEST_CASE_TEST_ENV_VALUE"

	// set up
	setEnvVar(t, envVarValue)
	mockGetter := mockEnvGetter{&mock.Mock{}}
	mockGetter.On("GetEnv", envVarKey).Return(envVarValue)

	// act
	e := newEnvWithCustomGet(envVarKey, &mockGetter)
	e.GetValue()
	actual := e.GetValue()

	// assert
	assert.Equal(t, envVarValue, actual)
	mockGetter.AssertNumberOfCalls(t, "GetEnv", 1)

	unsetEnvVar(t)
}

func setEnvVar(t *testing.T, value string) {
	err := os.Setenv(envVarKey, value)
	assert.NoError(t, err)
}

func unsetEnvVar(t *testing.T) {
	err := os.Unsetenv(envVarKey)
	assert.NoError(t, err)
}
