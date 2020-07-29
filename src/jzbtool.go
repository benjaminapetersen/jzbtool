package main

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"	
	"flag"
	// third party
	// pretty print, something like jq
	// "github.com/TylerBrock/colorjson"
	"github.com/tidwall/pretty"
)

var (
	jzbInput string 
	jsonInput string 
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

func init() {
	const (
		defaultJZB = ""
		defaltJSON = ""		
	)
	flag.StringVar(&jzbInput, "jzb", "", "pass in jzb as a string to decode it to JSON")
	flag.StringVar(&jzbInput, "jbz", "", "oops! its jzb, but we got your back. pass in jbz, we will assume you meant jzb and will also decode it to JSON")
	flag.StringVar(&jsonInput, "json", "", "pass in JSON as a string to encode it to jzb")	
	// note that boolean flags must be provided as -pretty=false or --pretty=false, not --pretty false as 
	// the existence of the flag alone indicates truthiness.  This is an exception for booleans, other 
	// flags do not require the equals
	flag.BoolVar(&jsonPretty, "pretty", false, "if jzb provided, pretty print the JSON output")	
	flag.BoolVar(&jsonColor, "color", false, "if jzb provided, color the JSON output")	
}

func main() {
	// fail to parse and things go bad
	flag.Parse()
	// leave a blank line before output
	fmt.Println("");

	if(len(jzbInput) != 0) {
		// original
		jsonToJzb(jzbInput)
		return 
		// if(jsonPretty && jsonColor) {
		// 	jzbToJsonColorPretty(jzbInput)
		// 	return
		// } 
		// if(jsonPretty) {
		// 	jzbToJsonPretty(jzbInput)		
		// 	return
		// }
		// jzbToJsonSimple(jzbInput)
		// // jzbToJson(jzbInput)
		// return
	}
	// turning JSON into jzb doesn't have all the fancy bells :) 
	if(len(jsonInput) != 0) {
		jsonToJzbSimple(jsonInput)
		// jsonToJzb(jsonInput)
		return		
	}
	fmt.Fprintln(os.Stderr, "no input provided. please provide from the following:")
	flag.PrintDefaults()	
	os.Exit(1)
}
