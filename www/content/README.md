# www/content

Site content

Conventions
- Each file contains one component, and the function name must be `Index`
- Content components all receive the same param, the `Content` view model
- The index files renders the dir path, e.g. `login/index.templ` renders `GET /login`
- Files not named index will add another path segment, e.g. `store/apples.templ` renders `GET /store/apples/`, and `store/oranges.templ` renders `GET /store/oranges/`

**TODO** See comments in `www/view/README.md` re. overrides
