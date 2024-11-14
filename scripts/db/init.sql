-- "Foreign key constraints are disabled by default",
-- disable them here explicitly anyway
-- https://sqlite.org/foreignkeys.html
-- Foreign keys are not used (by the app) to enforce relationships,
-- however, they are included in the schema for diagram visualization
pragma foreign_keys = off;

insert into taxonomy(taxonomy, descr, mod) values
-- See comments for Address table, and formats on Wikipedia
-- https://en.wikipedia.org/wiki/Address#Format_by_country_and_area
("address_it", "Address format for Italy", "000pt58M8fYM8MzqlOmoPyu0lbE"),
("address_za", "Address format for South Africa", "000pt58M8fYM8MzqlOmoPyu0lbE");

insert into eltag(eltag) values
-- Empty string is the default for term.eltag
(""),
("combobox"), -- component
("input"),
("select"),
("textarea");

insert into eltype(eltype) values
-- Empty string is the default for term.eltype
(""),
("email"),
("text");

insert into term(term, descr, mod) values
("address_za_street", "Street address or PO box", "000pt58M8fYM8MzqlOmoPyu0lbE"),
("address_za_suburb", "Suburb", "000pt58M8fYM8MzqlOmoPyu0lbE"),
("address_za_city", "City", "000pt58M8fYM8MzqlOmoPyu0lbE"),
("address_za_province", "Province", "000pt58M8fYM8MzqlOmoPyu0lbE"),
("address_za_postcode", "Postcode", "000pt58M8fYM8MzqlOmoPyu0lbE"),

("address_it_street", "Street name and number", "000pt58M8fYM8MzqlOmoPyu0lbE"),
("address_it_building", "Building, floor, and/or apartment number", "000pt58M8fYM8MzqlOmoPyu0lbE"),
("address_it_postbox", "Post office box number", "000pt58M8fYM8MzqlOmoPyu0lbE"),
("address_it_postcode", "Postcode, town and province abbreviation", "000pt58M8fYM8MzqlOmoPyu0lbE");

insert into taxonomy_x_term(taxonomy, term) values
("address_za", "address_za_street"),
("address_za", "address_za_suburb"),
("address_za", "address_za_city"),
("address_za", "address_za_province"),
("address_za", "address_za_postcode"),

("address_it", "address_it_street"),
("address_it", "address_it_building"),
("address_it", "address_it_postbox"),
("address_it", "address_it_postcode");

insert into field(term, deflt, eltag, eltype, elreq, idx, mod, mod_id) values
("address_za_street", "", "input", "text", 1, 0, "000pt58M8fYM8MzqlOmoPyu0lbE", "s"),
("address_za_suburb", "", "input", "text", 0, 1, "000pt58M8fYM8MzqlOmoPyu0lbE", "s"),
("address_za_city", "", "input", "text", 1, 2, "000pt58M8fYM8MzqlOmoPyu0lbE", "s"),
("address_za_province", "", "input", "text", 1, 3, "000pt58M8fYM8MzqlOmoPyu0lbE", "s"),
("address_za_postcode", "", "input", "text", 1, 4, "000pt58M8fYM8MzqlOmoPyu0lbE", "s"),

("address_it_street", "", "input", "text", 1, 0, "000pt58M8fYM8MzqlOmoPyu0lbE", "s"),
("address_it_building", "", "input", "text", 1, 1, "000pt58M8fYM8MzqlOmoPyu0lbE", "s"),
("address_it_postbox", "", "input", "text", 1, 2, "000pt58M8fYM8MzqlOmoPyu0lbE", "s"),
("address_it_postcode", "", "input", "text", 1, 3, "000pt58M8fYM8MzqlOmoPyu0lbE", "s");

insert into cat_state(state, custom) values
("stock", ""),
("hidden", ""),
("discontinued", ""),
("system", "");

-- TODO States listed here must correspond to the state machine that's used,
-- i.e. it's a subset of all the possibly states
insert into order_state(state) values
("cart"),
("pending"),
("confirmed"),
("reversed"),
("complete");

insert into addrtype(type) values
("billing"),
("delivery");

insert into role(role) values
("admin"),
("customer"),
("webhook");

insert into vat(country, pct, mod, mod_id) values
("ZAF", "1500", "000pt58M8fYM8MzqlOmoPyu0lbE", "s");
