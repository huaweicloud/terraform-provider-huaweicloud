package configurations

import "github.com/huaweicloud/golangsdk"

type Configuration struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	DataStoreVer  string `json:"datastore_version_name"`
	DataStoreName string `json:"datastore_name"`
}

type ListResult struct {
	golangsdk.Result
}

func (lr ListResult) Extract() ([]Configuration, error) {
	var a struct {
		Configurations []Configuration `json:"configurations"`
	}
	err := lr.Result.ExtractInto(&a)
	return a.Configurations, err
}
