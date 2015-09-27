package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

var playgroundUrl = regexp.MustCompile(`https?\:\/\/play\.golang\.org\/p\/([a-zA-Z0-9_-]+)`)

type compileResult struct {
	Errors string `json:"compile_errors"`
	Output string `json:"output"`
}

func shrink(src string, size int) string {
	runes := []rune(src)
	if len(runes) > size {
		return string(runes[:size]) + fmt.Sprintf("\nstring is too long (%d chars has been removed)", len(runes)-size)
	}
	return src
}

func playground(addr string) (string, error) {
	res := ""
	resp, err := http.Get(addr + ".go")
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}
	res = string(bytes)
	resp, err = http.PostForm("https://play.golang.org/compile", map[string][]string{
		"body": []string{res},
	})
	res = shrink(res, 1500)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()
	compileRes := &compileResult{}
	err = json.NewDecoder(resp.Body).Decode(compileRes)
	if err != nil {
		return res, err
	}
	res = fmt.Sprintf("```%s```", res)
	if compileRes.Errors != "" {
		res = fmt.Sprintf("\n Errors:\n```%s```", res, shrink(compileRes.Errors, 1000))
	}
	if compileRes.Output != "" {
		res = fmt.Sprintf("%s\nResult:\n```%s```", res, shrink(compileRes.Output, 1000))
	}
	return res, nil
}
