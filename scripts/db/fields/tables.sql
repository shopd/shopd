-- "How to get a list of column names on Sqlite3 database?"
-- https://stackoverflow.com/a/71413249/639133
select distinct name TableName
from sqlite_master where type = 'table' and TableName not like 'sqlite_%'
order by TableName
