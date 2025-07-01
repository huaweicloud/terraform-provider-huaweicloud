package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud"
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
	errCode := dealResources(basePath, provider)
	os.Exit(errCode)
}

func dealResources(packagePath string, provider *schema.Provider) int {
	missingCheckDeletedResourceCount := 0
	err := filepath.Walk(packagePath, func(path string, fInfo os.FileInfo, err error) error {
		if err != nil {
			log.Printf("scan path %s failed: %s\n", path, err)
			return err
		}

		if fInfo.IsDir() {
			missingCheckDeletedResourceCount += dealResource(path, provider)
		}

		return nil
	})
	if err != nil {
		fmt.Printf("ERROR: scan path failed: %s\n", err)
		return 1
	}
	fmt.Printf("missing check deleted resource count: %v\n", missingCheckDeletedResourceCount)
	return missingCheckDeletedResourceCount
}

func dealResource(path string, provider *schema.Provider) int {
	fSet := token.NewFileSet()
	packs, err := parser.ParseDir(fSet, path, nil, 0)
	if err != nil {
		fmt.Printf("Failed to parse package %s: %s\n", path, err)
		return 1
	}
	count := 0
	for _, pack := range packs {
		if isSkipDirectory(pack.Name) {
			continue
		}

		for filePath := range pack.Files {
			// 获得文件名并去除版本号
			fileName := filePath[strings.LastIndex(filePath, "/")+1 : len(filePath)-3]

			reg := regexp.MustCompile(`_v\d+$`)
			fileNameWithoutVersion := reg.ReplaceAllString(fileName, "")

			if !isExportResource(fileNameWithoutVersion, provider) || isSkipResource(fileNameWithoutVersion) {
				continue
			}

			resourceFileBytes, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("[ERROR] failed to read file %s: %s\n", filePath, err)
				return 1
			}

			fileStr := string(resourceFileBytes)

			if strings.HasPrefix(fileNameWithoutVersion, "resource_") {
				funcReg := regexp.MustCompile("(func )([r|R]esource.*Read)(\\(.* context.Context,[ |\n\t].* " +
					"\\*schema.ResourceData,[ |\n\t].* interface\\{}\\) diag.Diagnostics \\{)")
				allFuncMatch := funcReg.FindAllStringSubmatch(fileStr, -1)
				// 表示使用了其他资源的读方法
				if allFuncMatch == nil {
					// 表示老写法：func resourceClusterV1Read(d *schema.ResourceData, meta interface{}) error {
					funcReg = regexp.MustCompile("(func )([r|R]esource.*Read)(\\(.* \\*schema.ResourceData," +
						"[ |\n\t].* interface\\{}\\) error \\{)")
					allFuncMatch = funcReg.FindAllStringSubmatch(fileStr, -1)
					if allFuncMatch == nil {
						continue
					}
				}

				subStrs := strings.Split(fileStr, allFuncMatch[0][0])
				subStr := subStrs[1]

				endReg := regexp.MustCompile(`\n}\n`)
				allEndMatch := endReg.FindAllStringSubmatch(subStr, -1)

				subSubStrs := strings.Split(subStr, allEndMatch[0][0])
				subSubStr := subSubStrs[0]

				// 如果read方法中没有逻辑（动作类资源），那就不需要判断是否有checkDeleted
				emptyReg := regexp.MustCompile(`meta\.\(\*config\.Config\)`)
				allEmptyMatch := emptyReg.FindAllStringSubmatch(subSubStr, -1)
				if allEmptyMatch == nil {
					continue
				}

				checkDeletedReg := regexp.MustCompile(`return common.CheckDeletedDiag`)
				allCheckDeletedMatch := checkDeletedReg.FindAllStringSubmatch(subSubStr, -1)
				if len(allCheckDeletedMatch) == 0 {
					checkDeletedReg = regexp.MustCompile(`return CheckDeleted`)
					allCheckDeletedMatch = checkDeletedReg.FindAllStringSubmatch(subSubStr, -1)
					if len(allCheckDeletedMatch) == 0 {
						fmt.Printf("resource whose read method missing checkDeleted: %s\n", fileNameWithoutVersion)
						count++
					}
				}
			}
		}
	}
	return count
}

func isExportResource(resourceName string, provider *schema.Provider) bool {
	if strings.HasSuffix(resourceName, "_test") {
		return false
	}
	if strings.HasPrefix(resourceName, "resource_") {
		simpleFilename := strings.TrimPrefix(resourceName, "resource_")
		if _, ok := provider.ResourcesMap[simpleFilename]; ok {
			return true
		}
	}
	return false
}

func isSkipDirectory(name string) bool {
	var skipDirs = map[string]bool{
		"acceptance": true,
		"utils":      true,
		"internal":   true,
		"helper":     true,
		"deprecated": true,
	}

	return skipDirs[name]
}

func isSkipResource(name string) bool {
	var skipResources = map[string]bool{}

	return skipResources[name]
}
