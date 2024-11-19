# go/router

Define HTTP route handlers in here, and some logic. However, business logic must be defined in `go/model`

Route handlers respond with shared data types, e.g. `share.StockAvailableGet`. These types have corresponding JSON stubs and template files, e.g. `www/api/stock/available/GET.json` and `GET.templ`

Use go naming convention for shared data types, e.g. `StockItem.Sku`. This enables using stubs for prototyping, without having to first create the corresponding go data type.

Query param constants, e.g. `share.ParamUserID`, have the same naming as shared data type fields. However, the `share.Query` func is used for case-insensitive matching, e.g. `?userid=xxx`.

Forms are encoded as flat json. The `share.Params.ToRecords` method converts form post data for use with CSV functions, headers are case-insensitive.

Some forms have well defined fields. In this case the naming convention is to prefix *"Params"*, e.g. `ParamsLoginAttemptPost`.

Avoid duplicating testing of domain model behaviour. Router tests check
- Query param validation
- Parsing request body
- Status codes
