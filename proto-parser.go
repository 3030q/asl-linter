package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pterm/pterm"
	"os"

	"path/filepath"

	protoparser "github.com/yoheimuta/go-protoparser/v4"
)

var (
	debug      = flag.Bool("debug", false, "debug flag to output more parsing process detail")
	permissive = flag.Bool("permissive", true, "permissive flag to allow the permissive parsing rather than the just documented spec")
	unordered  = flag.Bool("unordered", false, "unordered flag to output another one without interface{}")
)

func protoParse(fileName string) map[string]interface{} {
	flag.Parse()
	protoFileName := fileName + ".proto"
	reader, err := os.Open(protoFileName)
	if err != nil {
		pterm.Error.Println(err)
		return nil
	} else {
		pterm.Success.Println("Successfully Opened", protoFileName)
	}
	defer reader.Close()

	got, err := protoparser.Parse(
		reader,
		protoparser.WithDebug(*debug),
		protoparser.WithPermissive(*permissive),
		protoparser.WithFilename(filepath.Base(protoFileName)),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse, err %v\n", err)
		return nil
	}

	var v interface{}
	v = got
	if *unordered {
		v, err = protoparser.UnorderedInterpret(got)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to interpret, err %v\n", err)
			return nil
		}
	}

	gotJSON, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal, err %v\n", err)
	}
	var result map[string]interface{}
	json.Unmarshal([]byte(gotJSON), &result)
	return result
}
