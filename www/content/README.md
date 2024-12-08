# www/content

Static site content

Conventions
- Each file contains one component, and the function name must be `Index`
- The generated HTML file will have the same name, e.g. `login/index.templ` will generate `login/index.html`
- File not named index will add another path segment, e.g. `store/apples.templ` will generate `store/apples/index.html`, and `store/oranges.templ` will generate `store/oranges/index.html`

The conventions above make it possible to create domain specific override files
