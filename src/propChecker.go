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
	"strings"
)

const SRVRS_FILE string = "/home/aorfanos/projects/gommander/src/propchecker/server.list"

type DnsLgApiIntResponse struct {
	Name	string	`json:"name"`
	Type	string	`json:"type"`
	Class	string	`json:"class"`
	Ttl	int64	`json:"ttl"`
	Rdlength	int64	`json:"rdlength"`
	Rdata	string	`json:"rdata"`
}

type DnsLgApiExtResponse struct {
	Question	string	`json:"question"`
	Answer		string	`json:"answer"`
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

	fmt.Printf("Server\tIPs\t\n")

	for server := range serverPool {
		go func(msg string) {
			fmt.Printf("%s\t%s\n", serverPool[server], msg)
		}(parseResponse(makeRequest(*domainF, serverPool[server], *recordF)))
	}
}

func altParse(request string) (string) {
	d := DnsLgApiExtResponse{}
	e := json.NewDecoder(strings.NewReader(request)).Decode(&d)
	errChck(e)

	return d.Answer.Name
}

func parseResponse(response string) (string){
	data := DnsLgApiExtResponse{}
	errChck(json.Unmarshal([]byte(response), &data))

	return data
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
