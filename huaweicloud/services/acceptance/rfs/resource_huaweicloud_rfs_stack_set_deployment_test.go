package rfs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccStackSetDeployment_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDomainId(t)
			acceptance.TestAccPreCheckRfsStackSetName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStackSetDeployment_basic(),
			},
		},
	})
}

func testAccStackSetDeployment_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_rfs_stack_set_deployment" "test" {
  stack_set_name = "%[1]s"
  
  deployment_targets {
    regions    = ["%[2]s"]
    domain_ids = ["%[3]s"]
  }

  template_body = <<-EOT
    resource "huaweicloud_vpc" "example" {
      name = "example-vpc"
      cidr = "192.168.0.0/16"
    }
  EOT
}
`, acceptance.HW_RFS_STACK_SET_NAME, acceptance.HW_REGION_NAME, acceptance.HW_DOMAIN_ID)
}
