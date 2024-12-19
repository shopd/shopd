# www/view

[View models](https://templ.guide/core-concepts/view-models/)

Naming convention is to have files corresponding to the last path element, e.g. `login.go` for the paths `/login` and `/api/login`. 

Content models are named the same as the path (`type Login`), and view models for the Hypermedia API append the method (`type LoginPost`).

The view model is not the HTTP request payload, it's the data used for server side rendering of the templ components.

Only for use in route handlers with the templ components in `/www`, generally the backend code must make use of domain models and shared data types.

View models may embed shared data types.

**TODO** Domains may override existing components, but not the corresponding view models. Custom components may define their own view models. Implementing this would require two things
- Ideally using templ cmd, e.g. `templ generate --override $DOMAIN_DIR/www` to override existing components
- Generate init code in the router package that calls `NewCustomRouter(r *gin.Engine)` to add custom routes defined in the domain dir
