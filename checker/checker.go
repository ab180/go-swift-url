package checker

/*
#include <stdlib.h>
*/
import "C"
import (
	_ "embed"

	"context"
	"errors"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed checker.wasm
var checkerWASM []byte

var errorUnexpectedResults = errors.New("unexpected results are received from WebAssembly runtime")
var errorOutOfMemory = errors.New("failed on write data to memory of WebAssembly runtime")

type Checker interface {
	Close()
	IsValid(url string) (bool, error)
	IsCanBeModified(url string) (bool, error)
}

type checker struct {
	context context.Context
	runtime wazero.Runtime
	module  api.Module
}

// Creates checker instance and initialize WebAssembly runtime
func New() (Checker, error) {
	context := context.Background()
	runtime := wazero.NewRuntime(context)

	_, error := wasi_snapshot_preview1.Instantiate(context, runtime)
	if error != nil {
		runtime.Close(context)
		return nil, error
	}

	module, error := runtime.Instantiate(context, checkerWASM)
	if error != nil {
		runtime.Close(context)
		return nil, error
	}

	_, error = module.ExportedFunction("_initialize").Call(context)
	if error != nil {
		module.Close(context)
		runtime.Close(context)
		return nil, error
	}

	checker := checker{
		context,
		runtime,
		module,
	}

	return &checker, nil
}

// Close WebAssembly runtime
func (checker *checker) Close() {
	checker.module.Close(checker.context)
	checker.runtime.Close(checker.context)
}

// Checks Swift's `URL` can be initialized with `url`
func (checker *checker) IsValid(url string) (bool, error) {
	bytes := []byte(url)
	length := len(bytes)

	pointer, error := checker.allocate(uint64(length))
	if error != nil {
		return false, error
	}

	isWritten := checker.module.ExportedMemory("memory").Write(uint32(pointer), bytes)
	if !isWritten {
		return false, errorOutOfMemory
	}

	results, error := checker.module.ExportedFunction("is_valid").Call(checker.context, pointer)
	if error != nil {
		return false, error
	}
	if len(results) == 0 {
		return false, errorUnexpectedResults
	}

	error = checker.deallocate(pointer)
	if error != nil {
		return false, error
	}

	return results[0] != 0, nil
}

// Checks `url` can be modified by Swift's `URLComponents.queryItems`
//
// When Swift's `URLComponents` can not be initialized with `url`, then return false
func (checker *checker) IsCanBeModified(url string) (bool, error) {
	bytes := []byte(url)
	length := len(bytes)

	pointer, error := checker.allocate(uint64(length))
	if error != nil {
		return false, error
	}

	isWritten := checker.module.ExportedMemory("memory").Write(uint32(pointer), bytes)
	if !isWritten {
		return false, errorOutOfMemory
	}

	results, error := checker.module.ExportedFunction("is_can_be_modified").Call(checker.context, pointer)
	if error != nil {
		return false, error
	}
	if len(results) == 0 {
		return false, errorUnexpectedResults
	}

	error = checker.deallocate(pointer)
	if error != nil {
		return false, error
	}

	return results[0] != 0, nil
}

// Allocates `length` bytes from WebAssembly runtime and return address
func (checker *checker) allocate(length uint64) (uint64, error) {
	results, error := checker.module.ExportedFunction("allocate").Call(checker.context, length)
	if error != nil {
		return 0, error
	}
	if len(results) == 0 {
		return 0, errorUnexpectedResults
	}

	return results[0], nil
}

// Deallocates `pointer` from WebAssembly runtime
func (checker *checker) deallocate(pointer uint64) error {
	_, error := checker.module.ExportedFunction("deallocate").Call(checker.context, pointer)
	if error != nil {
		return error
	}

	return nil
}
