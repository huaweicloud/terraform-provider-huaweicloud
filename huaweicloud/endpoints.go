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
	// ******* client for Compute start *******
	"ecs": ServiceCatalog{
		Name:    "ecs",
		Version: "v1",
	},
	"computeV11": ServiceCatalog{
		Name:    "ecs",
		Version: "v1.1",
	},
	"computeV2": ServiceCatalog{
		Name:    "ecs",
		Version: "v2.1",
	},
	"autoscalingV1": ServiceCatalog{
		Name:    "as",
		Version: "autoscaling-api/v1",
	},
	"imageV2": ServiceCatalog{
		Name:             "ims",
		Version:          "v2",
		WithOutProjectID: true,
	},
	"cceV3": ServiceCatalog{
		Name:    "cce",
		Version: "api/v3/projects",
	},
	"cceAddonV3": ServiceCatalog{
		Name:             "cce",
		Version:          "api/v3",
		WithOutProjectID: true,
	},
	"cciV1": ServiceCatalog{
		Name:             "cci",
		Version:          "apis/networking.cci.io/v1beta1",
		WithOutProjectID: true,
	},
	"FgsV2": ServiceCatalog{
		Name:    "functiongraph",
		Version: "v2",
	},
	// ******* client for Compute end *******

	"network": ServiceCatalog{
		Name:             "vpc",
		Version:          "v1",
		WithOutProjectID: true,
	},
	"networkV2": ServiceCatalog{
		Name:             "vpc",
		Version:          "v2.0",
		WithOutProjectID: true,
	},
	"vpc": ServiceCatalog{
		Name:             "vpc",
		Version:          "v1",
		WithOutProjectID: true,
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
	"bss": ServiceCatalog{
		Name:             "bss",
		Version:          "v1.0",
		WithOutProjectID: true,
	},
	"fwV2": ServiceCatalog{
		Name:             "vpc",
		Version:          "v2.0",
		WithOutProjectID: true,
	},
}
