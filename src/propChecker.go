package main

import (
	"fmt"
	"encoding/json"
	"flag"
	"net/http"
	"bufio"
	"log"
	"os"
	"io/ioutil"
)

const SRVRS_FILE string = "/home/aorfanos/projects/gommander/src/propchecker/server.list"

type DnsLgApiResponse struct {
	Name string `json:"name"`
	Type string `json: "type"`
	Class string `json: "class"`
	Ttl int64 `json: "ttl"`
	Rdlength int64 `json: "rdlength"`
	Rdata string `json: "rdata"`
}

func main() {

	//define flags - ending of F denotes a flag, accessed by pointer
	//flag.String(flagLabel, defaultVal, description)
	domainF := flag.String("domain", "google.com", "The domain to be examined")
	recordF := flag.String("type", "a", "Type of DNS record to return")
	poolF := flag.String("pool", SRVRS_FILE, "Pool to read from (values can be found http://www.dns-lg.com/")
	flag.Parse()

	fmt.Println("Propagation Checker init...")
	fmt.Printf("Checking for %s records of %s\n", *recordF, *domainF)

	serverPool, e := readLines(*poolF)
	errChck(e)

	fmt.Printf("Server\tStatus\t\n")

	for server := range serverPool {
		go func(msg string) {
			fmt.Printf("%s\t%s\n", serverPool[server], msg)
		}(makeRequest(*domainF, serverPool[server], *recordF))
	}
}

func readLines(path string) ([]string, error) {
	file, e := os.Open(path)
	errChck(e)
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()

}

func parseResponse(body []byte) (*DnsLgApiResponse, error) {
	var s = new(DnsLgApiResponse)
	err := json.Unmarshal(body, &s)
	if ( err != nil) {
		fmt.Println(err)
	}
	return s, err
}

func makeRequest(domain, server, reqType string) (string) {
	uri := fmt.Sprintf("http://www.dns-lg.com/%s/%s/%s", server, domain, reqType )
	resp, e := http.Get(uri)
	errChck(e)
	defer resp.Body.Close()

	data, e := ioutil.ReadAll(resp.Body)
	errChck(e)

	//return data

	return string(data)
}

func errChck(e error){
	if e != nil {
		if os.IsPermission(e) {
			log.Fatal(e)
		}
		log.Fatal(e)
	}
}
