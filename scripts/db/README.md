# scripts/db


## Naming Conventions

All tables must use [strict typing mode](https://www.sqlite.org/stricttables.html).

Consider the rules for generated code as per [renaming fields](https://docs.sqlc.dev/en/stable/howto/rename.html): *"Struct field names are generated from column names using a simple algorithm: split the column name on underscores and capitalize the first letter of each part..."*

**Table names** must preferably be a single short word, e.g. *"orders"*. From the docs linked above: *"By default, the struct name will be the singular version of the table name"*, however `emit_exact_table_names` is set to true in the `sqlc.yaml`

**One-to-many (meta or EAV) table names** prefixes the parent, e.g. *"cat_config"*, or *"user_hist"*. See [EAV](https://en.wikipedia.org/wiki/Entity%E2%80%93attribute%E2%80%93value_model)

The characters *"_x_"* delimit **Many-to-many (bridge) table names**. Tables in the bridge table names must be ordered alphabetically, e.g. *"taxonomy_x_term"* or *"session_x_user"*.

**Col names** must be snake_case and all lowercase, e.g. *"user_hist_id"* ID cols must be named *"${TABLE}_id"*, e.g. *"user_hist.user_hist_id"*, or *"user.user_id"*. Foreign keys must be named the same as the referenced col, e.g. *"cat.user_id"*. Tables with an ID, or single col primary key, must define it first.

Underscore delimits **index and trigger names**. They start by prefixing the parent table name, followed by column names ordered alphabetically, and ends with either *"_idx"* or *"_trigger"*, e.g *"tag_tag_idx"*, or *"msg_format_trigger"*.

**Terms** use snake_case e.g. *cat_config.term = "some_thing"*. 

**Taxonomies** also use snake_case, e.g. *"address_za"*. Terms may prefix the taxonomy, e.g. *"address_za_street"*

**Hash tables** (e.g. *"img"* and *"address"*) may be linked from any config table, the convention is for naming the term is *"hash_${TABLE}_${DESCR}"*, e.g. *cat_config.term = "hash_img_default"*, or *user_config.term = "hash_address_home"*.

Tables with the same name prefix are grouped together in the schema definition. Group statements in this order
- `create table` for reference data
- `create index` and `trigger` on reference data tables
- Create related bridge or meta tables
- Indexes and triggers on related tables

Table with the following cols must always include them last in this order
- `mod` (modified timestamp)
- `del` (deleted timestamp, set to *0* if not deleted)

When syncing data (e.g. updating FTS tables), the mod col may be used as a **pagination** token. Values for this col are unique, and loosely but not exactly sortable by order of creation


## FTS5

**TODO** Store uses FTS, admin pages use more precise queries. That means admin advanced search will match rows even if the FTS table was not, for whatever reason, populated with that data. Domain models are responsible for keeping the FTS tables in sync. There must be a way to rebuild FTS tables from scratch

Notes from the [fts5 page](https://www.sqlite.org/fts5.html)

"It is an error to add types, constraints or PRIMARY KEY declarations to a CREATE VIRTUAL TABLE statement used to create an FTS5 table"

"FTS5 table may be populated using INSERT, UPDATE or DELETE statements like any other table"

"If using the MATCH or = operators, the expression to the left of the MATCH operator is usually the name of the FTS5 table", or use "table-valued function syntax", e.g.
```sql
select * from email('fts5')
```

"By default, FTS5 full-text searches are case-independent"

"A description of the available auxiliary functions, and more details regarding configuration of the special **rank** column, are [available below](https://www.sqlite.org/fts5.html#_auxiliary_functions_)"

"advanced searches are requested by providing a more complicated **FTS5 query string**", see [full-text query syntax](https://www.sqlite.org/fts5.html#full_text_query_syntax)

