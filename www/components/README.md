# www/components

Re-usable components

## layout.templ

Layout is a special component that [wraps content](https://templ.guide/syntax-and-usage/template-composition#components-as-parameters) like this
```go
templ Layout(content templ.Component) {
	// Header...
	@content
	// Footer...
}
```
