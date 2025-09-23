package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSecmasterPlaybookVersions_basic(t *testing.T) {
	dataSource := "data.huaweicloud_secmaster_playbook_versions.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSecmasterPlaybookVersions_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "playbook_versions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "playbook_versions.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "playbook_versions.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "playbook_versions.0.enabled"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					resource.TestCheckOutput("is_enabled_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSecmasterPlaybookVersions_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_secmaster_playbook_versions" "test" {
  workspace_id = "%[2]s"
  playbook_id  = huaweicloud_secmaster_playbook.test.id

  depends_on = [huaweicloud_secmaster_playbook_version.test]
}

locals {
  status  = data.huaweicloud_secmaster_playbook_versions.test.playbook_versions[0].status
  type    = data.huaweicloud_secmaster_playbook_versions.test.playbook_versions[0].type
  enabled = data.huaweicloud_secmaster_playbook_versions.test.playbook_versions[0].enabled
}

data "huaweicloud_secmaster_playbook_versions" "filter_by_status" {
  workspace_id = "%[2]s"
  playbook_id  = huaweicloud_secmaster_playbook.test.id
  status       = local.status
}

data "huaweicloud_secmaster_playbook_versions" "filter_by_type" {
  workspace_id = "%[2]s"
  playbook_id  = huaweicloud_secmaster_playbook.test.id
  type         = tostring(local.type)
}

data "huaweicloud_secmaster_playbook_versions" "filter_by_enabled" {
  workspace_id = "%[2]s"
  playbook_id  = huaweicloud_secmaster_playbook.test.id
  enabled      = local.enabled ? "1" : "0"
}

locals {
  list_by_status  = data.huaweicloud_secmaster_playbook_versions.filter_by_status.playbook_versions
  list_by_type    = data.huaweicloud_secmaster_playbook_versions.filter_by_type.playbook_versions
  list_by_enabled = data.huaweicloud_secmaster_playbook_versions.filter_by_enabled.playbook_versions
}

output "is_status_filter_useful" {
  value = length(local.list_by_status) > 0 && alltrue(
    [for v in local.list_by_status[*].status : v == local.status]
  )
}

output "is_type_filter_useful" {
  value = length(local.list_by_type) > 0 && alltrue(
    [for v in local.list_by_type[*].type : v == local.type]
  )
}

output "is_enabled_filter_useful" {
  value = length(local.list_by_enabled) > 0 && alltrue(
    [for v in local.list_by_enabled[*].enabled : v == local.enabled]
  )
}
`, testPlaybookVersion_basic(name), acceptance.HW_SECMASTER_WORKSPACE_ID)
}
