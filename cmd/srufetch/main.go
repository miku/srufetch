// Fetch SRU endpoint, raw acquisition.
//
// Request w/o params will yield an ExplainResponse, e.g.
// http://sru.k10plus.de/gvk.
//
// Example request:
//
// http://sru.k10plus.de/gvk?version=1.2&operation=searchRetrieve&query=pica.ssg=24,1%20or%20pica.ssg=bbi%20or%20pica.sfk=bub%20or%20pica.osg=bbi&maximumRecords=10&startRecord=10
//
package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/sethgrid/pester"
)

var (
	startRecord      = flag.Int("s", 1, "SRU startRecord, zero won't work")
	maximumRecords   = flag.Int("m", 10, "maximum records per request")
	endpoint         = flag.String("e", "https://sru.bsz-bw.de/swb299", "endpoint")
	verbose          = flag.Bool("verbose", false, "increase log output")
	limit            = flag.Int("l", -1, "total limit to retrieve, -1 for no limit")
	recordRegex      = flag.Bool("x", false, "try to dig out record via regex (XXX: a simple xml.Encode failed)")
	query            = flag.String("q", `pica.rvk="A*"`, "sru query")
	recordSchema     = flag.String("a", "picaxml", "recordSchema (http://www.loc.gov/standards/sru/recordSchemas/)")
	showVersion      = flag.Bool("version", false, "show version")
	userAgent        = flag.String("ua", "Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0)", "set user agent")
	ignoreHTTPErrors = flag.Bool("ignore-http-errors", false, "do not fail on HTTP 400 or higher")

	Version   string
	BuildTime string
)

// SearchRetrieveResponse was generated 2019-07-17 14:05:42 by tir on sol.
type SearchRetrieveResponse struct {
	XMLName         xml.Name `xml:"searchRetrieveResponse"`
	Text            string   `xml:",chardata"`
	Zs              string   `xml:"zs,attr"`
	Version         string   `xml:"version"`         // 1.2
	NumberOfRecords int      `xml:"numberOfRecords"` // 151502
	Records         struct {
		Text   string `xml:",chardata"`
		Record []struct {
			Text          string `xml:",chardata"`
			RecordSchema  string `xml:"recordSchema"`
			RecordPacking string `xml:"recordPacking"` // xml, xml, xml, xml, xml, ...
			RecordData    struct {
				Text   string `xml:",chardata"`
				Record struct {
					Text         string `xml:",chardata"`
					Xmlns        string `xml:"xmlns,attr"`
					Leader       string `xml:"leader"` // cam a22        4500, cam ...
					Controlfield []struct {
						Text string `xml:",chardata"` // 1665767790, DE-627, 20190...
						Tag  string `xml:"tag,attr"`
					} `xml:"controlfield"`
					Datafield []struct {
						Text     string `xml:",chardata"`
						Tag      string `xml:"tag,attr"`
						Ind1     string `xml:"ind1,attr"`
						Ind2     string `xml:"ind2,attr"`
						Subfield []struct {
							Text string `xml:",chardata"` // 19,N20, dnb, 1185902589, ...
							Code string `xml:"code,attr"`
						} `xml:"subfield"`
					} `xml:"datafield"`
				} `xml:"record"`
			} `xml:"recordData"`
			RecordPosition string `xml:"recordPosition"` // 1, 2, 3, 4, 5, 6, 7, 8, 9...
		} `xml:"record"`
	} `xml:"records"`
	EchoedSearchRetrieveRequest struct {
		Text           string `xml:",chardata"`
		Version        string `xml:"version"`        // 1.2
		Query          string `xml:"query"`          // pica.ssg=24,1 or pica.ssg...
		StartRecord    string `xml:"startRecord"`    // 1
		MaximumRecords string `xml:"maximumRecords"` // 10
		RecordPacking  string `xml:"recordPacking"`  // xml
	} `xml:"echoedSearchRetrieveRequest"`
}

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Printf("%s %s\n", Version, BuildTime)
		os.Exit(0)
	}

	var vs = url.Values{}
	vs.Set("version", "1.2")
	vs.Set("operation", "searchRetrieve")
	vs.Set("query", *query)
	vs.Set("maximumRecords", strconv.Itoa(*maximumRecords))

	if *recordSchema != "" {
		vs.Set("recordSchema", *recordSchema)
	}

	var retrieved int

	re := regexp.MustCompile(`(?ms)(<record(.*?)</record>)`)

	if *recordRegex {
		fmt.Println("<collection>")
	}

	client := pester.New()
	client.Backoff = pester.ExponentialBackoff
	client.MaxRetries = 7

	for {
		vs.Set("startRecord", strconv.Itoa(*startRecord))

		link := fmt.Sprintf("%s?%s", *endpoint, vs.Encode())
		if *verbose {
			log.Println(link)
		}

		req, err := http.NewRequest("GET", link, nil)
		if err != nil {
			log.Fatal(err)
		}
		if *userAgent != "" {
			req.Header.Add("User-Agent", *userAgent)
		}

		// req.Header.Add("Accept-Encoding", "identity"), https://stackoverflow.com/q/21147562/89391
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			if *ignoreHTTPErrors {
				log.Printf("ignoring per flag: %s", resp.Status)
			} else {
				log.Fatal(resp.Status)
			}
		}
		var buf bytes.Buffer
		tee := io.TeeReader(resp.Body, &buf)

		dec := xml.NewDecoder(tee)
		var srr SearchRetrieveResponse
		if err := dec.Decode(&srr); err != nil {
			log.Fatal(err)
		}
		retrieved = retrieved + len(srr.Records.Record)
		if *limit > -1 && retrieved >= *limit {
			break
		}

		// Try to dig out: <record ... </record>
		if *recordRegex {
			for _, match := range re.FindAllString(buf.String(), -1) {
				match := strings.TrimSpace(match)
				fmt.Println(match)
			}
		} else {
			fmt.Println(buf.String())
		}
		buf.Reset()

		*startRecord = *startRecord + *maximumRecords
		if *startRecord > srr.NumberOfRecords {
			break
		}
	}

	if *recordRegex {
		fmt.Println("</collection>")
	}
}
