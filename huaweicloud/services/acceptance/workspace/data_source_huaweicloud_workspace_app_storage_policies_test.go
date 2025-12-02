package workspace

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataAppStoragePolicies_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_workspace_app_storage_policies.all"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataAppStoragePolicies_basic,
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameter.
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_system_policy_exist", "true"),
					resource.TestCheckOutput("is_policy_id_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataAppStoragePolicies_basic string = `
resource "huaweicloud_workspace_app_storage_policy" "test" {
  server_actions = ["GetObject"]
}

data "huaweicloud_workspace_app_storage_policies" "all" {
  depends_on = [huaweicloud_workspace_app_storage_policy.test]
}

# Filter all system policies
locals {
  system_policies = [for o in data.huaweicloud_workspace_app_storage_policies.all.policies: o if strcontains(o.id, "DEFAULT")]
}

output "is_system_policy_exist" {
  value = length(local.system_policies) > 0
}

# Filter by storage permission policy ID, in manual
locals {
  policy_id = huaweicloud_workspace_app_storage_policy.test.id

  policy_id_filter_result = [for o in data.huaweicloud_workspace_app_storage_policies.all.policies: o if o.id == local.policy_id]
}

output "is_policy_id_filter_useful" {
  value = length(local.policy_id_filter_result) == 1
}
`
