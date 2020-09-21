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
	"iam": ServiceCatalog{
		Name:             "iam",
		Version:          "v3",
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
	"dns": ServiceCatalog{
		Name:             "dns",
		Version:          "v2",
		Scope:            "global",
		WithOutProjectID: true,
	},

	"ecs": ServiceCatalog{
		Name:    "ecs",
		Version: "v1",
	},
	"compute": ServiceCatalog{
		Name:    "ecs",
		Version: "v2.1",
	},
	"network": ServiceCatalog{
		Name:    "vpc",
		Version: "v2.0",
	},
	"vpc": ServiceCatalog{
		Name:    "vpc",
		Version: "v1",
	},
	"volumev2": ServiceCatalog{
		Name:    "evs",
		Version: "v2",
	},
	"volumev3": ServiceCatalog{
		Name:    "evs",
		Version: "v3",
	},
	"sfs-turbo": ServiceCatalog{
		Name:    "sfs-turbo",
		Version: "v1",
	},
	"dcsv2": ServiceCatalog{
		Name:    "dcs",
		Version: "v2",
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
}
