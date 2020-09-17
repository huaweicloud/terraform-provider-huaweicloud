package huaweicloud

// ServiceCatalog defines a struct which was used to generate a service client for huaweicloud.
// the endpoint likes https://{Name}.{Region}.myhuaweicloud.com/{Version}/{project_id}/{ResourceBase}
// For more information, please refer to Config.NewServiceClient
type ServiceCatalog struct {
	Name             string
	Version          string
	Scope            string
	ResourceBase     string
	WithOutProjectID bool
}

var allServiceCatalog = map[string]ServiceCatalog{
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
}
