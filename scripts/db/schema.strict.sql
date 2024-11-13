-- See naming conventions in scripts/db/README.md

-- KSUID is used to generated key sortable unique IDs
-- https://github.com/segmentio/ksuid
-- Consider this article by the ksuid author
-- "if you don’t need to generate unique identifiers
-- in a distributed, stateless fashion, you don’t need UUIDs"
-- https://rbranson.medium.com/why-you-should-avoid-uuids-f3e2936d6ed3
-- However, using the same ID generator everywhere makes sense,
-- and it might be useful for synchronizing data

-- .............................................................................
--
-- config is for global settings
create table config (
	-- term is the unique key for a config value
	term text primary key,
	-- val is the config value
	val text not null default '',
	mod text not null check (mod <> ''),
	-- taxonomies can be used to create a collection of terms, however,
	-- config tables don't necessarily link back to a defined term
	foreign key (term) references term(term)
) strict;

create index config_mod_idx on config(mod);

-- .............................................................................
-- Hash tables help with de-duplication of data
--
-- img hash table
create table img (
	-- hash is computed on the original image
	hash text primary key,
	-- ext to use for a registered image format
	ext text not null,
	-- alt attribute
	alt text not null,
	mod text not null check (mod <> '')
) strict;

create index img_mod_idx on img(mod);

-- addr hash table is append or delete only
create table addr (
	-- hash is calculated on val with normalised white-space
	hash text primary key,
	-- taxonomy defines format fields as per wikipedia
	-- https://en.wikipedia.org/wiki/Address#Format_by_country_and_area
	-- See sample address taxonomies in scripts/db/init.sql
	taxonomy text not null,
	-- val is the address lines separated by newlines,
	-- excluding the recipient name, or other personal details
	val text not null default '',
	mod text not null check (mod <> '')
	-- foreign key (taxonomy) references taxonomy(taxonomy)
) strict;

create index addr_mod_idx on addr(mod);

-- addrtype lookup table lists valid address types, see e.g. order_addr
create table addrtype (
	type text primary key,
	foreign key (type) references order_addr(type)
) strict;

-- .............................................................................
--
-- tag is a lookup table for tags.
-- Applicable to tables that have a corresponding tag meta table
create table tag (
	-- tag is not intended for long strings, e.g. partial HTML or markdown.
	-- Tags may not contain white space, use config table if that is required
	tag text primary key,
	-- mod is useful for querying "most recent tags"
	mod text not null check (mod <> '')
) strict;

create index tag_mod_idx on tag(mod);

-- TODO Trigger to check tag does not contain white space

-- taxonomy is for grouping terms.
-- Inspired by Hugo Taxonomies
-- https://gohugo.io/content-management/taxonomies
-- "Taxonomies are classifications of logical relationships between content",
-- this table defined "logical relationships between data",
-- i.e. config (meta data) table terms (keys).
-- Hugo automatically creates the taxonomies "tags" and "categories",
-- these names are reserved and must not be used in the db
-- https://gohugo.io/content-management/taxonomies/#default-taxonomies
create table taxonomy (
	-- taxonomy is the unique name
	taxonomy text primary key,
	-- descr in short
	descr text not null,
	mod text not null check (mod <> '')
) strict;

create index taxonomy_mod_idx on taxonomy(mod);

-- taxonomy_x_term for associating a term with a taxonomy.
-- This is optional, terms do not have to be part of a taxonomy
create table taxonomy_x_term (
	taxonomy text not null,
	term text not null,
	-- TODO Make term the primary key?
	-- The model.NewFieldLabel func would work better this way,
	-- If there is only one taxonomy per term,
	-- the queries can always return the correct prefix.
	-- The only downside of this is potential duplication of field rows,
	-- in this case field data (e.g. field_opt for size) is likely different?
	primary key (term),
	foreign key (taxonomy) references taxonomy(taxonomy),
	foreign key (term) references term(term)
) strict;

-- taxonomy_x_term_taxonomy_idx index for listing taxonomy terms
create index taxonomy_x_term_taxonomy_idx on taxonomy_x_term(taxonomy);

-- taxonomy_x_term_term_idx index for joining term table
create index taxonomy_x_term_term_idx on taxonomy_x_term(term);

-- term table for listing config terms
create table term (
	-- term is the unique key for a config value,
	-- for use in config meta data tables.
	-- May prefix taxonomy if required to make the term unique
	term text primary key,
	-- descr in short
	descr text not null,
	mod text not null check (mod <> '')
) strict;

create index term_mod_idx on term(mod);

-- taxonomy_x_term_trigger checks for valid term before inserting rows.
-- Corresponding entry in taxonomy table is not required
create trigger taxonomy_x_term_trigger before
insert on taxonomy_x_term begin
select raise(fail, 'invalid term')
where length(new.taxonomy) > 0 and
	length(new.term) > 0 and
	new.term not in (
		select term from term where term = new.term
	);
end;

-- field table for data capture terms
create table field (
	term text primary key,
	-- deflt is a short default value for the term.
	-- If the term is linked to tags then this is the default tag
	deflt text not null default '',
	-- eltag is the HTML element tag name, must be listed in eltag table.
	-- Set empty value if not applicable
	eltag text not null default '',
	-- eltype is the element type attribute, must be listed in eltype table.
	-- Set empty value if not applicable
	eltype text not null default '',
	-- elreq toggles the elements "required" attribute.
	-- Note that SQLite "boolean values are stored as integers"
	-- https://stackoverflow.com/a/22186315/639133
	elreq integer not null check (elreq in (0, 1)) default 0,
	-- TODO Move the idx col to taxonomy_x_term table?
	-- idx for sorting data capture fields.
	-- Optional, fall back to sort on term
	idx integer not null default 0,
	mod text not null check (mod <> ''),
	mod_id text not null check (mod_id <> ''),
	-- foreign key (term) references term(term),
	foreign key (eltag) references eltag(eltag),
	foreign key (eltype) references eltype(eltype)
) strict;

create index field_mod_idx on field(mod);

-- eltag lookup table lists valid element tags.
-- https://developer.mozilla.org/en-US/docs/Web/API/Element/tagName
-- Consider the element might be a Web Component
create table eltag (eltag text primary key) strict;

-- eltag_trigger checks for valid eltag before inserting field rows
-- https://stackoverflow.com/a/42174064/639133
-- TODO Include invalid tag in error message,
-- wait for sqlite 3.47.0, and add context everywhere raise is used
-- https://sqlite.org/forum/forumpost/049f6fde23c88839
-- select raise(fail, 'invalid element tag %s' || new.eltag)
create trigger eltag_trigger before
insert on field begin
select raise(fail, 'invalid element tag')
where new.eltag not in (
		select eltag
		from eltag
		where eltag = new.eltag
	);
end;

-- eltype lookup table lists valid element types
-- https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input#input_types
create table eltype (eltype text primary key) strict;

-- eltype_trigger checks for valid eltype before inserting field rows
create trigger eltype_trigger before
insert on field begin
select raise(fail, 'invalid element type')
where new.eltype not in (
		select eltype
		from eltype
		where eltype = new.eltype
	);
end;

-- term_attr is used to assign attributes to a term.
-- The attribute can then be used as a filter.
-- It's like a tag, but only for terms
create table term_attr (
	term text not null,
	attr text not null,
	primary key (term, attr),
	foreign key (term) references term(term)
) strict;

-- field_elattr lists attribute value pairs to set on a field.
-- It can be used to configure built-in behaviour of HTML Elements,
-- or for attributes that are specific to JavaScript libraries
create table field_elattr (
	term text not null,
	-- elattr is the attribute to set on the HTML element, e.g. "pattern"
	-- https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input#attributes
	-- https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Regular_Expressions
	-- Or "x-validate-rules"
	-- https://github.com/mozey/alpine-util/pull/7
	elattr text not null,
	-- val to set on element attribute
	val text not null default '',
	primary key (term, elattr),
	foreign key (term) references field(term)
) strict;

-- field_elattr_taxonomy_term_idx unique index
create unique index field_elattr_term_idx on field_elattr(term);

-- field_opt for listing field options, e.g. for use with term.eltag="select"
create table field_opt (
	term text not null,
	-- val must not contain space or punctuation characters
	val text not null,
	-- descr in short
	descr text not null,
	primary key (term, val),
	foreign key (term) references field(term)
) strict;

-- field_opt_term_idx index for listing by term
create index field_opt_term_idx on field_opt(term);

-- term_tree is used for hierarchical classification
-- "Many taxonomies are hierarchies and have an intrinsic tree structure"
-- https://en.wikipedia.org/wiki/Taxonomy
create table term_tree (
	term text primary key,
	-- parent is empty for root nodes,
	-- otherwise it's set to the parent term
	parent text not null,
	foreign key (term) references term(term)
) strict;

-- term_tree_taxonomy_term_idx unique index
create unique index term_tree_term_idx on term_tree(term);


-- .............................................................................
--
-- cat is a complete list of items available in this DB
create table cat (
	-- sku (stock-keeping unit) is a unique code.
	-- Must not be updated after the row is created,
	-- the sku is used in URLs for the static pages.
	-- Differs from a serial number, 
	-- or code into an external system,
	-- use cat_config to record additional codes
	sku text not null,
	-- title when displaying catalog items
	title text not null,
	-- descr in short, use config for the latter
	descr text not null,
	-- state must be listed in cat_state.
	-- It's useful for specifying how a catalog item may be used.
	-- 1. The default state is "stock", that means stock this item.
	-- 2. Use "hidden" to hide items from the website, but keep them in the DB,
	-- as opposed to deleting them.
	-- 3. Items with state "discontinued" can't be added to orders,
	-- if state is set to discontinued it can't be changed again.
	-- When deleting items used in order_lines, state is set to discontinued,
	-- otherwise the rows may be deleted from the DB
	-- 4. Items with state "system" are not visible to customer users,
	-- they are used by the system, or admin users.
	-- See comments for cat_qty and order_line tables
	state text not null default 'stock',

	mod text not null check (mod <> ''),
	-- mod_id is the user_id that made the last change
	mod_id text not null check (mod_id <> ''),

	primary key (sku)
) strict;

-- cat_descr_idx for search
create index cat_descr_idx on cat(descr);

create index cat_mod_idx on cat(mod);

-- cat_price is separate, assuming updates are more frequent than for cat
create table cat_price (
	sku text primary key,
	-- price in the smallest possible unit, e.g. cent
	-- This value must be the exclusive price, if multiple VAT rates are used.
	-- All prices in listed the db must be for the default currency
	-- https://github.com/Rhymond/go-money
	-- https://martinfowler.com/eaaCatalog/money.html
	price integer not null check (price >= 0) default 0,
	foreign key (sku) references cat(sku)
) strict;

create table cat_qty (
	sku text not null,
	-- depot is an optional physical location
	depot text not null default '',
	-- qty is the number of items that is available (if applicable)
	qty integer not null check (qty >= 0) default 0,
	primary key (sku, depot),
	foreign key (sku) references cat(sku)
) strict;

-- cat_config meta table, see comments for config table
create table cat_config (
	sku text not null,
	term text not null,
	val text not null default '',
	primary key (sku, term),
	foreign key (sku) references cat(sku),
	-- foreign key (term) references term(term)
	check (sku <> '')
) strict;

-- cat_tag meta table for tagging catalog items.
-- This table is front-matter for generated content,
-- it's the default "tags" taxonomy for rendering static pages
create table cat_tag (
	sku text not null,
	tag text not null,
	primary key (sku, tag),
	foreign key (sku) references cat(sku)
) strict;

-- cat_tag_sku_idx to list all tags for a sku
create index cat_tag_sku_idx on cat_tag(sku);

-- cat_tag_tag_idx to list all skus for a tag
create index cat_tag_tag_idx on cat_tag(tag);

-- cat_state lookup table lists valid catalog states
create table cat_state (
	state text primary key,
	-- custom is a optional user defined label for this state
	custom text not null,
	foreign key (state) references cat(state)
) strict;

-- cat_state_trigger checks for valid state before inserting cat rows
create trigger cat_state_trigger before
insert on cat begin
select raise(fail, 'invalid catalog state')
where new.state not in (
		select state
		from cat_state
		where state = new.state
	);
end;

-- variant groups catalog items that are variants.
-- Using variants is optional.
-- A separate static page is rendered for each variant,
-- and a dropdown or list links the variant pages
--
-- SKU must be unique for each variant, and group_id might be a prefix.
-- E.g.
-- group_id="running-shoe"
-- sku="running-shoe-42-grey"
--
-- TODO Taxonomy "variants" with terms "variant_attr_x" and "variant_val_x"?
-- E.g.
-- variant_attr_1="size", variant_val_1="42"
-- variant_attr_2="colour", variant_val_2="grey"
create table variant (
	-- group_id is the unique ID for a group of variants
	group_id text not null default '',
	-- sku is the unique ID for each variant
	sku text not null,
	-- idx to order variants for display.
	-- Optional, fall back to sort on sku
	idx integer not null default 0,
	primary key (group_id, sku),
	foreign key (sku) references cat(sku)
) strict;

create table cat_img (
	sku text not null,
	hash text not null,
	-- descr for the image used in this context,
	-- i.e. the same img can also be linked to other ref tables
	descr text not null default '',
	-- idx for sorting images linked to a sku.
	-- Sort query result by idx and mod asc,
	-- the first row index 0 is the default image.
	-- When setting the default image,
	-- set idx for all other rows to greater than 0
	idx integer not null default 0,
	primary key (sku, hash),
	foreign key (hash) references img(hash)
) strict;

-- TODO Trigger on cat_img, hash must exist in img table

-- .............................................................................
--
-- user table
create table user (
	-- user_id is required because users might want to change their email
	user_id text primary key,
	-- email might be shared for different roles.
	-- Note that "there is a restriction in RFC 2821 on the length of an
	-- address in MAIL and RCPT commands of 256 characters"
	-- https://www.rfc-editor.org/errata/eid1690
	email text not null,
	-- username is used for multiple identities,
	-- it's not required if email is unique.
	-- The option to add username is only available when
	-- registering the same email for the second time?
	username text not null default '',
	-- descr in short.
	-- Used as a label (instead of email) when displaying users if not empty
	descr text not null default '',
	-- role for this user, e.g. "admin" or "customer"
	-- Bulk import (upsert) on this table is allowed,
	-- but not for "admin" users?
	role text not null default 'customer',
	-- verified timestamp is not set for new users,
	-- it is set the first time a user verifies, and then remains set
	verified integer not null check (verified >= 0) default 0,
	-- disabled timestamp disables a verified user,
	-- instead of permanently deleting the user.
	-- This might be useful if the account misbehaves,
	-- and/or to preserve user activity.
	-- Disabled users can't create sessions
	disabled integer not null check (disabled >= 0) default 0,
	mod text not null check (mod <> ''),
	-- TODO mod_id for this table?
	foreign key (role) references role(role)
) strict;

-- user_email_username_idx unique index
create unique index user_email_username_idx on user(email, username);

create index user_mod_idx on user(mod);

-- user_config meta table
create table user_config (
	user_id text not null,
	term text not null,
	val text not null default '',
	primary key (user_id, term),
	foreign key (user_id) references user(user_id)
	-- foreign key (term) references term(term)
) strict;

-- user_tag
create table user_tag (
	user_id text not null,
	tag text not null,
	primary key (user_id, tag),
	foreign key (user_id) references user(user_id)
) strict;

-- user_tag_user_id_idx to list all tags for a user_id
create index user_tag_user_id_idx on user_tag(user_id);

-- user_tag_tag_idx to list all user_ids for a tag
create index user_tag_tag_idx on user_tag(tag);

-- session table for keeping track of active sessions.
-- Users can have one active session per email and username,
-- and multiple usernames (and roles) per email.
-- An asymmetrically signed JWT is used as a bearer token with requests.
-- Claims (e.g. user_id) are encoded in the JWT.
-- Claims are not secret since the session may be decoded with a public key,
-- but tamper proof because the JWT is signed with a private key.
-- Reads on this table are minimized since the private key is secret,
-- i.e. invalid or expired tokens do not require db reads
-- https://github.com/shopd/shopd-issues/issues/33
-- Each unique email and username combo has one session,
-- the relationship between the user and session tables is one to one.
-- To logout a user, reset the corresponding session row
create table session (
	-- user_id is unique per email and username,
	-- and each user has one role
	user_id text primary key,
	-- otp is used to ensure only the most recent loginTokens
	-- can be used to create a new accessToken.
	-- Otherwise any previous loginToken would work,
	-- if it was signed with the same private key.
	-- Can only be used to generate an accessToken once,
	-- reset the otp when the session is verified
	-- https://en.wikipedia.org/wiki/One-time_password
	otp text not null,
	-- verified is used with login links.
	-- Initially the session is not verified.
	-- Clicking the link verifies the session
	verified integer not null check (verified >= 0) default 0,
	-- attempts is the number of pending login attempts.
	-- It's incremented each time the user requests a session,
	-- after a specified number of tries the admin is notified,
	-- and the system is prevented from spamming the email address.
	-- Also incremented when the user types tries to verify an OTP,
	-- to prevent brute force guessing
	attempts integer not null check (attempts >= 0) default 0,
	-- mod date
	mod text not null check (mod <> ''),
	-- mod_id not applicable,
	-- only system can edit this table
	foreign key (user_id) references user(user_id)
) strict;

create index session_mod_idx on session(mod);

-- session_act is session activity.
-- The user table is one to one with session,
-- for config use the user_config table.
-- This table is append only, that means
-- insert, select and delete only, no updates
create table session_act (
	user_id text not null,
	-- msg for this activity entry
	msg text not null,
	mod text primary key,
	-- mod_id not applicable,
	-- only system can edit this table
	foreign key (user_id) references session(user_id)
) strict;

create index session_act_mod_idx on session_act(mod);

create table role (
	role text primary key
) strict;

-- user_role_trigger checks for valid role before inserting rows.
-- The user role must be listed in the role table
create trigger user_role_trigger before
insert on user begin
select raise(fail, 'invalid user role')
where new.role not in (
		select role
		from role
		where role = new.role
	);
end;

-- role_perm to configure custom roles.
-- For example the "sync" user role might consist of the permissions
-- the "ExportOrders", and "ExportCatalog", or "ImportCatalog" etc.
-- Roles and permission are not intended to be composable,
-- they're created bearing in mind specific functionality and conventions.
-- This makes the system less flexible, but hopefully easier to understand.
-- The basic convention is this, the default role is "customer",
-- and routes starting with "/admin" is only for "admin" users.
-- Admin user can still place orders as a customer, i.e. non-admin routes.
-- Potentially admin users might place orders on behalf of other users,
-- or view the site as another user.
-- The "custom" and "admin" roles do not require entries in this table.
-- Other roles only have the permission listed in here, i.e. whitelist
create table role_perm (
	role text not null,
	perm text not null,
	-- path is the api path
	path text not null,
	primary key (role, perm, path),
	foreign key (role) references role(role)
	-- TODO Trigger for perm but not path?
) strict;

-- .............................................................................
--

-- TODO Create a msg domain model,
-- to start just send email in the request handler.
-- Then refactor to buffer to NATS later

-- .............................................................................
--
-- orders table
-- "order" is a reserved word, avoid having to type escape chars
create table orders (
	order_id text primary key,
	-- order_no is initially set to empty string,
	-- therefore it can't have a unique index
	order_no text not null,
	state text not null,
	-- notes for this order.
	-- May be written by the customer, or by an admin during processing
	notes text not null,
	-- user_id that created this order
	user_id text not null,
	-- paid is set if the order has been paid in full,
	-- for partial payments see the order_tran table
	paid integer not null check (paid in (0, 1)) default 0,
	mod text not null check (mod <> ''),
	-- mod_id is the user_id that made the last change
	mod_id text not null check (mod_id <> '')
) strict;

create index orders_mod_idx on orders(mod);

create table order_tax (
	order_id text not null,
	-- order_line_id if applicable, otherwise empty.
	-- For example, the tax calculation method could be "basket"
	-- https://github.com/shopd/shopd-issues/issues/70
	order_line_id text not null,
	-- fixed is set if tax was calculated as a fixed amount
	fixed integer not null check (fixed in (0, 1)) default 0,
	-- tax is the calculated amount before rounding.
	-- Floating point sum in SQLite,
	-- better to apply the "scaling factor" (rounding) in code
	-- https://g.co/gemini/share/33caaf098314
	-- Should order tax lines be summed as float or int?
	-- https://g.co/gemini/share/cfa33d2d719d
	tax real not null check (tax >= 0) default 0,
	-- pct is non-zero if a percentage tax was applied,
	-- e.g. use 1500 for 15% VAT
	pct integer not null check (pct >= 0) default 0,
	-- mod is the timestamp when this tax line was added.
	-- Tax lines can't be edited, only add or delete is allowed
	mod text not null check (mod <> ''),
	primary key (order_id, order_line_id, mod),
	foreign key (order_id) references orders(order_id)
) strict;

-- order_config meta table
create table order_config (
	order_id text not null,
	-- order_line_id if applicable, otherwise empty
	order_line_id text not null,
	term text not null,
	val text not null default '',
	primary key (order_id, order_line_id, term),
	foreign key (order_id) references orders(order_id)
	-- foreign key (term) references term(term)
) strict;

-- order_line meta table
create table order_line (
	-- order_line_id can be used to sort order lines by creation timestamp
	order_line_id text primary key,
	-- order_id these lines belong to
	order_id text not null,
	-- state overrides orders.state if not empty.
	-- Note that null is not used in the db. However,
	-- the code may check for null to toggle bulk update cols.
	-- Therefore most shared data struct fields support null
	state text not null,
	-- sku is the unique catalog item.
	-- System processes and admin users can add "system" skus to an order,
	-- they are useful for things like discounts, coupons, vouchers, etc.
	-- TODO Consider discounts, coupons, vouchers, etc.
	-- "A coupon grants you a discount on your order.
	-- A voucher, on the other hand, is considered a monetary substitute,
	-- which is determined by the amount stated on the voucher"
	sku text not null,
	-- price in the smallest possible unit, e.g. cents
	price integer not null default 0,
	-- qty is the number of items
	qty integer not null check (qty >= 1) default 0,
	foreign key (order_id) references orders(order_id)
	-- foreign key (sku) references cat(sku)
) strict;

-- order_act is order activity, state history, and admin notes
create table order_act (
	order_id text not null,
	-- order_line_id if applicable, otherwise empty
	order_line_id text not null,
	-- state is the order state at a point in time
	state text not null,
	-- msg for this activity entry
	msg text not null,
	-- user_id that created this activity line,
	-- zero if the row was created by the system
	user_id text not null,
	-- admin is set if this activity entry is visible to admin users only.
	-- Useful for adding admin only notes in the msg col
	admin integer not null check (admin in (0, 1)) default 0,
	-- mod records when the order activity occurred
	mod text primary key,
	-- mod_id not applicable, this table is append only,
	-- that means insert, select and delete only, no updates
	foreign key (order_id) references orders(order_id)
) strict;

create index order_act_mod_idx on order_act(mod);

-- order_tag bridge table.
-- Unlike cat_tag, the values in this table is not front-matter,
-- i.e. it is not used for rendering static pages
create table order_tag (
	order_id text not null,
	-- order_line_id if applicable, otherwise empty
	order_line_id text not null,
	tag text not null,
	primary key (order_id, order_line_id, tag),
	foreign key (order_id) references orders(order_id)
) strict;

-- order_tag_order_id_idx to list all tags for an order_id
create index order_tag_order_id_idx on order_tag(order_id);

-- order_tag_tag_idx to list all order_ids for a tag
create index order_tag_tag_idx on order_tag(tag);

-- order_state lookup table lists valid order states
create table order_state (
	state text primary key,
	foreign key (state) references orders(state)
) strict;

-- order_state_trigger checks for valid state before inserting order rows
create trigger order_state_trigger before
insert on orders begin
select raise(fail, 'invalid order state')
where new.state not in (
		select state
		from order_state
		where state = new.state
	);
end;

create table order_addr (
	order_id text not null,
	-- type of address, e.g. delivery, billing, etc
	-- Order may only have one address per type
	type text not null,
	-- hash of the address
	hash text not null,
	primary key (order_id, type),
	foreign key (hash) references addr(hash)
) strict;

-- order_addr_trigger checks for valid type before inserting order_addr
create trigger order_addr_trigger before
insert on order_addr begin
select raise(fail, 'invalid address type')
where new.type not in (
		select type
		from addrtype
		where type = new.type
	);
end;


-- order_tran links transactions to an order
create table order_tran (
	order_id text not null,
	tran_id text not null,
	primary key (order_id, tran_id),
	foreign key (tran_id) references tran(tran_id)
) strict;


-- .............................................................................

-- tran table for recording payments, credit notes, etc
create table tran (
	tran_id text primary key,
	-- account_id is empty if not applicable
	account_id text not null,
	-- state indicated if the tran was successful, see comments on
	-- other state machines in go/hooks/README/OrderState.png
	state text not null,
	descr text not null,
	-- amount in the smallest possible unit, e.g. cents
	amount integer not null check (amount >= 0),
	-- currency code for the transaction
	currency text not null,
	-- user_id that made this transaction
	user_id text not null,
	mod text not null check (mod <> ''),
	-- mod_id not applicable,
	-- only system can edit this table
	foreign key (account_id) references account(account_id)
) strict;

create index tran_mod_idx on tran(mod);

-- tran_config meta table, e.g.
-- "method=cash", "method=card", "processor=stripe", "ref=foo"
create table tran_config (
	tran_id text not null,
	term text not null,
	val text not null default '',
	user_id text not null,
	primary key (tran_id, term),
	foreign key (tran_id) references tran(tran_id)
	-- foreign key (term) references term(term)
) strict;

-- tran_tag bridge table
create table tran_tag (
	tran_id text not null,
	tag text not null,
	primary key (tran_id, tag),
	foreign key (tran_id) references tran(tran_id)
) strict;

-- tran_tag_tran_id_idx to list all tags for an tran_id
create index tran_tag_tran_id_idx on tran_tag(tran_id);

-- tran_tag_tag_idx to list all tran_ids for a tag
create index tran_tag_tag_idx on tran_tag(tag);

-- TODO account might be overkill for now,
-- but create it now for future reference.
-- In the meantime tran table is good enough?
-- All accounts in here are for users with customer role,
-- hosted site subscriptions are accounts in a back office store?
create table account (
	account_id text primary key,
	descr text not null
) strict;

create table account_x_user (
	account_id text not null,
	user_id text not null,
	primary key (account_id, user_id),
	foreign key (account_id) references account(account_id)
) strict;


-- .............................................................................

-- vat table lists the default vat rate by country code
create table vat (
	-- country is the ISO 3166-1 alpha-3 country code
	-- https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes
	country text primary key,
	-- pct is the percentage tax to apply,
	-- e.g. use 1500 for 15% VAT
	pct integer not null check (pct >= 0) default 0,
	mod text not null check (mod <> ''),
	mod_id text not null check (mod_id <> '')
) strict;

-- tax table may be used to specify additional tax,
-- or override the default vat rate for a country code.
-- TODO Using this table requires cat_price to use exclusive price,
-- otherwise it would be impossible to calculate with multiple tax lines?
create table tax (
	country text not null,
	-- tag to match skus on cat_tag table
	tag text not null,
	-- sku to match on cat table
	sku text not null,
	-- vat is set to override the default vat rate for a country code.
	-- Specificity of the override is country, tag, and then sku.
	-- This tax row is an additional tax line if vat override is not set
	vat integer not null check (vat in (0, 1)) default 0,
	-- pct is the percentage tax to apply,
	-- e.g. use 1500 for 15% VAT.
	-- Some countries may have multiple VAT rates
	-- https://g.co/gemini/share/978a8127426d
	-- South-Africa only has a standard, or zero, VAT rate
	-- https://g.co/gemini/share/37f88b439988
	pct integer not null check (pct >= 0) default 0,
	-- tax is a fixed amount of tax in the smallest unit, e.g. cents.
	-- Use the value in this col instead of pct if non-zero
	tax integer not null check (tax >= 0) default 0,
	-- descr for this tax line,
	-- e.g. "Essential goods" "Luxury goods", "Services".
	-- Leave empty to use default description, i.e. "Tax"
	descr text not null,
	mod text not null check (mod <> ''),
	mod_id text not null check (mod_id <> ''),
	primary key (country, tag, sku),
	foreign key (country) references vat(country)
) strict;


-- .............................................................................

-- TODO price override tables have the same structure as discount tables?

create table discount (
	discount_id text primary key,
	-- pct is the percentage discount to apply, or zero if not applicable
	pct integer not null check (pct >= 0) default 0,
	-- discount is a fixed discount, instead of percentage.
	-- Use the value in this col instead of pct if non-zero
	discount integer not null check (discount >= 0) default 0,
	-- descr to describe what this discount is for
	descr text not null,
	mod text not null check (mod <> ''),
	mod_id text not null check (mod_id <> '')
) strict;

-- discount_country if the discount is for specified countries
create table discount_country (
	discount_id text not null,
	country text not null,
	primary key (discount_id, country),
	foreign key (discount_id) references discount(discount_id)
) strict;

-- discount_tag if the discount is for specified tags in cat_tag
create table discount_tag (
	discount_id text not null,
	tag text not null,
	primary key (discount_id, tag),
	foreign key (discount_id) references discount(discount_id)
) strict;

-- discount_sku if the discount is for specified skus
create table discount_sku (
	discount_id text not null,
	sku text not null,
	primary key (discount_id, sku),
	foreign key (discount_id) references discount(discount_id)
) strict;

-- discount_user if the discount is for specified users
create table discount_user (
	discount_id text not null,
	user_id text not null,
	primary key (discount_id, user_id),
	foreign key (discount_id) references discount(discount_id)
) strict;

-- discount_range if the discount applies for a date time range
create table discount_range (
	discount_id text not null,
	-- start date time
	start text not null,
	-- end date time
	end text not null,
	primary key (discount_id, start, end),
	foreign key (discount_id) references discount(discount_id)
) strict;

-- discount_opt lists discounts that may optionally be applied by admin users
-- on checkout, or afterward creating the order but before receiving payment?
create table discount_opt (
	discount_id text primary key,
	foreign key (discount_id) references discount(discount_id)
) strict;
