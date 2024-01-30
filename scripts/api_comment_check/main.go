package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	commentFormat = "// @API Product httpMethod requestPath"
	httpMethods   = map[string]bool{
		"GET":     true,
		"POST":    true,
		"PUT":     true,
		"DELETE":  true,
		"PATCH":   true,
		"HEAD":    true,
		"OPTIONS": true,
		"CONNECT": true,
		"TRACE":   true,
	}
)

var (
	// 命令行参数
	basePath string
)

func init() {
	flag.StringVar(&basePath, "basePath", "../../huaweicloud/", "base Path")
}

func main() {
	flag.Parse()

	errCode := checkApiComments(basePath)
	os.Exit(errCode)
}

func checkApiComments(packagePath string) int {
	errCode := 0
	err := filepath.Walk(packagePath, func(path string, fInfo os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("scan path %s failed: %s\n", path, err)
			return err
		}

		if fInfo.IsDir() && !isSkipDirectory(path) {
			errCode += dealFileApiComments(path)
		}

		return nil
	})
	if err != nil {
		fmt.Printf("ERROR: scan path failed: %s\n", err)
		return 1
	}
	return errCode
}

func dealFileApiComments(path string) int {
	fSet := token.NewFileSet()
	packs, err := parser.ParseDir(fSet, path, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("[ERROR] Failed to parse package %s: %s\n", path, err)
		return 1
	}

	var errorCount int
	for _, pack := range packs {
		for _, file := range pack.Files {
			for _, group := range file.Comments {
				for _, comment := range group.List {
					errorCount += apiCommentCheck(comment.Text)
				}
			}
		}
	}
	return errorCount
}

func apiCommentCheck(commentStr string) int {
	commentStr = strings.TrimSpace(commentStr)
	apiReg := regexp.MustCompile(`(// @API)\s*(.*)`)
	apiMatch := apiReg.FindAllStringSubmatch(commentStr, -1)
	if len(apiMatch) == 0 {
		return 0
	}
	index := strings.Index(commentStr, "// @API")
	if commentStr[index+len("// @API"):index+len("// @API")+1] != " " {
		fmt.Printf("[ERROR] the format of `%s` is not correct, should be (%s)\n", commentStr, commentFormat)
		return 1
	}

	str := strings.TrimSpace(apiMatch[0][2])
	standardStr := ""
	// remove the extra spaces
	for i, s := range str {
		if string(s) != " " {
			// only one space is keep
			if i > 0 && string(str[i-1]) == " " {
				standardStr += " "
			}
			standardStr += string(s)
		}
	}
	parts := strings.Split(standardStr, " ")
	if len(parts) != 3 {
		fmt.Printf("[ERROR] the format of `%s` is not correct, should be (%s)\n", commentStr, commentFormat)
		return 1
	}

	requestMethod := parts[1]
	if !httpMethods[requestMethod] {
		fmt.Printf("[ERROR] the http method of `%s` is not correct\n", commentStr)
		return 1
	}

	url := parts[2]
	if !strings.HasPrefix(url, "/") {
		fmt.Printf("[ERROR] the request url of `%s` is not correct, the URL path should start with /\n", commentStr)
		return 1
	}
	if strings.HasSuffix(url, "/") {
		fmt.Printf("[ERROR] the request url of `%s` is not correct, the URL path should not end with /\n", commentStr)
		return 1
	}
	return 0
}

func isSkipDirectory(path string) bool {
	var skipKeys = []string{
		"acceptance", "utils", "internal", "helper", "obs", "deprecated",
	}

	for _, sub := range skipKeys {
		if strings.Contains(path, sub) {
			return true
		}
	}
	return false
}
