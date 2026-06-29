package secmaster

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/secmaster"
)

func getResourceSiemShipperFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("secmaster", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	return secmaster.GetSiemShipperInfo(client, state.Primary.Attributes["workspace_id"], state.Primary.ID)
}

func TestAccResourceSiemShipper_basic(t *testing.T) {
	var (
		rName       = "huaweicloud_secmaster_siem_shipper.test"
		shipperName = acceptance.RandomAccResourceName()

		object interface{}
		rc     = acceptance.InitResourceCheck(
			rName,
			&object,
			getResourceSiemShipperFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceSiemShipper_basic(shipperName),
				ExpectError: regexp.MustCompile(`空间ID不存在`),
			},
		},
	})
}

func testAccResourceSiemShipper_basic(shipperName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_siem_shipper" "test" {
  consumption_type = 0
  dataspace_id     = "xxx"
  dataspace_name   = "xxx"
  domain_id        = "xxx"
  pipe_id          = "xxx"
  pipe_name        = "xxx"
  project_id       = "xxx"
  region           = "cn-north-4"
  shipper_name     = "%[1]s"
  version          = "v1"
  workspace_id     = "%[2]s"
  workspace_name   = "xxx"

  shipper_destination {
    data_param                 = jsonencode({})
    destination_dataspace      = "xxx"
    destination_dataspace_name = "xxx"
    destination_identity_role  = "pipe-strategy"
    destination_pipe           = "xxx"
    destination_pipe_name      = "xxx"
    destination_region         = "cn-north-4"
    destination_shipper_type   = 0
    destination_workspace      = "xxx"
    destination_workspace_name = "xxx"
  }

  shipper_source {
    region                = "cn-north-4"
    source_dataspace      = "xxx"
    source_dataspace_name = "xxx"
    source_identity_role  = "pipe-strategy"
    source_pipe           = "xxx"
    source_pipe_name      = "xxx"
    source_type           = 0
    source_workspace      = "xxx"
    source_workspace_name = "xxx"
  }
}
`, shipperName, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
