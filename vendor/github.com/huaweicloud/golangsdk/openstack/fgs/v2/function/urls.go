package function

import "github.com/huaweicloud/golangsdk"

const (
	FGS      = "fgs"
	FUNCTION = "functions"
	CODE     = "code"
	CONFIG   = "config"
	VERSION  = "versions"
	ALIAS    = "aliases"
	INVOKE   = "invocations"
	ASINVOKE = "invocations-async"
)

func createURL(c *golangsdk.ServiceClient) string {
	return listURL(c)
}

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(FGS, FUNCTION)
}

func deleteURL(c *golangsdk.ServiceClient, functionUrn string) string {
	return c.ServiceURL(FGS, FUNCTION, functionUrn)
}

//function code
func getCodeURL(c *golangsdk.ServiceClient, functionUrn string) string {
	return c.ServiceURL(FGS, FUNCTION, functionUrn, CODE)
}

func updateCodeURL(c *golangsdk.ServiceClient, functionUrn string) string {
	return getCodeURL(c, functionUrn)
}

//function metadata
func getMetadataURL(c *golangsdk.ServiceClient, functionUrn string) string {
	return c.ServiceURL(FGS, FUNCTION, functionUrn, CONFIG)
}

func updateMetadataURL(c *golangsdk.ServiceClient, functionUrn string) string {
	return getMetadataURL(c, functionUrn)
}

//function invoke
func invokeURL(c *golangsdk.ServiceClient, functionUrn string) string {
	return c.ServiceURL(FGS, FUNCTION, functionUrn, INVOKE)
}

func asyncInvokeURL(c *golangsdk.ServiceClient, functionUrn string) string {
	return c.ServiceURL(FGS, FUNCTION, functionUrn, ASINVOKE)
}

//function version
func createVersionURL(c *golangsdk.ServiceClient, functionUrn string) string {
	return c.ServiceURL(FGS, FUNCTION, functionUrn, VERSION)
}

func listVersionURL(c *golangsdk.ServiceClient, functionUrn string) string {
	return createVersionURL(c, functionUrn)
}

//function alias
func createAliasURL(c *golangsdk.ServiceClient, functionUrn string) string {
	return c.ServiceURL(FGS, FUNCTION, functionUrn, ALIAS)
}

func updateAliasURL(c *golangsdk.ServiceClient, functionUrn, aliasName string) string {
	return c.ServiceURL(FGS, FUNCTION, functionUrn, ALIAS, aliasName)
}

func deleteAliasURL(c *golangsdk.ServiceClient, functionUrn, aliasName string) string {
	return updateAliasURL(c, functionUrn, aliasName)
}

func getAliasURL(c *golangsdk.ServiceClient, functionUrn, aliasName string) string {
	return updateAliasURL(c, functionUrn, aliasName)
}

func listAliasURL(c *golangsdk.ServiceClient, functionUrn string) string {
	return createAliasURL(c, functionUrn)
}
