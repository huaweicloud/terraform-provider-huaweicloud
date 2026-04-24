package rfs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccRfsContinueDeployStack_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRfsStackName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config:      testAccRfsContinueDeployStack_basic(),
				ExpectError: regexp.MustCompile(`stack deployment failed`),
			},
		},
	})
}

func testAccRfsContinueDeployStack_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rfs_stacks" "test" {}

locals {
  stack_name = "%s"
  stack_ids  = [for s in data.huaweicloud_rfs_stacks.test.stacks : s.stack_id if s.stack_name == local.stack_name]
  stack_id   = try(local.stack_ids[0], null)
}

resource "huaweicloud_rfs_continue_deploy_stack" "test" {
  stack_name = local.stack_name
  stack_id   = local.stack_id
}
`, acceptance.HW_RFS_STACK_NAME)
}
