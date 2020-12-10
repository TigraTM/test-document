package migrations

var Schema = `
create table if not exists documents (
    name text,
    date text,
    number integer,
    sum text
)`