package share

const GET = "GET"
const PATCH = "PATCH"
const POST = "POST"
const PUT = "PUT"
const DELETE = "DELETE"

// .............................................................................
// Define query params here (if there is no matching const in schema.params.go).
// Camel case is used for consistent naming with data type structs,
// fields must be public and therefore start with uppercase.
// Go templating expects public fields.

const ParamEnv = "Env"
