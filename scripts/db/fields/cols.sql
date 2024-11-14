-- "How to get a list of column names on Sqlite3 database?"
-- https://stackoverflow.com/a/71413249/639133
with tables as (select name tablename, sql 
from sqlite_master where type = 'table' and tablename not like 'sqlite_%')
-- select fields.name, fields.type, tablename
select distinct fields.name ColName
from tables cross join pragma_table_info(tables.tablename) fields
order by ColName
