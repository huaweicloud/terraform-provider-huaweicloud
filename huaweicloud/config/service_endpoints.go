package config

// For cloud services like GaussDB, the endpoints are different in different regions.
// Therefore, we have add a map to manage their endpoints.
// please refer to https://developer.huaweicloud.com/intl/en-us/endpoint
var serviceRegionCatalogName = map[string]interface{}{
	// the prefix(catalog name) of GaussDB is "gaussdb" in other regions
	"gaussdb": map[string]string{
		"ap-southeast-1": "gaussdbformysql",
		"ap-southeast-2": "gaussdbformysql",
		"cn-north-11":    "gaussdbformysql",
		"tr-west-1":      "gaussdbformysql",
		"af-south-1":     "gaussdb-mysql",
	},
}

func getServiceCatalogNameByRegion(service, region string) string {
	if v, ok := serviceRegionCatalogName[service]; ok {
		if regionMap, ok := v.(map[string]string); ok {
			return regionMap[region]
		}
	}

	return ""
}
