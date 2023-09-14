import Foundation

/// Checks `URL` can be initialized with `string`
/// - Note: `stringPointer` must be initialized by `allocate`
@_cdecl("is_valid")
func isValid(_ stringPointer: UnsafePointer<Int8>) -> Bool {
    let string = String(cString: stringPointer)

    return URL(string: string) != nil
}

/// Checks url string can be modified by `URLComponents.queryItems`
/// - Note: when `URLComponents` can not be initialized with `string`, then return false
/// - Note: `stringPointer` must be initialized by `allocate`
@_cdecl("is_can_be_modified")
func isCanBeModified(_ stringPointer: UnsafePointer<Int8>) -> Bool {
    let string = String(cString: stringPointer)
    guard var editor = URLComponents(string: string) else {
        return false
    }

    let queryItems = editor.queryItems
    editor.queryItems = queryItems
    guard let modified = editor.string else {
        return false
    }

    return modified != string
}

/// Allocates `length` bytes from WebAssembly runtime and return address
@_cdecl("allocate")
func allocate(_ length: Int) -> UnsafeMutableRawPointer {
    return malloc(length)
}

/// Deallocates `pointer` from WebAssembly runtime
@_cdecl("deallocate")
func deallocate(_ pointer: UnsafeMutableRawPointer) {
    free(pointer)
}
