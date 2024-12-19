# go/share

Data structures **shared** between the domain models, view models, data access layer, and other packages.

Data structures used for parsing JSON may use [flexible types](https://github.com/mozey/ft). This makes the system more forgiving. For example, when parsing JSON, the **sku** field could be a string or a number. However, this field is required and therefore it uses **ft.String** instead of **ft.NString** (the latter allows [null](https://www.json.org/json-en.html)).

Struct tags for JSON are not used. Variables in **templ, html and text templates** use Go naming convention.

Constants shared between packages can also be defined here, e.g. HTTP headers.

This package is mainly for implementing data structures, and *"must have a minimal amount of logic"*. Function and methods implemented in here must not rely on services or other packages in this repo
