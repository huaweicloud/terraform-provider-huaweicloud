package config

// ServiceCatalog defines a struct which was used to generate a service client for huaweicloud.
// the endpoint likes https://{Name}.{Region}.myhuaweicloud.com/{Version}/{project_id}/{ResourceBase}
// For more information, please refer to Config.NewServiceClient
type ServiceCatalog struct {
	Name             string
	Version          string
	Scope            string
	Admin            bool
	ResourceBase     string
	WithOutProjectID bool
}

// multiCatalogKeys is a map of primary and derived catalog keys for services with multiple clients.
// If we add another version of a service client, don't forget to update it.
var multiCatalogKeys = map[string][]string{
	"iam":      {"identity", "iam_no_version"},
	"bss":      {"bssv2"},
	"ecs":      {"ecsv21", "ecsv11"},
	"evs":      {"evsv21"},
	"cce":      {"ccev1", "cce_addon"},
	"cci":      {"cciv1_bata"},
	"vpc":      {"networkv2", "vpcv3", "security_group", "fwv2"},
	"elb":      {"elbv2", "elbv3"},
	"dns":      {"dns_region"},
	"kms":      {"kmsv1"},
	"mrs":      {"mrsv2"},
	"rds":      {"rdsv1"},
	"waf":      {"waf-dedicated"},
	"geminidb": {"geminidbv31"},
	"dli":      {"dliv2"},
	"dcs":      {"dcsv1"},
	"dis":      {"disv3"},
	"dms":      {"dmsv2"},
	"dws":      {"dwsv2"},
	"apig":     {"apigv2"},
}

// GetServiceDerivedCatalogKeys returns the derived catalog keys of a service.
func GetServiceDerivedCatalogKeys(mainKey string) []string {
	return multiCatalogKeys[mainKey]
}

var allServiceCatalog = map[string]ServiceCatalog{
	// catalog for global service
	// identity is used for openstack keystone APIs
	"identity": {
		Name:             "iam",
		Version:          "v3",
		Scope:            "global",
		Admin:            true,
		WithOutProjectID: true,
	},
	"iam_no_version": {
		Name:             "iam",
		Version:          "",
		Scope:            "global",
		Admin:            true,
		WithOutProjectID: true,
	},
	// iam is used for huaweicloud IAM APIs
	"iam": {
		Name:             "iam",
		Version:          "v3.0",
		Scope:            "global",
		Admin:            true,
		WithOutProjectID: true,
	},
	"cdn": {
		Name:             "cdn",
		Version:          "v1.0",
		Scope:            "global",
		WithOutProjectID: true,
	},
	"eps": {
		Name:             "eps",
		Version:          "v1.0",
		Scope:            "global",
		Admin:            true,
		WithOutProjectID: true,
	},
	"bss": {
		Name:             "bss",
		Version:          "v1.0",
		Scope:            "global",
		WithOutProjectID: true,
	},
	"bssv2": {
		Name:             "bss",
		Version:          "v2",
		Scope:            "global",
		WithOutProjectID: true,
	},

	// ******* catalog for Compute *******
	"ecs": {
		Name:    "ecs",
		Version: "v1",
	},
	"ecsv11": {
		Name:    "ecs",
		Version: "v1.1",
	},
	"ecsv21": {
		Name:    "ecs",
		Version: "v2.1",
	},
	"autoscaling": {
		Name:    "as",
		Version: "autoscaling-api/v1",
	},
	"ims": {
		Name:             "ims",
		Version:          "v2",
		WithOutProjectID: true,
	},
	"ccev1": {
		Name:             "cce",
		Version:          "api/v1",
		WithOutProjectID: true,
	},
	"cce": {
		Name:    "cce",
		Version: "api/v3/projects",
	},
	"cce_addon": {
		Name:             "cce",
		Version:          "api/v3",
		WithOutProjectID: true,
	},
	"aom": {
		Name:    "aom",
		Version: "svcstg/icmgr/v1",
	},
	"cci": {
		Name:             "cci",
		Version:          "api/v1",
		WithOutProjectID: true,
	},
	"cciv1_bata": {
		Name:             "cci",
		Version:          "apis/networking.cci.io/v1beta1",
		WithOutProjectID: true,
	},
	"fgs": {
		Name:    "functiongraph",
		Version: "v2",
	},
	"swr": {
		Name:             "swr-api",
		Version:          "v2",
		WithOutProjectID: true,
	},
	"bms": {
		Name:    "bms",
		Version: "v1",
	},

	// ******* catalog for storage ******
	"evs": {
		Name:    "evs",
		Version: "v2",
	},
	"evsv21": {
		Name:    "evs",
		Version: "v2.1",
	},
	"sfs": {
		Name:    "sfs",
		Version: "v2",
	},
	"sfs-turbo": {
		Name:    "sfs-turbo",
		Version: "v1",
	},
	"cbr": {
		Name:    "cbr",
		Version: "v3",
	},
	"csbs": {
		Name:    "csbs",
		Version: "v1",
	},
	"vbs": {
		Name:    "vbs",
		Version: "v2",
	},

	// ******* catalog for network ******
	"vpc": {
		Name:             "vpc",
		Version:          "v1",
		WithOutProjectID: true,
	},
	"networkv2": {
		Name:             "vpc",
		Version:          "v2.0",
		WithOutProjectID: true,
	},
	"vpcv3": {
		Name:    "vpc",
		Version: "v3",
	},
	"security_group": {
		Name:    "vpc",
		Version: "v1",
	},
	"nat": {
		Name:    "nat",
		Version: "v2",
	},
	"elbv2": {
		Name:             "elb",
		Version:          "v2.0",
		WithOutProjectID: true,
	},
	"elbv3": {
		Name:    "elb",
		Version: "v3",
	},
	"elb": {
		Name:    "elb",
		Version: "v2",
	},
	"fwv2": {
		Name:             "vpc",
		Version:          "v2.0",
		WithOutProjectID: true,
	},
	"vpcep": {
		Name:    "vpcep",
		Version: "v1",
	},
	"dns": {
		Name:             "dns",
		Version:          "v2",
		Scope:            "global",
		WithOutProjectID: true,
	},
	"dns_region": {
		Name:             "dns",
		Version:          "v2",
		WithOutProjectID: true,
	},

	// catalog for database
	"rdsv1": {
		Name:    "rds",
		Version: "rds/v1",
	},
	"rds": {
		Name:    "rds",
		Version: "v3",
	},
	"dds": {
		Name:    "dds",
		Version: "v3",
	},
	"geminidb": {
		Name:    "gaussdb-nosql",
		Version: "v3",
	},
	"geminidbv31": {
		Name:    "gaussdb-nosql",
		Version: "v3.1",
	},
	"gaussdb": {
		Name:    "gaussdb",
		Version: "mysql/v3",
	},
	"opengauss": {
		Name:    "gaussdb-opengauss",
		Version: "v3",
	},
	"drs": {
		Name:    "drs",
		Version: "v3",
	},

	// catalog for management service
	"ces": {
		Name:    "ces",
		Version: "V1.0",
	},
	"cts": {
		Name:    "cts",
		Version: "v1.0",
	},
	"lts": {
		Name:    "lts",
		Version: "v2",
	},
	"smn": {
		Name:         "smn",
		Version:      "v2",
		ResourceBase: "notifications",
	},

	// catalog for Security service
	"anti-ddos": {
		Name:    "antiddos",
		Version: "v1",
	},
	"kms": {
		Name:             "kms",
		Version:          "v1.0",
		WithOutProjectID: true,
	},
	"kmsv1": {
		Name:    "kms",
		Version: "v1",
	},
	"waf": {
		Name:         "waf",
		Version:      "v1",
		ResourceBase: "waf",
	},
	"waf-dedicated": {
		Name:         "waf",
		Version:      "v1",
		ResourceBase: "premium-waf",
	},

	// catalog for Enterprise Intelligence
	"mrs": {
		Name:    "mrs",
		Version: "v1.1",
	},
	"mrsv2": {
		Name:    "mrs",
		Version: "v2",
	},
	"modelarts": {
		Name:    "modelarts",
		Version: "v1",
	},
	"dws": {
		Name:    "dws",
		Version: "v1.0",
	},
	"dwsv2": {
		Name:    "dws",
		Version: "v2",
	},
	"dli": {
		Name:    "dli",
		Version: "v1.0",
	},
	"dliv2": {
		Name:    "dli",
		Version: "v2.0",
	},
	"dis": {
		Name:    "dis",
		Version: "v2",
	},
	"disv3": {
		Name:    "dis",
		Version: "v3",
	},
	"css": {
		Name:    "css",
		Version: "v1.0",
	},
	"cs": {
		Name:    "cs",
		Version: "v1.0",
	},
	"ges": {
		Name:    "ges",
		Version: "v1.0",
	},
	"cloudtable": {
		Name:    "cloudtable",
		Version: "v2",
	},
	"cdm": {
		Name:    "cdm",
		Version: "v1.1",
	},

	// catalog for Application
	"apig": {
		Name:             "apig",
		Version:          "v1.0",
		ResourceBase:     "apigw",
		WithOutProjectID: true,
	},
	"apigv2": {
		Name:         "apig",
		Version:      "v2",
		ResourceBase: "apigw",
	},
	"bcs": {
		Name:    "bcs",
		Version: "v2",
	},
	"dcsv1": {
		Name:    "dcs",
		Version: "v1.0",
	},
	"dcs": {
		Name:    "dcs",
		Version: "v2",
	},
	"dms": {
		Name:    "dms",
		Version: "v1.0",
	},
	"dmsv2": {
		Name:    "dms",
		Version: "v2",
	},

	// catalog for IEC which is a global service
	"iec": {
		Name:             "iecs",
		Version:          "v1",
		Scope:            "global",
		Admin:            true,
		WithOutProjectID: true,
	},

	// catalog for Others
	"rts": {
		Name:    "rts",
		Version: "v1",
	},
	"oms": {
		Name:    "oms",
		Version: "v1",
	},
	"mls": {
		Name:    "mls",
		Version: "v1.0",
	},
	"scm": {
		Name:             "scm",
		Version:          "v3",
		WithOutProjectID: true,
	},
}
