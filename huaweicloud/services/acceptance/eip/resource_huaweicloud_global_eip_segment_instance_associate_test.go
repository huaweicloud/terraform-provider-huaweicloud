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

// The test environment is missing
func getGolbalEipSegmentFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region      = acceptance.HW_REGION_NAME
		globalEipId = state.Primary.ID
	)

	client, err := cfg.NewServiceClient("geip", region)
	if err != nil {
		return nil, fmt.Errorf("error creating GEIP client: %s", err)
	}

	return eip.GetGlobalsegment(client, cfg.DomainID, globalEipId)
}
func TestAccSegmentInstanceAssociate_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_global_eip_segment_instance_associate.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGolbalEipSegmentFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGlobalEipSegmentId(t)
			acceptance.TestAccPreCheckGlobalEipSegmentInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGlobalSegmentAssociateBasic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
		},
	})
}

func testGlobalSegmentAssociateBasic() string {
	return fmt.Sprintf(`
resource "huaweicloud_global_eip_segment_instance_associate" "test" {
  global_eip_segment_id = "%[1]s"

  global_eip_segment {
    region        = "cn-south-1"
    instance_id   = "%[2]s"
	instance_type = "DC-CONNECT-GATEWAY"
    project_id    = "f0ab1d29eb70461fa7ef2cbd133f573e"
  }
}
`, acceptance.HW_GLOBAL_EIP_SEGMENT_ID, acceptance.HW_GLOBAL_EIP_SEGMENT_INSTANCE_ID)
}
