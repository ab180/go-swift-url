# go-swift-url

Using Swift URL from Go over WebAssembly.

## Notice

WebAssembly module works differently with `>= iOS 17`.

WebAssembly module is built with Swift 5.8. And it's `Core Foundation` of `Foundation` is from `apple/swift-corelibs-foundation@swift-5.8-RELEASE` (It has iOS 15's implementation). And since iOS 17, iOS's implementation checks url validity less strictly and normalizes url automatically.

```
// iOS 15, 16: nil
// iOS 17: Optional(https://airbridge.io/?%7B%7D)
URL(string: "https://airbridge.io/?{}")
```

> `Foundation`'s URL, NSURL class wraps `Core Foundation`'s CFURL class.

## Install

```sh
go get github.com/ab180/go-swift-url
```

## Usage

```go
import "github.com/ab180/go-swift-url/checker"

checker, error := checker.New()
if error != nil {
	test.Error(error)
}
defer checker.Close()

// true
isValid, _ := checker.IsValid("https://example.example")
// false
isValid, _ = checker.IsValid("!@#$%^&*()")
// false
isCanBeModified, _ := checker.IsCanBeModified("https://example.example")
// true
isCanBeModified, _ = checker.IsCanBeModified("https://example.example/?url=example%3A%2F%2F")
```

### IsValid

Checks Swift's `URL` can be initialized with `url`.

### IsCanBeModified

Checks `url` can be modified by Swift's `URLComponents.queryItems`.

When Swift's `URLComponents` can not be initialized with `url`, then return false.

#### Example

```swift
import Foundation

var url = URL(string: "https://example.example/?url=example%3A%2F%2F")!
var editor = URLComponents(url: url, resolvingAgainstBaseURL: false)!
let queryItems = editor.queryItems
editor.queryItems = queryItems

// https://example.example/?url=example://
print(editor.url)
```

## Develop

1. Modify `checker.swift`.
2. Run `./script/build.sh`.
    - If you add more functions, then you must add below line to `script/build.sh`.
    - `-Xlinker --export={FUNCTION_NAME} \`
3. Modify `checker.go`.
4. Commit and push `checker.swift`, `checker.wasm`, `checker.go`.

### Requirement

- <https://github.com/swiftwasm/swift>

> Requirement is needed for development not for usage.
