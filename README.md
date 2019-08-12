# srufetch

Basic SRU endpoint retrieval.

Request w/o params will yield an ExplainResponse, e.g. http://sru.k10plus.de/gvk.

Example request: [sru.k10plus.de/gvk?version=1.2...](http://sru.k10plus.de/gvk?version=1.2&operation=searchRetrieve&query=pica.ssg=24,1%20or%20pica.ssg=bbi%20or%20pica.sfk=bub%20or%20pica.osg=bbi&maximumRecords=10&startRecord=10)

```
$ srufetch -h
Usage of ./srufetch:
  -e string
        endpoint (default "http://sru.k10plus.de/gvk")
  -l int
        total limit to retrieve, -1 for no limit (default -1)
  -m int
        maximum records per request (default 10)
  -q string
        sru query (default "pica.ssg=24,1 or pica.ssg=bbi or pica.sfk=bub or pica.osg=bbi")
  -s int
        SRU startRecord, zero won't work (default 1)
  -verbose
        increase log output
  -version
        show version
  -x    try to dig out record via regex (XXX: a simple xml.Encode failed)

$ srufetch -x -verbose -e http://sru.k10plus.de/gvk -q "pica.ssg=24,1 or pica.ssg=bbi or pica.sfk=bub or pica.osg=bbi" > data.xml
$ yaz-marcdump -i marcxml -o marc data.xml data.mrc
```
