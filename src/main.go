package main

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"flag"
	// third party
	// pretty print, something like jq
	// "github.com/TylerBrock/colorjson"
	"github.com/tidwall/pretty"
)

var (
	// jzbInput string
	// jsonInput string
	jsonPretty bool
	jsonColor bool
)


type JzbCorruptError struct {
	err error
}

func (e JzbCorruptError) Error() string {
	return fmt.Sprintf("%s", e.err)
}

func RawMessageToJzb(j json.RawMessage) (string, error) {
	buf := &bytes.Buffer{}
	b64 := base64.NewEncoder(base64.URLEncoding, buf)
	compress := zlib.NewWriter(b64)

	if i, err := compress.Write(j); err != nil {
		return "", err
	} else if i != len(j) {
		return "", fmt.Errorf("ToJzb() writer did not accept all input bytes")
	} else if err := compress.Close(); err != nil {
		return "", err
	} else if err := b64.Close(); err != nil {
		return "", err
	} else {
		s := string(buf.Bytes())
		for s != "" && s[len(s)-1] == '=' {
			s = s[:len(s)-1]
		}

		return s, nil
	}
}

func ToJzb(thing interface{}) (string, error) {

	if j, err := json.Marshal(thing); err != nil {
		return "", err
	} else {
		return RawMessageToJzb(j)
	}
}

func JzbToRawMessageStream(jzb string) (io.Reader, error) {
	switch len(jzb) % 4 {
	case 1:
		return nil, fmt.Errorf("bad base64 encoding")
	case 2:
		jzb += "=="
	case 3:
		jzb += "="
	}

	compressed := base64.NewDecoder(base64.URLEncoding, strings.NewReader(jzb))
	uncompressed, err := zlib.NewReader(compressed)

	if err != nil {
		return nil, fmt.Errorf("Error with compressed input: %s", err.Error())
	}

	return uncompressed, nil
}

// This converts all numbers to float64, which means 1492517388551 remarshals as 1.492517388551e+12
func FromJzb(jzb string, thing interface{}) error {
	return FromJzbUseNumber(jzb, thing, false)
}

func FromJzbUseNumber(jzb string, thing interface{}, useNumber bool) error {
	uncompressed, err := JzbToRawMessageStream(jzb)
	if err != nil {
		return err
	}

	buffered := bufio.NewReader(uncompressed)
	var reader io.Reader
	if b, err := buffered.Peek(1); err != nil {
		if err != nil {
			return JzbCorruptError{fmt.Errorf("Error peeking into compressed input: %s", err.Error())}
		}
	} else if b[0] == '"' {
		// this was double stringified. gross.
		var s string
		decoder := json.NewDecoder(buffered)
		if useNumber {
			decoder.UseNumber()
		}
		if err := decoder.Decode(&s); err != nil {
			return fmt.Errorf("Cannot parse stringified json: %s", err.Error())
		}

		reader = strings.NewReader(s)
	} else {
		reader = buffered
	}

	decoder := json.NewDecoder(reader)
	if useNumber {
		decoder.UseNumber()
	}
	if err := decoder.Decode(thing); err != nil {
		return fmt.Errorf("Cannot parse json: %s", err.Error())
	}

	return nil
}

func convertToUrlSafe(in string) string {
	out := strings.Replace(in, "+", "-", -1)
	return strings.Replace(out, "/", "_", -1)
}

// original
func jsonToJzb(j string) {
	if jzb, err := ToJzb(string(j)); err != nil {
		fmt.Printf("jsonToJzb err = %v\n", err)
	} else {
		fmt.Printf("\nJSON: %s\n", j)
		fmt.Printf("JZB : %s\n", jzb)
	}
}

// strip off the extra output so that it can be piped to jq, etc
func jsonToJzbSimple(j string) {
	if jzb, err := ToJzb(string(j)); err != nil {
		fmt.Printf("jsonToJzb err = %v\n", err)
	} else {
		// fmt.Printf("\nJSON: %s\n", j)
		fmt.Printf("%s\n", jzb)
	}
}

// original
func jzbToJson(jzb string) {
	jzb = convertToUrlSafe(jzb)
	var orig interface{}
	if err := FromJzb(jzb, &orig); err != nil {
		fmt.Printf("jzbToJson err = %v\n", err)
	} else {
		j, _ := json.Marshal(orig)
		fmt.Printf("\nJZB : %s\n", jzb)
		fmt.Printf("JSON: %s\n", j)
	}
}

// strip off the extra output so that it can be piped to jq, etc
func jzbToJsonSimple(jzb string) {
	jzb = convertToUrlSafe(jzb)
	var orig interface{}
	if err := FromJzb(jzb, &orig); err != nil {
		fmt.Printf("jzbToJson (simple) err = %v\n", err)
	} else {
		j, _ := json.Marshal(orig)
		// fmt.Printf("\nJZB : %s\n", jzb)
		fmt.Printf("%s\n", j)
	}
}

// attempt to remove the need for jq by importing a pretty print library
func jzbToJsonPretty(jzb string) {
	jzb = convertToUrlSafe(jzb)
	var orig interface{}
	if err := FromJzb(jzb, &orig); err != nil {
		fmt.Printf("jzbToJson (pretty) err = %v\n", err)
	} else {
		j, _ := json.Marshal(orig)
		result := pretty.Pretty(j)
		// fmt.Printf("\nJZB : %s\n", jzb)
		fmt.Printf("%s\n", result)
	}
}
// this is a little copy-pasta
func jzbToJsonColorPretty(jzb string) {
	jzb = convertToUrlSafe(jzb)
	var orig interface{}
	if err := FromJzb(jzb, &orig); err != nil {
		fmt.Printf("jzbToJson (color,pretty) err = %v\n", err)
	} else {
		j, _ := json.Marshal(orig)
		result := pretty.Pretty(j)
		colored := pretty.Color(result, nil)
		// fmt.Printf("\nJZB : %s\n", jzb)
		fmt.Printf("%s\n", colored)
	}
}

func isJSONString(s string) bool {
	var js string
	return json.Unmarshal([]byte(s), &js) == nil

}

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil

}

func isValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

func pluckJZBFromQuery(toTest string) string {
	val, _ := url.ParseQuery(toTest)
	// ignore jzb/jbz human error, the URL should be built by machines
	return val["jzb"][0]
}


func processAndPrintJZB(toProcess string) {
	// leave a blank line before output, just for readability
	fmt.Println("")
	if len(toProcess) != 0 {
		if jsonPretty && jsonColor {
			jzbToJsonColorPretty(toProcess)
			return
		}
		if jsonPretty {
			jzbToJsonPretty(toProcess)
			return
		}
		jzbToJsonSimple(toProcess)
	}
}


func init() {
	//const (
	//	defaultJZB = ""
	//	defaultJSON = ""
	//)
	// no reason to need these flags anymore, they are the default behavior
	// flag.StringVar(&jzbInput, "jzb", "", "pass in jzb as a string to decode it to JSON")
	// flag.StringVar(&jzbInput, "jbz", "", "oops! its jzb, but we got your back. pass in jbz, we will assume you meant jzb and will also decode it to JSON")
	// flag.StringVar(&jsonInput, "json", "", "pass in JSON as a string to encode it to jzb")

	// note that boolean flags must be provided as -pretty=false or --pretty=false, not --pretty false as
	// the existence of the flag alone indicates truthiness.  This is an exception for booleans, other
	// flags do not require the equals
	// assuming a jq style output first, but allowing the ability to turn it off
	flag.BoolVar(&jsonPretty, "pretty", true, "optionally turn off pretty printing")
	flag.BoolVar(&jsonColor, "color", true, "optionally turn off color when pretty printing")
}

func main() {
	flag.Parse()

	// we should get one single arg. flag.Args() is better than os.Args in that it
	// cleans out the flags and returns what wemains (and drops the program name itself)
	arguments := flag.Args()

	// first arg is the program itself
	// second should be our string
	if len(arguments) == 0 {
		// leave a blank line before output, just for readability
		fmt.Println("")
		fmt.Fprintln(os.Stderr, "no input provided. please provide an input, for example:")
		fmt.Println("  jzbtool eJxSUgIEAAD__wBoAEU")
		fmt.Println("  jzbtool '{\"name\": \"Jane\"}'")
		flag.PrintDefaults()
		os.Exit(1)
	}

	argument := arguments[0]

	// 1. if its json, encode the jzb and print it
	if isJSON(argument) {
		// fmt.Println("Processing as JSON")
		jsonToJzbSimple(argument)
		return
	}

	// 2. if its a URL, pluck the jzb, then do the jzb decoding
	if isValidUrl(argument) {
		// fmt.Println("Processing as URL")
		jzbStr := pluckJZBFromQuery(argument)
		processAndPrintJZB(jzbStr)
		return
	}

	// 3. otherwise, treat whatever string we get as a JZB and attempt to decode it.
	// fmt.Println("Processing string as a JZB")
	processAndPrintJZB(argument)
}



// to manually test:
// go run main.go 'https://www.helloworld.com?foo=bar&jzb=eJxSqo5RykvMTY1RslKIUfJKzEuNUapVAgQAAP__TdEGrg' 	// pretty prints {"name": "Jane"}
// go run main.go eJxSqo5RykvMTY1RslKIUfJKzEuNUapVAgQAAP__TdEGrg 											// pretty prints {"name": "Jane"}
// go run main/go '{"name": "Jane"}' 																		// prints eJxSqo5RykvMTY1RslKIUfJKzEuNUapVAgQAAP__TdEGrg
// go run main.go eJxSqo5RykvMTY1RslKIUfJKzEuNUapVAgQAAP__TdEGrg 											// pretty prints {"name": "Jane"}
// go run main.go eJxSqo5RykvMTY1RslKIUfJKzEuNUapVAgQAAP__TdEGrg --color=false
// go run main.go --pretty=false --color=false eJxSqo5RykvMTY1RslKIUfJKzEuNUapVAgQAAP__TdEGrg
// go run main.go -pretty=false -color=false eJxSqo5RykvMTY1RslKIUfJKzEuNUapVAgQAAP__TdEGrg
// go run main.go -color=false eJxSqo5RykvMTY1RslKIUfJKzEuNUapVAgQAAP__TdEGrg
// will not work:
// must have = for bool flags:
//   go run main.go -pretty false -color false eJxSqo5RykvMTY1RslKIUfJKzEuNUapVAgQAAP__TdEGrg
// flags must come before args:
//   go run main.go eJxSqo5RykvMTY1RslKIUfJKzEuNUapVAgQAAP__TdEGrg -pretty=false -color=false
