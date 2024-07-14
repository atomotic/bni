default:
    just -l

links:
    curl -s "http://bni.bncf.firenze.sbn.it/bniweb/menu.jsp" \
        | pup 'a attr{href}' \
        | grep elenco_fasc \
        | sd "&amp;" "&" \
        | sd "elenco_fasc" "scaricaxml" \
        > links.txt

downloads:
    parallel wcurl --curl-options="--remote-header-name" "http://bni.bncf.firenze.sbn.it/bniweb/{}" :::: links.txt
    mkdir xml
    move *.xml xml/

load:
    go build
    parallel -j1 ./bni {} ::: xml/*.xml

export:
    duckdb bni.ddb "copy bni to test.parquet (format parquet);"
