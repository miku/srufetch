# srufetch

Basic SRU endpoint retrieval. For a more serious implementation, see
[yaz-client](https://software.indexdata.com/yaz/doc/yaz-client.html).

Request w/o params will yield an ExplainResponse, e.g. http://sru.k10plus.de/gvk.

Example request: [sru.k10plus.de/gvk?version=1.2...](http://sru.k10plus.de/gvk?version=1.2&operation=searchRetrieve&query=pica.ssg=24,1%20or%20pica.ssg=bbi%20or%20pica.sfk=bub%20or%20pica.osg=bbi&maximumRecords=10&startRecord=10)

```shell
$ srufetch -h
Usage of srufetch:
  -a string
        recordSchema (http://www.loc.gov/standards/sru/recordSchemas/) (default "picaxml")
  -e string
        endpoint (default "https://sru.bsz-bw.de/swb299")
  -ignore-http-errors
        do not fail on HTTP 400 or higher
  -l int
        total limit to retrieve, -1 for no limit (default -1)
  -m int
        maximum records per request (default 10)
  -p duration
        time to sleep between requests (default 100ms)
  -q string
        sru query (default "pica.rvk=\"A*\"")
  -r    randomize the number of records [1, m)
  -s int
        SRU startRecord, zero won't work (default 1)
  -sru-version string
        set SRU version (default "1.1")
  -ua string
        set user agent (default "Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0)")
  -verbose
        increase log output
  -version
        show version
  -x    try to dig out record via regex (XXX: a simple xml.Decode failed)
  -xr string
        (go) regex to parse records (default "(?ms)(<[a-z:]*record(.*?)</[a-z:]*record>)")
```

## Examples

```shell
$ srufetch -x -verbose -e http://sru.k10plus.de/gvk \
           -q "pica.ssg=24,1 or pica.ssg=bbi or pica.sfk=bub or pica.osg=bbi" > data.xml
$ yaz-marcdump -i marcxml -o marc data.xml data.mrc
```

Some endpoints are a bit flaky, use `-r`, `-ignore-http-erros` to mitigate.

```shell
$ srufetch -ignore-http-errors -r -verbose -m 100 -x -a "marc" \
    -e "https://muscat.rism.info/sru" -q "*" \
    -xr "(?ms)(<[a-z:]*recordData(.*?)</[a-z:]*recordData>)"
```
