package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Global values
var (
	ToScan, bl      int
	total           int
	Results         []string
	DorkFile        string
	DenyFile        string
	Concurrency     bool
	JSON            bool
	dir             = "."
	TotalOutputFile = "All.txt"
)

type gen struct {
	PageTitle string
	Links     []string
}

var blck = []string{"go.microsoft.com"}

// init parses the command line flags
// and default values
func init() {
	flag.StringVar(&DorkFile, "dork", "dorks.txt", "Path of your dork file!")
	flag.StringVar(&DenyFile, "deny", "deny.txt", "Path of your deny file!")
	flag.BoolVar(&Concurrency, "turbo", true, "Concurrectly searches dorks. Turn it off (-turbo=false) if you have a lot of dorks and time.")
	flag.BoolVar(&JSON, "json", false, "json output in every api endpoint")
	flag.Parse()
}

// load loads given file values by flag
// or loads the default ones!
func load() {
	dorkfileBytes, err := ioutil.ReadFile(DorkFile)
	if err != nil {
		fmt.Println(DorkFile, ":: file not found\nError:", err)
		os.Exit(1)
		return
	}
	Dorks := strings.Split(string(dorkfileBytes), "\n")
	denyfileBytes, err := ioutil.ReadFile(DenyFile)
	if err != nil {
		fmt.Println(DenyFile, ":: file not found\nError:", err)
		fmt.Println("But still I am going to scan all dorks!")
	} else {
		sp := strings.Split(string(denyfileBytes), "\n")
		sp = append(sp, "go.microsoft.com")
		blck = sp
	}
	ToScan = len(Dorks)
	if Concurrency {
		for i := range Dorks {
			go Botify(Dorks[i])
		}
		return
	}
	for i := range Dorks {
		fmt.Println("Working on:", Dorks[i])
		Botify(Dorks[i])
	}
}

// main function to load
func main() {
	println()
	println()
	println()
	fmt.Println("BingBot!\nCreated By Anik Hasibul (@AnikHasibul)!")
	fmt.Println("•••••••••••••••••")
	fmt.Println("Results on:")
	fmt.Println("http://localhost:1338")
	fmt.Println("•••••••••••••••••")
	fmt.Println("Searching....\nFor exit:\npress CTRL+C\nor visit\nhttp://127.0.0.1:1338/exit")
	println()
	println()
	println()
	load()
	// Server and handler functions
	http.HandleFunc("/", live)
	http.HandleFunc("/exit", Exit)
	http.HandleFunc("/reload", Reload)
	http.HandleFunc("/domain", Host)
	erb := http.ListenAndServe(":1338", nil)
	if erb != nil {
		fmt.Println("Got an error on http serving!\n", erb, "\nBut I am working for 10 second...")
		time.Sleep(10000 * time.Millisecond)
	}
	// count the collected sites
	defer Count()
}

// Botify searches the dorks on bing
func Botify(searchStr string) {
	if searchStr == "" {
		return
	}
	start := time.Now()
	// pre: previous page's results
	// nex: next(current) page's results
	// sey: final results to add to collection
	var pre, nex, sey []string
	// report the time took to collect sites
	defer func() {
		since := String(time.Since(start))
		fmt.Println(since, ":: DONE ::", searchStr)
	}()

	/* Searching starts here */
	var page int
	nex = []string{""}
	// Search until we reach the last page
	for {
		//Bing page indicator
		page += 10
		/* HTTP client stuffs */
		client := http.Client{}
		req, ReqErr := http.NewRequest(http.MethodGet, fmt.Sprintf("https://www.bing.com/search?q=%s&first=%d", url.QueryEscape(searchStr), page), nil)
		if ReqErr != nil {
			fmt.Println("Error: ", ReqErr)
		}
		// Spoof a old browser user agent
		// for getting a simple html page
		req.Header.Set("User-Agent",
			"Nokia2700c/10.0.011 (SymbianOS/9.4; U; Series60/5.0 Opera/5.0; Profile/MIDP-2.1 Configuration/CLDC-1.1 ) AppleWebKit/525 (KHTML, like Gecko) Safari/525 3gpp-gba",
		)

		resp, RespErr := client.Do(req)
		if RespErr != nil {
			fmt.Println("Error: ", RespErr)
		}
		body, BodyErr := ioutil.ReadAll(resp.Body)
		// we don't need resp.Body anymore
		resp.Body.Close()
		if BodyErr != nil {
			fmt.Println("Error: ", BodyErr)
		}
		/* HTTP client ends here! we got the response in body */
		// split when a search result link captured
		htmlExtractionSlice := strings.Split(string(body), `<a _ctf="rdr_T" href="http`)
		// We didn't get a valid page of bing results and it contains less than 10 results or no results!
		if len(htmlExtractionSlice) < 10 || eq(pre, nex) {
			// Let's save the sites because we are not going to scan this dork
			err := ioutil.WriteFile(
				url.PathEscape(searchStr)+".bingbot.txt",
				[]byte(strings.Join(sey[len(pre):], "\n")),
				0666)
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		pre = nex
		nex = []string{}
		// lets extract the links from <a> tag
		for i := range htmlExtractionSlice {
			if i == 0 {
				// ignore the first link (it is a internal bing link)
				continue
			}
			ExtractedHref := strings.Split(htmlExtractionSlice[i], `"`)
			if IsBlocked(ExtractedHref[0]) {
				// ignore denied links or sites
				continue
			}
			nex = append(nex, ExtractedHref[0])
			// append them with previous results
			Results = append(Results, "http"+ExtractedHref[0])
			sey = append(sey, "http"+ExtractedHref[0])
			total++
		}
	}
}

/* HTTP handler for  / or the root handler */
func live(w http.ResponseWriter, _ *http.Request) {
	if JSON {
		OutputWithJSON(w, Results)
		return
	}
	ser := gen{
		PageTitle: fmt.Sprintf("BingBot :: By Anik Hasibul [ %d Site Collected]", total),
		Links:     Results,
	}

	t, ert := template.New("lol").Parse(`
	
	<h1>{{.PageTitle}}<h1>
<ul>
    {{range .Links}}
           <li> <a href="{{.}}">{{.}}</a></li><br>
    {{end}}
</ul>`)
	if ert != nil {
		fmt.Println("Parse error:", ert)
	}
	err := t.ExecuteTemplate(w, "lol", ser)
	if err != nil {
		fmt.Println("Execution error:", err)
	}
}

// Reload reloads the Dorks and Deny loader
// without exiting the bot
/* HTTP handler for reload */
func Reload(w http.ResponseWriter, _ *http.Request) {

	fmt.Fprintln(w, "Realoaded! Check your terminal output!")
	fmt.Println("====================")
	load()
	fmt.Println("====================")
}

// Host saves only the domain name of collected site
/* HTTP handler for /domain */
func Host(w http.ResponseWriter, _ *http.Request) {
	var s []string
	for i := range Results {
		v, _ := url.ParseRequestURI(Results[i])
		s = append(s, "http://"+v.Host)
		// list only domains (not the full url)
		//[BUG only http:// listing available, https:// are not available ]
	}
	// decline the duplicate hosts
	s = Unique(s)
	// Save it
	err := ioutil.WriteFile(
		filepath.Join(dir, "hosts.txt"),
		[]byte(strings.Join(s, "\n")),
		0666,
	)
	if err != nil {
		log.Println(err)
	}
	if JSON {
		OutputWithJSON(w, s)
	}
	fmt.Fprintf(
		w,
		"Done![%s]",
		filepath.Join(dir, "hosts.txt"),
	)
}

// Exit exits the bot via web
/* HTTP handler for /exit */
func Exit(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "Exited!")
	time.Sleep(1 * time.Second)
	os.Exit(0)
}

/* Synchronizing the output with zzz.txt */
func eq(a, b []string) bool {
	// check the output is not the previous pages results!
	// if it contains new results, then list them else return false finish the scan
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	// save it if the results are new
	go Save()
	return true
}

// Save saves the results on each and every page scan
func Save() {
	err := ioutil.WriteFile(filepath.Join(dir, TotalOutputFile), []byte(strings.Join(Results, "\n")), 0666)
	if err != nil {
		log.Println(err)
	}
}

// IsBlocked returns the true if the site/link matches any of Deny params
func IsBlocked(str string) bool {
	for i := range blck {
		if blck[i] != "" {
			if strings.Contains(str, blck[i]) {
				bl++
				return true
			}
		}
	}
	return false
}

// Unique removes duplicate values from domain array
func Unique(sites []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range sites {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// Count counts the total collected sites
func Count() {
	fmt.Println("\n\nTotal", total, "url has been saved!")
}

// String formats interface{}'s values to string
func String(s ...interface{}) string {
	return fmt.Sprintf("%v", s...)
}

// OutputWithJSON outputs json encoded data
func OutputWithJSON(w http.ResponseWriter, data interface{}) {

	output, err := json.MarshalIndent(Results, "", "    ")
	if err != nil {
		log.Println(err)
		OutputWithJSON(w, err.Error())
		return
	}
	w.Header().Set(
		"Content-Type",
		"application/json",
	)
	fmt.Fprintln(w, string(output))
}
