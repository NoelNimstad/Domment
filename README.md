# Domment

## What is domment?
Domment is (or will hopefully become, that is) a universal comment-based documentation system, similair to JSDoc. It is built in Go and currently only supports Go for testing purposes.

## How do I write domments?
```go
/*!
 * @? Description of the following function
 */
func double(a int) int {
    return a * 2
}
```