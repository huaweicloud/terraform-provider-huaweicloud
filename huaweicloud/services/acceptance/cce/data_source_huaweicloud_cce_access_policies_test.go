package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAccessPolicies_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cce_access_policies.test"
	name := acceptance.RandomAccResourceNameWithDash()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckUserId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAccessPolicies_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "access_policy_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "access_policy_list.0.kind"),
					resource.TestCheckResourceAttrSet(dataSource, "access_policy_list.0.api_version"),
					resource.TestCheckResourceAttrSet(dataSource, "access_policy_list.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "access_policy_list.0.policy_id"),
					resource.TestCheckResourceAttrSet(dataSource, "access_policy_list.0.clusters.#"),
					resource.TestCheckResourceAttrSet(dataSource, "access_policy_list.0.access_scope.#"),
					resource.TestCheckResourceAttrSet(dataSource, "access_policy_list.0.access_scope.0.namespaces.#"),
					resource.TestCheckResourceAttrSet(dataSource, "access_policy_list.0.policy_type"),
					resource.TestCheckResourceAttrSet(dataSource, "access_policy_list.0.principal.#"),
					resource.TestCheckResourceAttrSet(dataSource, "access_policy_list.0.principal.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "access_policy_list.0.principal.0.ids.#"),
					resource.TestCheckResourceAttrSet(dataSource, "access_policy_list.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "access_policy_list.0.update_time"),
					resource.TestCheckOutput("cluster_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceAccessPolicies_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_access_policy" "test" {
  name     = "%[2]s"
  clusters = [huaweicloud_cce_cluster.test.id]

  access_scope {
    namespaces = ["default"]
  }

  policy_type = "CCEClusterAdminPolicy"

  principal {
    type = "user"
    ids  = ["%[3]s"]
  }
}
`, testAccCluster_basic(name), name, acceptance.HW_USER_ID)
}

func testAccDataSourceAccessPolicies_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cce_access_policies" "test" {
  depends_on = [huaweicloud_cce_access_policy.test]
}

data "huaweicloud_cce_access_policies" "cluster_id_filter" {
  depends_on = [huaweicloud_cce_access_policy.test]

  cluster_id = huaweicloud_cce_cluster.test.id
}
locals{
  cluster_id = huaweicloud_cce_cluster.test.id
}
output "cluster_id_filter_is_useful" {
  value = length(data.huaweicloud_cce_access_policies.cluster_id_filter.access_policy_list) > 0 && alltrue(
    [for v in data.huaweicloud_cce_access_policies.cluster_id_filter.access_policy_list[*].clusters : alltrue(
      [for vv in v : vv == local.cluster_id]
    )]
  )
}
`, testAccDataSourceAccessPolicies_base(name))
}
