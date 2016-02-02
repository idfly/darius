package darius

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type configMock struct {
	mock.Mock
}

func (mock *configMock) read(file string) ([]byte, error) {
	args := mock.Called(file)
	return []byte(args.String(0)), args.Error(1)
}

func (mock *configMock) glob(pattern string) ([]string, error) {
	args := mock.Called(pattern)
	return args.Get(0).([]string), args.Error(1)
}

func newConfigTest() (Config, *configMock) {
	mock := &configMock{}
	config := Config{ReadFile: mock.read, Glob: mock.glob}
	return config, mock
}

func TestConfigLoadLoadsFile(test *testing.T) {
	config, mock := newConfigTest()
	mock.On("read", "FILE").Return("KEY: VALUE", nil)
	result, err := config.Load("FILE")
	assert.NoError(test, err)
	assert.Equal(test, map[interface{}]interface{}{"KEY": "VALUE"}, result)
}

func TestConfigLoadIncludesFile(test *testing.T) {
	config, mock := newConfigTest()
	mock.On("read", "FILE1").Return(`{KEY: "${include FILE2}"}`, nil)
	mock.On("read", "FILE2").Return("VALUE", nil)
	result, err := config.Load("FILE1")
	assert.NoError(test, err)
	assert.Equal(test, map[interface{}]interface{}{"KEY": "VALUE"}, result)
}

func TestConfigLoadLoadsArrayRecursive(test *testing.T) {
	config, mock := newConfigTest()
	mock.On("read", "FILE").Return(`{"KEY": [{"KEY": "VALUE"}]}`, nil)
	result, err := config.Load("FILE")
	assert.NoError(test, err)
	array := []interface{}{map[interface{}]interface{}{"KEY": "VALUE"}}
	assert.Equal(test, map[interface{}]interface{}{"KEY": array}, result)
}

func TestConfigLoadIncludesFileWithOldSyntax(test *testing.T) {
	config, mock := newConfigTest()
	mock.On("read", "FILE1").Return(`{KEY: "$include FILE2"}`, nil)
	mock.On("read", "FILE2").Return("VALUE", nil)
	result, err := config.Load("FILE1")
	assert.NoError(test, err)
	assert.Equal(test, map[interface{}]interface{}{"KEY": "VALUE"}, result)
}

func TestConfigLoadIncludesFileFromSubkey(test *testing.T) {
	config, mock := newConfigTest()
	mock.On("read", "FILE1").Return(`{KEY: {KEY: "${include FILE2}"}}`, nil)
	mock.On("read", "FILE2").Return("VALUE", nil)
	result, err := config.Load("FILE1")
	assert.NoError(test, err)
	expected := map[interface{}]interface{}{"KEY": "VALUE"}
	assert.Equal(test, map[interface{}]interface{}{"KEY": expected}, result)
}

func TestConfigLoadIncludesFileFromIncludee(test *testing.T) {
	config, mock := newConfigTest()
	mock.On("read", "FILE1").Return(`{KEY: "${include FILE2}"}`, nil)
	mock.On("read", "FILE2").Return(`{KEY: "${include FILE3}"}`, nil)
	mock.On("read", "FILE3").Return("VALUE", nil)
	result, err := config.Load("FILE1")
	assert.NoError(test, err)
	expected := map[interface{}]interface{}{"KEY": "VALUE"}
	assert.Equal(test, map[interface{}]interface{}{"KEY": expected}, result)
}

func TestConfigLoadReturnsErrorOnRecursiveInclude(test *testing.T) {
	config, mock := newConfigTest()
	mock.On("read", "FILE1").Return(`{KEY: "${include FILE2}"}`, nil)
	mock.On("read", "FILE2").Return(`{KEY: "${include FILE1}"}`, nil)
	_, err := config.Load("FILE1")
	assert.Error(test, err)
}

func TestConfigLoadIncludesWithRootPath(test *testing.T) {
	config, mock := newConfigTest()
	mock.On("read", "FILE").Return(`{KEY: "${include /ROOT/FILE}"}`, nil)
	mock.On("read", "/ROOT/FILE").Return(`VALUE`, nil)
	result, err := config.Load("FILE")
	assert.NoError(test, err)
	assert.Equal(test, map[interface{}]interface{}{"KEY": "VALUE"}, result)
}

func TestConfigLoadIncludesWithRelativePath(test *testing.T) {
	config, mock := newConfigTest()
	mock.On("read", "/ROOT/FILE1").Return(`{KEY: "${include FILE2}"}`, nil)
	mock.On("read", "/ROOT/FILE2").Return(`VALUE`, nil)
	result, err := config.Load("/ROOT/FILE1")
	assert.NoError(test, err)
	assert.Equal(test, map[interface{}]interface{}{"KEY": "VALUE"}, result)
}

func TestConfigLoadMergesFilesWithGlob(test *testing.T) {
	config, mock := newConfigTest()
	mock.On("read", "/ROOT/FILE1").Return(`{KEY: "${include PATH/*}"}`, nil)
	mock.On("read", "/FILE1").Return(`{KEY1: VALUE1}`, nil)
	mock.On("read", "/FILE2").Return(`{KEY2: VALUE2}`, nil)
	mock.On("glob", "/ROOT/PATH/*").Return([]string{"/FILE1", "/FILE2"}, nil)
	result, err := config.Load("/ROOT/FILE1")
	assert.NoError(test, err)
	value := map[interface{}]interface{}{"KEY1": "VALUE1", "KEY2": "VALUE2"}
	expected := map[interface{}]interface{}{"KEY": value}
	assert.Equal(test, expected, result)
}
