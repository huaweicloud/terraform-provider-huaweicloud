package huaweicloud

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

var allServiceCatalog = map[string]ServiceCatalog{
	// catalog for global service
	// identity is used for openstack keystone APIs
	"identity": ServiceCatalog{
		Name:             "iam",
		Version:          "v3",
		Scope:            "global",
		Admin:            true,
		WithOutProjectID: true,
	},
	// iam is used for huaweicloud IAM APIs
	"iam": ServiceCatalog{
		Name:             "iam",
		Version:          "v3.0",
		Scope:            "global",
		Admin:            true,
		WithOutProjectID: true,
	},
	"cdn": ServiceCatalog{
		Name:             "cdn",
		Version:          "v1.0",
		Scope:            "global",
		WithOutProjectID: true,
	},
	"eps": ServiceCatalog{
		Name:             "eps",
		Version:          "v1.0",
		Scope:            "global",
		Admin:            true,
		WithOutProjectID: true,
	},
	"bss": ServiceCatalog{
		Name:             "bss",
		Version:          "v1.0",
		Scope:            "global",
		WithOutProjectID: true,
	},
	"bssv2": ServiceCatalog{
		Name:             "bss",
		Version:          "v2",
		Scope:            "global",
		WithOutProjectID: true,
	},

	// ******* catalog for Compute *******
	"ecs": ServiceCatalog{
		Name:    "ecs",
		Version: "v1",
	},
	"ecsv11": ServiceCatalog{
		Name:    "ecs",
		Version: "v1.1",
	},
	"ecsv21": ServiceCatalog{
		Name:    "ecs",
		Version: "v2.1",
	},
	"autoscaling": ServiceCatalog{
		Name:    "as",
		Version: "autoscaling-api/v1",
	},
	"ims": ServiceCatalog{
		Name:             "ims",
		Version:          "v2",
		WithOutProjectID: true,
	},
	"cce": ServiceCatalog{
		Name:    "cce",
		Version: "api/v3/projects",
	},
	"cce_addon": ServiceCatalog{
		Name:             "cce",
		Version:          "api/v3",
		WithOutProjectID: true,
	},
	"cciv1": ServiceCatalog{
		Name:             "cci",
		Version:          "apis/networking.cci.io/v1beta1",
		WithOutProjectID: true,
	},
	"fgsv2": ServiceCatalog{
		Name:    "functiongraph",
		Version: "v2",
	},

	// ******* catalog for storage ******
	"volumev2": ServiceCatalog{
		Name:    "evs",
		Version: "v2",
	},
	"evs": ServiceCatalog{
		Name:    "evs",
		Version: "v3",
	},
	"sfs": ServiceCatalog{
		Name:    "sfs",
		Version: "v2",
	},
	"sfs-turbo": ServiceCatalog{
		Name:    "sfs-turbo",
		Version: "v1",
	},
	"csbs": ServiceCatalog{
		Name:    "csbs",
		Version: "v1",
	},
	"vbs": ServiceCatalog{
		Name:    "vbs",
		Version: "v2",
	},

	// ******* catalog for network ******
	"vpc": ServiceCatalog{
		Name:             "vpc",
		Version:          "v1",
		WithOutProjectID: true,
	},
	"networkv2": ServiceCatalog{
		Name:             "vpc",
		Version:          "v2.0",
		WithOutProjectID: true,
	},
	"security_group": ServiceCatalog{
		Name:    "vpc",
		Version: "v1",
	},
	"natv2": ServiceCatalog{
		Name:             "nat",
		Version:          "v2.0",
		WithOutProjectID: true,
	},
	"nat_gatewayv2": ServiceCatalog{
		Name:    "nat",
		Version: "v2",
	},
	"elb": ServiceCatalog{
		Name:             "elb",
		Version:          "v1.0",
		WithOutProjectID: true,
	},
	"elbv2": ServiceCatalog{
		Name:             "elb",
		Version:          "v2.0",
		WithOutProjectID: true,
	},
	"fwv2": ServiceCatalog{
		Name:             "vpc",
		Version:          "v2.0",
		WithOutProjectID: true,
	},
	"vpcep": ServiceCatalog{
		Name:    "vpcep",
		Version: "v1",
	},
	"dns": ServiceCatalog{
		Name:             "dns",
		Version:          "v2",
		Scope:            "global",
		WithOutProjectID: true,
	},
	"dns_region": ServiceCatalog{
		Name:             "dns",
		Version:          "v2",
		WithOutProjectID: true,
	},

	// catalog for database
	"rdsv1": ServiceCatalog{
		Name:    "rds",
		Version: "rds/v1",
	},
	"rds": ServiceCatalog{
		Name:    "rds",
		Version: "v3",
	},
	"dds": ServiceCatalog{
		Name:    "dds",
		Version: "v3",
	},
	"cassandra": ServiceCatalog{
		Name:    "gaussdb-nosql",
		Version: "v3",
	},
	"gaussdb": ServiceCatalog{
		Name:    "gaussdb",
		Version: "mysql/v3",
	},
	"opengauss": ServiceCatalog{
		Name:    "gaussdb",
		Version: "opengauss/v3",
	},

	// catalog for management service
	"ces": ServiceCatalog{
		Name:    "ces",
		Version: "V1.0",
	},
	"cts": ServiceCatalog{
		Name:    "cts",
		Version: "v1.0",
	},
	"lts": ServiceCatalog{
		Name:    "lts",
		Version: "v2",
	},
	"smn": ServiceCatalog{
		Name:         "smn",
		Version:      "v2",
		ResourceBase: "notifications",
	},

	// catalog for Security service
	"anti-ddos": ServiceCatalog{
		Name:    "antiddos",
		Version: "v1",
	},
	"kms": ServiceCatalog{
		Name:             "kms",
		Version:          "v1.0",
		WithOutProjectID: true,
	},

	// catalog for Enterprise Intelligence
	"mrs": ServiceCatalog{
		Name:    "mrs",
		Version: "v1.1",
	},
	"dws": ServiceCatalog{
		Name:    "dws",
		Version: "v1.0",
	},
	"dli": ServiceCatalog{
		Name:    "dli",
		Version: "v1.0",
	},
	"disv2": ServiceCatalog{
		Name:    "dis",
		Version: "v2",
	},
	"css": ServiceCatalog{
		Name:    "css",
		Version: "v1.0",
	},
	"cs": ServiceCatalog{
		Name:    "cs",
		Version: "v1.0",
	},
	"ges": ServiceCatalog{
		Name:    "ges",
		Version: "v1.0",
	},
	"cloudtable": ServiceCatalog{
		Name:    "cloudtable",
		Version: "v2",
	},
	"cdm": ServiceCatalog{
		Name:    "cdm",
		Version: "v1.1",
	},

	// catalog for Application
	"apig": ServiceCatalog{
		Name:             "apig",
		Version:          "v1.0",
		ResourceBase:     "apigw",
		WithOutProjectID: true,
	},
	"dcsv1": ServiceCatalog{
		Name:    "dcs",
		Version: "v1.0",
	},
	"dcsv2": ServiceCatalog{
		Name:    "dcs",
		Version: "v2",
	},
	"dms": ServiceCatalog{
		Name:    "dms",
		Version: "v1.0",
	},
	"dmsv2": ServiceCatalog{
		Name:    "dms",
		Version: "v2",
	},

	// catalog for edge / IoT
	"iec": ServiceCatalog{
		Name:             "iecs",
		Version:          "v1",
		Scope:            "global",
		WithOutProjectID: true,
	},

	// catalog for Others
	"rts": ServiceCatalog{
		Name:    "rts",
		Version: "v1",
	},
	"oms": ServiceCatalog{
		Name:    "oms",
		Version: "v1",
	},
	"mls": ServiceCatalog{
		Name:    "mls",
		Version: "v1.0",
	},
}
