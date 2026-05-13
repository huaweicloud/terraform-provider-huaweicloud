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

func getSegmentBandwidthAssociateFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("geip", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating GEIP client: %s", err)
	}

	return eip.GetEipSegmentDetail(client, cfg.DomainID, state.Primary.ID)
}

func TestAccSegmentBandwidthAssociate_basic(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_global_eip_segment_bandwidth_associate.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getSegmentBandwidthAssociateFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGlobalEipSegmentId(t)
			acceptance.TestAccPreCheckGlobalEipBandwidthId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSegmentBandwidthAssociate_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "global_eip_segment_id", acceptance.HW_GLOBAL_EIP_SEGMENT_ID),
					resource.TestCheckResourceAttr(rName, "internet_bandwidth_id", acceptance.HW_GLOBAL_EIP_BANDWIDTH_ID),
					resource.TestCheckResourceAttrSet(rName, "internet_bandwidth.0.id"),
					resource.TestCheckResourceAttrSet(rName, "internet_bandwidth.0.size"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSegmentBandwidthAssociate_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_global_eip_segment_bandwidth_associate" "test" {
  global_eip_segment_id = "%[1]s"
  internet_bandwidth_id = "%[2]s"
}
`, acceptance.HW_GLOBAL_EIP_SEGMENT_ID, acceptance.HW_GLOBAL_EIP_BANDWIDTH_ID)
}
