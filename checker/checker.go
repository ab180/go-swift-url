package checker

import (
	_ "embed"

	"errors"

	"github.com/bytecodealliance/wasmtime-go/v13"
)

//go:embed checker.wasm
var checkerWASM []byte

var errorExportNotFound = errors.New("expected export is not found from WebAssembly runtime")
var errorUnexpectedResults = errors.New("unexpected results are received from WebAssembly runtime")

type Checker interface {
	IsValid(url string) (bool, error)
	IsCanBeModified(url string) (bool, error)
}

type checker struct {
	instance                *wasmtime.Instance
	module                  *wasmtime.Module
	store                   *wasmtime.Store
	memory                  *wasmtime.Memory
	functionIsValid         *wasmtime.Func
	functionIsCanBeModified *wasmtime.Func
	functionAllocate        *wasmtime.Func
	functionDeallocate      *wasmtime.Func
}

// Creates checker instance and initialize WebAssembly runtime
func New() (Checker, error) {
	store := wasmtime.NewStore(wasmtime.NewEngine())
	store.SetWasi(wasmtime.NewWasiConfig())

	module, error := wasmtime.NewModule(store.Engine, checkerWASM)
	if error != nil {
		return nil, error
	}

	linker := wasmtime.NewLinker(store.Engine)
	error = linker.DefineWasi()
	if error != nil {
		return nil, error
	}

	instance, error := linker.Instantiate(store, module)
	if error != nil {
		return nil, error
	}

	initialize := instance.GetFunc(store, "_initialize")
	if initialize == nil {
		return nil, errorExportNotFound
	}

	_, error = initialize.Call(store)
	if error != nil {
		return nil, error
	}

	memory := instance.GetExport(store, "memory").Memory()
	if memory == nil {
		return nil, errorExportNotFound
	}

	functionIsValid := instance.GetFunc(store, "is_valid")
	if functionIsValid == nil {
		return nil, errorExportNotFound
	}

	functionIsCanBeModified := instance.GetFunc(store, "is_can_be_modified")
	if functionIsCanBeModified == nil {
		return nil, errorExportNotFound
	}

	functionAllocate := instance.GetFunc(store, "allocate")
	if functionAllocate == nil {
		return nil, errorExportNotFound
	}

	functionDeallocate := instance.GetFunc(store, "deallocate")
	if functionDeallocate == nil {
		return nil, errorExportNotFound
	}

	checker := checker{
		instance,
		module,
		store,
		memory,
		functionIsValid,
		functionIsCanBeModified,
		functionAllocate,
		functionDeallocate,
	}

	return &checker, nil
}

// Checks Swift's `URL` can be initialized with `url`
func (checker *checker) IsValid(url string) (bool, error) {
	bytes := append([]byte(url), 0)
	length := int32(len(bytes))

	pointer, error := checker.allocate(length)
	if error != nil {
		return false, error
	}

	copy(checker.memory.UnsafeData(checker.store)[pointer:], bytes)

	result, error := checker.functionIsValid.Call(checker.store, pointer)
	if error != nil {
		return false, error
	}

	isValid, isSuccess := result.(int32)
	if !isSuccess {
		return false, errorUnexpectedResults
	}

	error = checker.deallocate(pointer)
	if error != nil {
		return false, error
	}

	return isValid != 0, nil
}

// Checks `url` can be modified by Swift's `URLComponents.queryItems`
//
// When Swift's `URLComponents` can not be initialized with `url`, then return false
func (checker *checker) IsCanBeModified(url string) (bool, error) {
	bytes := append([]byte(url), 0)
	length := int32(len(bytes))

	pointer, error := checker.allocate(length)
	if error != nil {
		return false, error
	}

	copy(checker.memory.UnsafeData(checker.store)[pointer:], bytes)

	result, error := checker.functionIsCanBeModified.Call(checker.store, pointer)
	if error != nil {
		return false, error
	}

	isCanBeModified, isSuccess := result.(int32)
	if !isSuccess {
		return false, errorUnexpectedResults
	}

	error = checker.deallocate(pointer)
	if error != nil {
		return false, error
	}

	return isCanBeModified != 0, nil
}

// Allocates `length` bytes from WebAssembly runtime and return address
func (checker *checker) allocate(length int32) (int32, error) {
	result, error := checker.functionAllocate.Call(checker.store, length)
	if error != nil {
		return 0, error
	}

	pointer, isSuccess := result.(int32)
	if !isSuccess {
		return 0, errorUnexpectedResults
	}

	return pointer, nil
}

// Deallocates `pointer` from WebAssembly runtime
func (checker *checker) deallocate(pointer int32) error {
	_, error := checker.functionDeallocate.Call(checker.store, pointer)
	if error != nil {
		return error
	}

	return nil
}
