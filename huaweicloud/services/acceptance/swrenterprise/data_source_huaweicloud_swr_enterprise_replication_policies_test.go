package swrenterprise

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrEnterpriseReplicationPolicies_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_enterprise_replication_policies.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSwrSignatureWithImageEnabled(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrEnterpriseReplicationPolicies_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "policies.#"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.repo_scope_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.override"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.enabled"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.src_registry.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.dest_registry.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.filters.0.value"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.filters.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.trigger.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.updated_at"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("registry_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSwrEnterpriseReplicationPolicies_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_swr_enterprise_replication_policies" "test" {
  depends_on = [huaweicloud_swr_enterprise_replication_policy.test]
  
  instance_id = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
}

data "huaweicloud_swr_enterprise_replication_policies" "filter_by_name" {
  instance_id = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  name        = huaweicloud_swr_enterprise_replication_policy.test.name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_swr_enterprise_replication_policies.filter_by_name.policies) > 0 && alltrue(
	[for v in data.huaweicloud_swr_enterprise_replication_policies.filter_by_name.policies[*].name : 
	  v == huaweicloud_swr_enterprise_replication_policy.test.name]
  )
}

data "huaweicloud_swr_enterprise_replication_policies" "filter_by_registry_id" {
  depends_on = [huaweicloud_swr_enterprise_replication_policy.test]

  instance_id = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  registry_id = huaweicloud_swr_enterprise_instance_registry.test.registry_id
}

output "registry_id_filter_is_useful" {
  value = length(data.huaweicloud_swr_enterprise_replication_policies.filter_by_registry_id.policies) > 0 && alltrue(
	[for v in data.huaweicloud_swr_enterprise_replication_policies.filter_by_registry_id.policies[*].dest_registry[0].id : 
	  v == huaweicloud_swr_enterprise_instance_registry.test.registry_id]
  )
}
`, testAccSwrEnterpriseReplicationPolicy_basic(name))
}
