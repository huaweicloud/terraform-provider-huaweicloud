package workspace

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccAppServerGroupRestrict_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_workspace_app_server_group_restrict.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAppServerGroupRestrict_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "max_session", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestMatchResourceAttr(dataSourceName, "max_group_count", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestCheckOutput("is_max_session_set", "true"),
					resource.TestCheckOutput("is_max_group_count_set", "true"),
				),
			},
		},
	})
}

const testAccAppServerGroupRestrict_basic = `
data "huaweicloud_workspace_app_server_group_restrict" "test" {}

locals {
  restrict = data.huaweicloud_workspace_app_server_group_restrict.test
}

output "is_max_session_set" {
  value = try(local.restrict.max_session >= 0, false)
}

output "is_max_group_count_set" {
  value = try(local.restrict.max_group_count >= 0, false)
}
`
