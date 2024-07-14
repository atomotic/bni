# BNI – Bibliografia Nazionale Italiana

Download Unimarc XML files from [BNIweb](http://bni.bncf.firenze.sbn.it/bniweb/menu.jsp) and convert to Parquet with Duckdb.

The Parquet dump is available here https://atomotic.github.io/bni/bni.parquet (70M) and can be used with
[DuckDB Shell](https://shell.duckdb.org/#queries=v0,select-id%2Cisbn%2Ctitle-from-'https%3A%2F%2Fatomotic.github.io%2Fbni%2Fbni.parquet'-where-title-like-'%25esoteri%25'-limit-20~)

## Steps

The following steps are available inside the [Justfile](https://github.com/casey/just)

Scrape all XML urls (tools needed: [pup](https://github.com/ericchiang/pup) and [sd](https://github.com/chmln/sd))

```bash
curl -s "http://bni.bncf.firenze.sbn.it/bniweb/menu.jsp" \
    | pup 'a attr{href}' \
    | grep elenco_fasc \
    | sd "&amp;" "&" \
    | sd "elenco_fasc" "scaricaxml" \
    > links.txt
```

Download all XML files (tool needed: [wcurl](https://github.com/Debian/wcurl))

```bash
parallel wcurl --curl-options="--remote-header-name" "http://bni.bncf.firenze.sbn.it/bniweb/{}" :::: links.txt
mkdir xml
move *.xml xml/
```

Load all XML files to DuckDB (tools needed: [Go](https://golang.org) and [gnu parallel](https://www.gnu.org/software/parallel/))

```bash
go build
parallel -j1 ./bni {} ::: xml/*.xml
```

Export from DuckDB to Parquet

```
duckdb bni.ddb "copy bni to bni.parquet (format parquet);"
```

Size comparison

```
du -h bni.ddb bni.parquet
1.2G    bni.ddb
67M     bni.parquet
```

## Example query

```
duckdb
```

The schema: `data` contains the full Unimarc record converted to JSON

```
DESCRIBE SELECT * FROM 'https://atomotic.github.io/bni/bni.parquet';
┌─────────────┬─────────────┬─────────┬─────────┬─────────┬─────────┐
│ column_name │ column_type │  null   │   key   │ default │  extra  │
│   varchar   │   varchar   │ varchar │ varchar │ varchar │ varchar │
├─────────────┼─────────────┼─────────┼─────────┼─────────┼─────────┤
│ id          │ VARCHAR     │ YES     │         │         │         │
│ isbn        │ VARCHAR     │ YES     │         │         │         │
│ title       │ VARCHAR     │ YES     │         │         │         │
│ data        │ VARCHAR     │ YES     │         │         │         │
│ source      │ VARCHAR     │ YES     │         │         │         │
└─────────────┴─────────────┴─────────┴─────────┴─────────┴─────────┘
```

```
D .mode line
D SELECT id,title,isbn,source FROM 'https://atomotic.github.io/bni/bni.parquet' WHERE title LIKE '%biblioteco%' LIMIT 5;

    id = USM1959877
 title = Biblioteche e biblioteconomia
  isbn = 9788843075294
source = xml/Monografie201503.xml

    id = PAV0095007
 title = I fondamenti della biblioteconomia
  isbn = 9788870758474
source = xml/Monografie201601.xml

    id = SBT0014568
 title = Conferimento della laurea magistrale ad honorem in scienze archivistiche e biblioteconomiche a Michele Casalini
  isbn = 9788864538822
source = xml/Monografie201904.xml

    id = MOD1738924
 title = Guida alla biblioteconomia moderna
  isbn = 9788893574013
source = xml/Monografie202204.xml

    id = SBT0045209
 title = Principi, approcci e applicazioni della biblioteconomia comparata
  isbn = 9788855186063
source = xml/Monografie202301.xml

```
