package config

import (
	"fmt"
)

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
	Product          string
}

// multiCatalogKeys is a map of primary and derived catalog keys for services with multiple clients.
// If we add another version of a service client, don't forget to update it.
var multiCatalogKeys = map[string][]string{
	"iam":          {"identity", "identity_ext", "iam_no_version"},
	"bss":          {"bssv2"},
	"ecs":          {"ecsv21", "ecsv11"},
	"evs":          {"evsv21", "evsv1", "evsv5"},
	"cce":          {"ccev1", "cce_addon"},
	"cci":          {"cciv1_bata"},
	"vpc":          {"networkv2", "vpcv3", "fwv2"},
	"elb":          {"elbv2", "elbv3"},
	"dns":          {"dns_region", "dnsv21"},
	"dds":          {"ddsv31"},
	"kms":          {"kmsv1", "kmsv3"},
	"mrs":          {"mrsv2"},
	"nat":          {"natv3"},
	"rds":          {"rdsv1", "rdsv31"},
	"waf":          {"waf-dedicated"},
	"geminidb":     {"geminidbv31"},
	"dataarts":     {"dataarts-dlf"},
	"dli":          {"dliv2", "dliv3"},
	"dcs":          {"dcsv1"},
	"dis":          {"disv3"},
	"dms":          {"dmsv2"},
	"drs":          {"drsv5"},
	"dws":          {"dwsv2"},
	"apig":         {"apigv2"},
	"modelarts":    {"modelartsv2"},
	"servicestage": {"servicestagev2"},
	"smn":          {"smn-tag"},
	"ces":          {"cesv2"},
	"ims":          {"imsv1"},
	"config":       {"rms"}, // config is named as Resource Management Service(RMS) before
	"tms":          {"tmsv2"},
	"anti-ddos":    {"anti-ddosv2"},
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
		Admin:            true,
		WithOutProjectID: true,
		Product:          "IAM",
	},
	"identity_ext": {
		Name:             "iam",
		Version:          "v3-ext",
		Admin:            true,
		WithOutProjectID: true,
		Product:          "IAM",
	},
	"iam_no_version": {
		Name:             "iam",
		Version:          "",
		Admin:            true,
		WithOutProjectID: true,
		Product:          "IAM",
	},
	"accessanalyzer": {
		Name:             "access-analyzer",
		Version:          "",
		Admin:            true,
		WithOutProjectID: true,
		Product:          "AccessAnalyzer",
	},

	// iam is used for huaweicloud IAM APIs
	"iam": {
		Name:             "iam",
		Version:          "v3.0",
		Admin:            true,
		WithOutProjectID: true,
		Product:          "IAM",
	},
	"sts": {
		Name:             "sts",
		Version:          "",
		Admin:            true,
		WithOutProjectID: true,
		Product:          "IAM",
	},
	"identitycenter": {
		Name:             "identitycenter",
		Version:          "v1",
		Scope:            "global",
		WithOutProjectID: true,
		Product:          "IdentityCenter",
	},
	"identitystore": {
		Name:             "identitystore",
		Version:          "v1",
		WithOutProjectID: true,
		Product:          "IdentityCenter",
	},
	"cdn": {
		Name:             "cdn",
		Version:          "v1.0",
		Scope:            "global",
		WithOutProjectID: true,
		Product:          "CDN",
	},
	"eps": {
		Name:             "eps",
		Version:          "v1.0",
		Scope:            "global",
		Admin:            true,
		WithOutProjectID: true,
		Product:          "EPS",
	},
	"bss": {
		Name:             "bss",
		Version:          "v1.0",
		Scope:            "global",
		WithOutProjectID: true,
		Product:          "BSS",
	},
	"bssv2": {
		Name:             "bss",
		Version:          "v2",
		Scope:            "global",
		WithOutProjectID: true,
		Product:          "BSS",
	},
	// ******* catalog for Compute *******
	"ecs": {
		Name:    "ecs",
		Version: "v1",
		Product: "ECS",
	},
	"ecsv11": {
		Name:    "ecs",
		Version: "v1.1",
		Product: "ECS",
	},
	"ecsv21": {
		Name:    "ecs",
		Version: "v2.1",
		Product: "ECS",
	},
	"autoscaling": {
		Name:    "as",
		Version: "autoscaling-api/v1",
		Product: "AS",
	},
	"imsv1": {
		Name:             "ims",
		Version:          "v1",
		WithOutProjectID: true,
		Product:          "IMS",
	},
	"ims": {
		Name:             "ims",
		Version:          "v2",
		WithOutProjectID: true,
		Product:          "IMS",
	},
	"cms": {
		Name:             "cms",
		Version:          "v2",
		Scope:            "global",
		WithOutProjectID: true,
		Product:          "CMS",
	},

	// ******* catalog for Container *******
	"ccev1": {
		Name:             "cce",
		Version:          "api/v1",
		WithOutProjectID: true,
		Product:          "CCE",
	},
	"cce": {
		Name:    "cce",
		Version: "api/v3/projects",
		Product: "CCE",
	},
	"cce_addon": {
		Name:             "cce",
		Version:          "api/v3",
		WithOutProjectID: true,
		Product:          "CCE",
	},
	"swr": {
		Name:             "swr-api",
		Version:          "v2",
		WithOutProjectID: true,
		Product:          "SWR",
	},
	"cci": {
		Name:             "cci",
		Version:          "api/v1",
		WithOutProjectID: true,
		Product:          "CCI",
	},
	"cciv1_bata": {
		Name:             "cci",
		Version:          "apis/networking.cci.io/v1beta1",
		WithOutProjectID: true,
		Product:          "CCI",
	},
	"ucs": {
		Name:             "ucs",
		Version:          "v1",
		Scope:            "global",
		Admin:            true,
		WithOutProjectID: true,
		Product:          "UCS",
	},
	"asm": {
		Name:             "asm",
		Version:          "v1",
		WithOutProjectID: true,
		Product:          "ASM",
	},

	"aom": {
		Name:    "aom",
		Version: "svcstg/icmgr/v1",
		Product: "AOM",
	},
	"coc": {
		Name:    "coc",
		Version: "v1",
		Scope:   "global",
		Product: "COC",
	},
	"fgs": {
		Name:    "functiongraph",
		Version: "v2",
		Product: "FunctionGraph",
	},
	"bms": {
		Name:    "bms",
		Version: "v1",
		Product: "BMS",
	},
	"rfs": {
		Name:    "rfs",
		Version: "v1",
		Product: "RFS",
	},

	// ******* catalog for storage ******
	"evsv1": {
		Name:    "evs",
		Version: "v1",
		Product: "EVS",
	},
	"evs": {
		Name:    "evs",
		Version: "v2",
		Product: "EVS",
	},
	"evsv21": {
		Name:    "evs",
		Version: "v2.1",
		Product: "EVS",
	},
	"evsv5": {
		Name:    "evs",
		Version: "v5",
		Product: "EVS",
	},
	"sfs": {
		Name:    "sfs",
		Version: "v2",
		Product: "SFS",
	},
	"sfs-turbo": {
		Name:    "sfs-turbo",
		Version: "v1",
		Product: "SFSTurbo",
	},
	"cbh": {
		Name:    "cbh",
		Version: "v1",
		Product: "CBH",
	},
	"cbr": {
		Name:    "cbr",
		Version: "v3",
		Product: "CBR",
	},
	"csbs": {
		Name:    "csbs",
		Version: "v1",
		Product: "CSBS",
	},
	"vbs": {
		Name:    "vbs",
		Version: "v2",
		Product: "VBS",
	},
	"sdrs": {
		Name:    "sdrs",
		Version: "v1",
		Product: "SDRS",
	},

	// ******* catalog for network ******
	"vpc": {
		Name:             "vpc",
		Version:          "v1",
		WithOutProjectID: true,
		Product:          "VPC",
	},
	"networkv2": {
		Name:             "vpc",
		Version:          "v2.0",
		WithOutProjectID: true,
		Product:          "VPC",
	},
	"vpcv3": {
		Name:    "vpc",
		Version: "v3",
		Product: "VPC",
	},
	"geip": {
		Name:             "eip",
		Version:          "v3",
		Scope:            "global",
		WithOutProjectID: true,
		Product:          "EIP",
	},
	"nat": {
		Name:    "nat",
		Version: "v2",
		Product: "NAT",
	},
	"natv3": {
		Name:    "nat",
		Version: "v3",
		Product: "NAT",
	},
	"elbv2": {
		Name:             "elb",
		Version:          "v2.0",
		WithOutProjectID: true,
		Product:          "ELB",
	},
	"elbv3": {
		Name:    "elb",
		Version: "v3",
		Product: "ELB",
	},
	"elb": {
		Name:    "elb",
		Version: "v2",
		Product: "ELB",
	},
	"fwv2": {
		Name:             "vpc",
		Version:          "v2.0",
		WithOutProjectID: true,
		Product:          "VPC",
	},
	"vpcep": {
		Name:    "vpcep",
		Version: "v1",
		Product: "VPCEP",
	},
	"dns": {
		Name:             "dns",
		Version:          "v2",
		Scope:            "global",
		WithOutProjectID: true,
		Product:          "DNS",
	},
	"dns_region": {
		Name:             "dns",
		Version:          "v2",
		WithOutProjectID: true,
		Product:          "DNS",
	},
	"dnsv21": {
		Name:             "dns",
		Version:          "v2.1",
		WithOutProjectID: true,
		Product:          "DNS",
	},
	"workspace": {
		Name:    "workspace",
		Version: "v2",
		Product: "Workspace",
	},
	"appstream": {
		Name:    "appstream",
		Version: "v1",
		Product: "Workspace",
	},
	"er": {
		Name:    "er",
		Version: "v3",
		Product: "ER",
	},
	"vpn": {
		Name:    "vpn",
		Version: "v5",
		Product: "VPN",
	},
	"ga": {
		Name:             "ga",
		Version:          "v1",
		WithOutProjectID: true,
		Scope:            "global",
		Product:          "GA",
	},
	"dc": {
		Name:    "dcaas",
		Version: "v3",
		Product: "DC",
	},
	"cfw": {
		Name:    "cfw",
		Version: "v1",
		Product: "CFW",
	},

	// catalog for database
	"rdsv1": {
		Name:    "rds",
		Version: "rds/v1",
		Product: "RDS",
	},
	"rds": {
		Name:    "rds",
		Version: "v3",
		Product: "RDS",
	},
	"rdsv31": {
		Name:    "rds",
		Version: "v3.1",
		Product: "RDS",
	},
	"ram": {
		Name:             "ram",
		Version:          "v1",
		WithOutProjectID: true,
		Scope:            "global",
		Product:          "RAM",
	},
	"dds": {
		Name:    "dds",
		Version: "v3",
		Product: "DDS",
	},
	"ddsv31": {
		Name:    "dds",
		Version: "v3.1",
		Product: "DDS",
	},
	"deh": {
		Name:    "deh",
		Version: "v1.0",
		Product: "DEH",
	},
	"geminidb": {
		Name:    "gaussdb-nosql",
		Version: "v3",
		Product: "GaussDBforNoSQL",
	},
	"geminidbv31": {
		Name:    "gaussdb-nosql",
		Version: "v3.1",
		Product: "GaussDBforNoSQL",
	},
	"gaussdb": {
		Name:    "gaussdb",
		Version: "v3",
		Product: "GaussDBforMySQL",
	},
	"opengauss": {
		Name:    "gaussdb-opengauss",
		Version: "v3",
		Product: "GaussDB",
	},
	"drs": {
		Name:    "drs",
		Version: "v3",
		Product: "DRS",
	},
	"drsv5": {
		Name:    "drs",
		Version: "v5",
		Product: "DRS",
	},

	// catalog for management service
	"ces": {
		Name:    "ces",
		Version: "V1.0",
		Product: "CES",
	},
	"cesv2": {
		Name:    "ces",
		Version: "v2",
		Product: "CES",
	},
	"cts": {
		Name:    "cts",
		Version: "v1.0",
		Product: "CTS",
	},
	"lts": {
		Name:    "lts",
		Version: "v2",
		Product: "LTS",
	},
	"apm": {
		Name:             "apm2",
		Version:          "v1",
		Product:          "APM",
		WithOutProjectID: true,
	},
	"smn": {
		Name:         "smn",
		Version:      "v2",
		ResourceBase: "notifications",
		Product:      "SMN",
	},
	"smn-tag": {
		Name:    "smn",
		Version: "v2",
		Product: "SMN",
	},
	"sms-intl": {
		Name:             "sms.ap-southeast-3",
		Version:          "v3",
		Scope:            "global",
		WithOutProjectID: true,
		Product:          "SMS",
	},
	"sms": {
		Name:             "sms.cn-north-4",
		Version:          "v3",
		Scope:            "global",
		WithOutProjectID: true,
		Product:          "SMS",
	},
	"tms": {
		Name:             "tms",
		Version:          "v1.0",
		Scope:            "global",
		Admin:            true, // The 'X-Domain-Id' is required for TMS service requests.
		WithOutProjectID: true,
		Product:          "TMS",
	},
	"tmsv2": {
		Name:             "tms",
		Version:          "v2.0",
		Scope:            "global",
		Admin:            true, // The 'X-Domain-Id' is required for TMS service requests.
		WithOutProjectID: true,
		Product:          "TMS",
	},
	"rms": {
		Name:             "rms",
		Scope:            "global",
		Version:          "v1",
		WithOutProjectID: true,
		Product:          "Config",
	},
	"rgc": {
		Name:             "rgc",
		Version:          "v1",
		WithOutProjectID: true,
		Product:          "RGC",
	},
	"organizations": {
		Name:             "organizations",
		Version:          "v1",
		Scope:            "global",
		WithOutProjectID: true,
		Product:          "Organizations",
	},
	// catalog for Meeting service, only used for API scan
	"meeting": {
		Name:             "api.meeting",
		Version:          "v1",
		Scope:            "global",
		WithOutProjectID: true,
		Product:          "Meeting",
	},

	// catalog for Security service
	"aad": {
		Name:             "aad",
		Version:          "v1",
		Scope:            "global",
		WithOutProjectID: true,
		Product:          "AAD",
	},
	"anti-ddos": {
		Name:    "antiddos",
		Version: "v1",
		Product: "Anti-DDoS",
	},
	"anti-ddosv2": {
		Name:    "antiddos",
		Version: "v2",
		Product: "Anti-DDoS",
	},
	"kms": {
		Name:             "kms",
		Version:          "v1.0",
		WithOutProjectID: true,
		Product:          "DEW",
	},
	"kmsv1": {
		Name:    "kms",
		Version: "v1",
		Product: "DEW",
	},
	"kmsv3": {
		Name:    "kms",
		Version: "v3",
		Product: "DEW",
	},
	"waf": {
		Name:         "waf",
		Version:      "v1",
		ResourceBase: "waf",
		Product:      "WAF",
	},
	"waf-dedicated": {
		Name:         "waf",
		Version:      "v1",
		ResourceBase: "premium-waf",
		Product:      "WAF",
	},
	"dbss": {
		Name:    "dbss",
		Version: "v2",
		Product: "DBSS",
	},
	"hss": {
		Name:    "hss",
		Version: "v5",
		Product: "HSS",
	},
	"secmaster": {
		Name:    "secmaster",
		Version: "v1",
		Product: "SecMaster",
	},

	// catalog for Enterprise Intelligence
	"cae": {
		Name:    "cae",
		Product: "CAE",
	},
	"mrs": {
		Name:    "mrs",
		Version: "v1.1",
		Product: "MRS",
	},
	"mrsv2": {
		Name:    "mrs",
		Version: "v2",
		Product: "MRS",
	},
	"modelarts": {
		Name:    "modelarts",
		Version: "v1",
		Product: "ModelArts",
	},
	"modelartsv2": {
		Name:    "modelarts",
		Version: "v2",
		Product: "ModelArts",
	},
	"dataarts": {
		Name:    "dayu",
		Version: "v1",
		Product: "DataArtsStudio",
	},
	"dataarts-dlf": {
		Name:    "dayu-dlf",
		Version: "v1",
		Product: "DataArtsStudio",
	},
	"dws": {
		Name:    "dws",
		Version: "v1.0",
		Product: "DWS",
	},
	"dwsv2": {
		Name:    "dws",
		Version: "v2",
		Product: "DWS",
	},
	"dli": {
		Name:    "dli",
		Version: "v1.0",
		Product: "DLI",
	},
	"dliv2": {
		Name:    "dli",
		Version: "v2.0",
		Product: "DLI",
	},
	"dliv3": {
		Name:    "dli",
		Version: "v3",
		Product: "DLI",
	},
	"dis": {
		Name:    "dis",
		Version: "v2",
		Product: "DIS",
	},
	"disv3": {
		Name:    "dis",
		Version: "v3",
		Product: "DIS",
	},
	"css": {

		Name:    "css",
		Version: "v1.0",
		Product: "CSS",
	},
	"cssv2": {

		Name:    "css",
		Version: "v2.0",
		Product: "CSS",
	},
	"cs": {
		Name:    "cs",
		Version: "v1.0",
		Product: "CloudStream",
	},
	"ges": {
		Name:    "ges",
		Version: "v1.0",
		Product: "GES",
	},
	"cloudtable": {
		Name:    "cloudtable",
		Version: "v2",
		Product: "CloudTable",
	},
	"cdm": {
		Name:    "cdm",
		Version: "v1.1",
		Product: "CDM",
	},

	// catalog for Application
	"apig": {
		Name:             "apig",
		Version:          "v1.0",
		ResourceBase:     "apigw",
		WithOutProjectID: true,
		Product:          "APIG",
	},
	"apigv2": {
		Name:         "apig",
		Version:      "v2",
		ResourceBase: "apigw",
		Product:      "APIG",
	},
	"bcs": {
		Name:    "bcs",
		Version: "v2",
		Product: "BCS",
	},
	"cse": {
		Name:    "cse",
		Version: "v2",
		Product: "CSE",
	},
	"dcsv1": {
		Name:             "dcs",
		Version:          "v1.0",
		WithOutProjectID: true,
		Product:          "DCS",
	},
	"dcs": {
		Name:             "dcs",
		Version:          "v2",
		WithOutProjectID: true,
		Product:          "DCS",
	},
	"dms": {
		Name:             "dms",
		Version:          "v1.0",
		WithOutProjectID: true,
		Product:          "DMS",
	},
	"dmsv2": {
		Name:             "dms",
		Version:          "v2",
		WithOutProjectID: true,
		Product:          "DMS",
	},
	"servicestage": {
		Name:    "servicestage",
		Version: "v1",
		Product: "ServiceStage",
	},
	"servicestagev2": {
		Name:    "servicestage",
		Version: "v2",
		Product: "ServiceStage",
	},
	"eg": {
		Name:    "eg",
		Version: "v1",
		Product: "EG",
	},

	// catalog for IEC which is a global service
	"iec": {
		Name:             "iecs",
		Version:          "v1",
		Scope:            "global",
		Admin:            true,
		WithOutProjectID: true,
		Product:          "IEC",
	},

	// catalog for Others
	"rts": {
		Name:    "rts",
		Version: "v1",
		Product: "RTS",
	},
	"oms": {
		Name:    "oms",
		Version: "v1",
		Product: "OMS",
	},
	"scm": {
		Name:             "scm",
		Version:          "v3",
		WithOutProjectID: true,
		Product:          "CCM",
	},
	"ccm": {
		Name:             "ccm",
		Version:          "v1",
		WithOutProjectID: true,
		Product:          "CCM",
	},

	// catalog for cc
	"cc": {
		Name:             "cc",
		Version:          "v3",
		Scope:            "global",
		WithOutProjectID: true,
		Admin:            true,
		Product:          "CC",
	},

	"cpts": {
		Name:    "cpts",
		Version: "v1",
		Product: "CPTS",
	},

	"live": {
		Name:    "live",
		Version: "v1",
		Product: "Live",
	},

	"mpc": {
		Name:    "mpc",
		Version: "v1",
		Product: "MPC",
	},

	"iotda": {
		Name:    "iotda",
		Version: "v5",
		Product: "IoTDA",
	},

	"vod": {
		Name:    "vod",
		Version: "v1",
		Product: "VOD",
	},

	"cmdb": {
		Name:    "cmdb",
		Version: "v1",
		Scope:   "global",
		Product: "AOM",
	},

	"ddm": {
		Name:             "ddm",
		WithOutProjectID: true,
		Product:          "DDM",
	},

	// catalog for Developer Services
	"codehub": {
		Name:    "codehub-ext",
		Product: "CodeHub",
	},

	"projectman": {
		Name:    "projectman-ext",
		Version: "v4",
		Product: "ProjectMan",
	},

	"codearts_deploy": {
		Name:    "codearts-deploy",
		Version: "v2",
		Product: "CodeArtsDeploy",
	},

	"vss": {
		Name:    "vss",
		Version: "v3",
		Scope:   "global",
		Product: "CodeArtsInspector",
	},

	"codearts_pipeline": {
		Name:    "cloudpipeline-ext",
		Version: "v5",
		Product: "CodeArtsPipeline",
	},

	// catalog for Data Security Center
	"dsc": {
		Name:    "sdg",
		Product: "DSC",
	},

	// catalog for Cloud Phone
	"cph": {
		Name:    "cph",
		Product: "CPH",
	},

	// catalog for Joint-Operation Cloud only
	// it should be at the end of this map, and no necessary to put the key into allServiceCatalog
	"mls": {
		Name:    "mls",
		Version: "v1.0",
		Product: "MLS",
	},
	"natv2": {
		Name:             "nat",
		Version:          "v2.0",
		WithOutProjectID: true,
		Product:          "NAT",
	},

	// catalog for KooGallery
	"mkt": {
		Name:             "mkt",
		Version:          "api/mkp-openapi-public/v1",
		Scope:            "global",
		WithOutProjectID: true,
		Product:          "KooGallery",
	},
}

// GetServiceEndpoint try to get the endpoint from customizing map
func GetServiceEndpoint(c *Config, srv, region string) string {
	if endpoint, ok := c.Endpoints[srv]; ok {
		return endpoint
	}

	// get the endpoint from build-in catalog
	catalog, ok := allServiceCatalog[srv]
	if !ok {
		return ""
	}

	// update the service catalog name if necessary
	if name := getServiceCatalogNameByRegion(srv, region); name != "" {
		catalog.Name = name
	}

	var ep string
	if catalog.Scope == "global" && !c.RegionClient {
		ep = fmt.Sprintf("https://%s.%s/", catalog.Name, c.Cloud)
	} else {
		ep = fmt.Sprintf("https://%s.%s.%s/", catalog.Name, region, c.Cloud)
	}
	return ep
}

// GetServiceCatalog returns the catalog object of a service
func GetServiceCatalog(service string) *ServiceCatalog {
	if catalog, ok := allServiceCatalog[service]; ok {
		return &catalog
	}
	return nil
}
