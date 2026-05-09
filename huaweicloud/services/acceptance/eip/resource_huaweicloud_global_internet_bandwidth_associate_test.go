package eip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/eip"
)

func getGolbalEipFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region      = acceptance.HW_REGION_NAME
		globalEipId = state.Primary.ID
	)

	client, err := cfg.NewServiceClient("geip", region)
	if err != nil {
		return nil, fmt.Errorf("error creating GEIP client: %s", err)
	}

	return eip.GetGlobalEIPWithInternetBandwidth(client, cfg.DomainID, globalEipId)
}
func TestAccInternetBandwidthAssociate_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_global_internet_bandwidth_associate.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGolbalEipFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGlobalEipId(t)
			acceptance.TestAccPreCheckGlobalInternetBandwidthId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGlobalInternetBandwidthAssociateBasic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
		},
	})
}

func testGlobalInternetBandwidthAssociateBasic() string {
	return fmt.Sprintf(`
resource "huaweicloud_global_internet_bandwidth_associate" "test" {
  global_eip_id = "%s"

  global_eip {
    internet_bandwidth_id = "%s"
  }
}
`, acceptance.HW_GLOBAL_EIP_ID, acceptance.HW_GLOBAL_INTERNET_BANDWIDTH_ID)
}
