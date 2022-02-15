package sweep_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/rms/v1/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

// Mapping relationship between terraform resources name and rms. value: {provider}_{type}:{resourceName}
var defaultResourceMapping = map[string]string{
	"vpc_vpcs": "huaweicloud_vpc",
}

// All resources under the current user
var SweepResourcesList map[string][]resources.Resource

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func init() {
	SweepResourcesList = make(map[string][]resources.Resource)

	err := initProviderConfig()
	if err != nil {
		logp.Printf("[ERROR]Error init huaweicloud provider config: %s", err)
	}

	region := acceptance.HW_REGION_NAME
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	rmsClient, err := config.RmsV1Client(region)
	if err != nil {
		logp.Printf("Error creating RMS client: %s", err)
		return
	}

	pages, err := resources.List(rmsClient, resources.ListOpts{Region: region}).AllPages()

	if err != nil {
		logp.Printf("Unable to retrieve all resources in the region %s:%s ", region, err)
		return
	}

	allResources, err := resources.ExtractResources(pages)
	if err != nil {
		logp.Printf("Unable to retrieve all resources in the region %s:%s ", region, err)
		return
	}

	for _, v := range allResources {
		key := fmt.Sprintf("%s_%s", v.Provider, v.Type)
		if rName, ok := defaultResourceMapping[key]; ok {
			b := SweepResourcesList[rName]
			SweepResourcesList[rName] = append(b, v)
		} else {
			b := SweepResourcesList["unKnow_resources"]
			SweepResourcesList["unKnow_resources"] = append(b, v)
		}
	}
	logp.Printf("[WARN]%d resource is not in sweeper, please check: %s ", len(SweepResourcesList),
		SweepResourcesList["unKnow_resources"])
}

func initProviderConfig() error {
	testProvider := acceptance.TestAccProvider
	raw := make(map[string]interface{})
	diags := testProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		return fmtp.Errorf("Unexpected error when configure HuaweiCloud provider: %s", diags[0].Summary)
	}
	return nil
}

func getResources(tfResourceName string) []resources.Resource {
	return SweepResourcesList[tfResourceName]
}
