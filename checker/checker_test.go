package checker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValid(test *testing.T) {
	checker, error := New()
	if error != nil {
		test.Error(error)
	}

	isValid, error := checker.IsValid("https://example.example")
	assert.Equal(test, isValid, true)
	assert.Equal(test, error, nil)

	isValid, error = checker.IsValid("!@#$%^&*()")
	assert.Equal(test, isValid, false)
	assert.Equal(test, error, nil)

	isValid, error = checker.IsValid("https://가나다.example")
	assert.Equal(test, isValid, false)
	assert.Equal(test, error, nil)
}

func TestIsCanBeModified(test *testing.T) {
	checker, error := New()
	if error != nil {
		test.Error(error)
	}

	isCanBeModified, _ := checker.IsCanBeModified("https://example.example")
	assert.Equal(test, isCanBeModified, false)

	isCanBeModified, _ = checker.IsCanBeModified("https://example.example/?url=example%3A%2F%2F")
	assert.Equal(test, isCanBeModified, true)
}

func BenchmarkIsValid(benchmark *testing.B) {
	checker, error := New()
	if error != nil {
		benchmark.Error(error)
	}

	benchmark.ResetTimer()
	for index := 0; index < benchmark.N; index++ {
		_, _ = checker.IsValid("https://example.example")
	}
	benchmark.StopTimer()
}

func BenchmarkIsCanBeModified(benchmark *testing.B) {
	checker, error := New()
	if error != nil {
		benchmark.Error(error)
	}

	benchmark.ResetTimer()
	for index := 0; index < benchmark.N; index++ {
		_, _ = checker.IsCanBeModified("https://example.example")
	}
	benchmark.StopTimer()
}
