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
	defer checker.Close()

	isValid, _ := checker.IsValid("https://example.example")
	assert.Equal(test, isValid, true)

	isValid, _ = checker.IsValid("!@#$%^&*()")
	assert.Equal(test, isValid, false)

	isValid, _ = checker.IsValid("https://가나다.example")
	assert.Equal(test, isValid, false)
}

func TestIsCanBeModified(test *testing.T) {
	checker, error := New()
	if error != nil {
		test.Error(error)
	}
	defer checker.Close()

	isCanBeModified, _ := checker.IsCanBeModified("https://example.example")
	assert.Equal(test, isCanBeModified, false)

	isCanBeModified, _ = checker.IsCanBeModified("https://example.example/?url=example%3A%2F%2F")
	assert.Equal(test, isCanBeModified, true)
}
