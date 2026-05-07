package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dds"
)

func getResourceIpAddressFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dds", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DDS client: %s", err)
	}

	return dds.GetInstanceInfo(client, state.Primary.ID)
}

func TestAccResourceIpAddress_basic(t *testing.T) {
	var (
		rName  = "huaweicloud_dds_ip_address.test"
		object interface{}
		rc     = acceptance.InitResourceCheck(
			rName,
			&object,
			getResourceIpAddressFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDDSInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccIpAddress_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "instance_id", acceptance.HW_DDS_INSTANCE_ID),
					resource.TestCheckResourceAttr(rName, "type", "config"),
					resource.TestCheckResourceAttr(rName, "password", "Test@1234"),
				),
			},
		},
	})
}

func testAccIpAddress_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dds_ip_address" "test" {
  instance_id = "%s"
  type        = "config"
  password    = "Test@1234"
}
`, acceptance.HW_DDS_INSTANCE_ID)
}
