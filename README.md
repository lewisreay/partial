# Partial

This package allows you to determine which fields in a struct do not have a zero value, by tag.

The initial use case was for determining which fields need updating from a PATCH request, but the use case could extend much further.

Go's basic types are covered automatically. Custom types require implementing the Partials interface.

`example.go` displays an example of how to use this package.
