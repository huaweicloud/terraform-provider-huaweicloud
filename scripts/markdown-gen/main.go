package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud"
)

var (
	sourceName string
	isData     bool

	commandLine flag.FlagSet
)

func init() {
	commandLine.Init(os.Args[0], flag.ExitOnError)

	commandLine.BoolVar(&isData, "d", false, "Indicates the input name is a data source")
	commandLine.StringVar(&sourceName, "name", "", "The resource or data source name.")

	commandLine.Usage = func() {
		fmt.Fprintf(commandLine.Output(), "Usage of %s:\n\n", os.Args[0])
		commandLine.PrintDefaults()
	}
}

func main() {
	commandLine.Parse(os.Args[1:]) //nolint: errcheck

	if sourceName == "" {
		fmt.Printf("-name must be specified\n")
		os.Exit(1)
	}

	provider := huaweicloud.Provider()
	schemaResource := getSchemaResource(provider, sourceName, isData)
	if schemaResource == nil {
		fmt.Printf("cannot find the resource or data source [%s]\n", sourceName)
		os.Exit(2)
	}

	b := &strings.Builder{}
	if err := Render(b, schemaResource, sourceName); err != nil {
		fmt.Printf("falied to render %s: %s", sourceName, err)
		fmt.Printf("\n%s\n", b.String())
		os.Exit(3)
	}

	outFile := genPathFile(sourceName)
	output, err := os.Create(outFile)
	if err != nil {
		fmt.Printf("falied to create output file: %s", err)
		os.Exit(4)
	}

	_, err = output.WriteString(b.String())
	if err != nil {
		fmt.Printf("falied to write file: %s", err)
		output.Close()
		os.Exit(5)
	}

	output.Close()
	os.Exit(0)
}

func getSchemaResource(provider *schema.Provider, name string, isData bool) *schema.Resource {
	var resourceMap map[string]*schema.Resource

	if !isData {
		resourceMap = provider.ResourcesMap
	} else {
		resourceMap = provider.DataSourcesMap
	}

	return resourceMap[name]
}

func genPathFile(name string) string {
	path := name
	ret := strings.SplitN(name, "_", 2)
	if len(ret) == 2 {
		path = ret[1]
	}

	return path + ".md"
}
