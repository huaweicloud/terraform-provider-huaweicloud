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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud"
)

var (
	commentFormat = "// @API Product httpMethod requestPath"
	httpMethods   = map[string]struct{}{
		"GET":     {},
		"POST":    {},
		"PUT":     {},
		"DELETE":  {},
		"PATCH":   {},
		"HEAD":    {},
		"OPTIONS": {},
		"CONNECT": {},
		"TRACE":   {},
	}

	ignoreFile = map[string]struct{}{
		"resource_schema":                       {},
		"resource_huaweicloud_vpc_bandwidth_v1": {},
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
	provider := huaweicloud.Provider()

	errCode := checkApiComments(basePath, provider)
	os.Exit(errCode)
}

func checkApiComments(packagePath string, provider *schema.Provider) int {
	errCode := 0
	err := filepath.Walk(packagePath, func(path string, fInfo os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("scan path %s failed: %s\n", path, err)
			return err
		}

		if fInfo.IsDir() && !isSkipDirectory(path) {
			errCode += dealFileApiComments(path, provider)
		}

		return nil
	})
	if err != nil {
		fmt.Printf("ERROR: scan path failed: %s\n", err)
		return 1
	}
	return errCode
}

func dealFileApiComments(path string, provider *schema.Provider) int {
	fSet := token.NewFileSet()
	packs, err := parser.ParseDir(fSet, path, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("[ERROR] Failed to parse package %s: %s\n", path, err)
		return 1
	}

	var errorCount int
	for _, pack := range packs {
		for filePath, file := range pack.Files {
			// 获得文件名并去除版本号
			fileName := filePath[strings.LastIndex(filePath, "/")+1 : len(filePath)-3]
			if _, ok := ignoreFile[fileName]; ok {
				continue
			}

			reg := regexp.MustCompile(`_v\d+$`)
			fileNameWithoutVersion := reg.ReplaceAllString(fileName, "")
			if !isExportResource(fileNameWithoutVersion, provider) {
				continue
			}

			hasApiComments := false
			for _, group := range file.Comments {
				for _, comment := range group.List {
					hasApiComment, count := apiCommentCheck(comment.Text)
					errorCount += count
					if hasApiComment {
						hasApiComments = hasApiComment
					}
				}
			}
			if !hasApiComments {
				errorCount++
				fmt.Printf("[ERROR] the file `%s` has not API comments\n", fileName)
			}
		}
	}
	return errorCount
}

func isExportResource(resourceName string, p *schema.Provider) bool {
	if strings.HasSuffix(resourceName, "_test") {
		return false
	}
	if strings.HasPrefix(resourceName, "resource_") {
		simpleFilename := strings.TrimPrefix(resourceName, "resource_")
		// check whether the resource is internal
		if v, ok := p.ResourcesMap[simpleFilename]; ok {
			return v.Description != "schema: Internal"
		}
	}

	if strings.HasPrefix(resourceName, "data_source_") {
		simpleFilename := strings.TrimPrefix(resourceName, "data_source_")
		// check whether the data source is internal
		if v, ok := p.DataSourcesMap[simpleFilename]; ok {
			return v.Description != "schema: Internal"
		}
	}
	return false
}

func apiCommentCheck(commentStr string) (bool, int) {
	commentStr = strings.TrimSpace(commentStr)
	apiReg := regexp.MustCompile(`(// @API)\s*(.*)`)
	apiMatch := apiReg.FindAllStringSubmatch(commentStr, -1)
	if len(apiMatch) == 0 {
		return false, 0
	}
	index := strings.Index(commentStr, "// @API")
	if commentStr[index+len("// @API"):index+len("// @API")+1] != " " {
		fmt.Printf("[ERROR] the format of `%s` is not correct, should be (%s)\n", commentStr, commentFormat)
		return true, 1
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
		return true, 1
	}

	requestMethod := parts[1]
	if _, ok := httpMethods[requestMethod]; !ok {
		fmt.Printf("[ERROR] the http method of `%s` is not correct\n", commentStr)
		return true, 1
	}

	url := parts[2]
	if !strings.HasPrefix(url, "/") {
		fmt.Printf("[ERROR] the request url of `%s` is not correct, the URL path should start with /\n", commentStr)
		return true, 1
	}
	if strings.HasSuffix(url, "/") {
		fmt.Printf("[ERROR] the request url of `%s` is not correct, the URL path should not end with /\n", commentStr)
		return true, 1
	}
	return true, 0
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
