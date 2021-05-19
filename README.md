# Partial

Partial allows you to determine which fields in a struct do not have a zero value field, by tag. This project was inspired by an issue I was having with partial DB updates.

Go's basic types are covered automatically. Custom types require implementing the Partials interface.

TODO:

Add tests
