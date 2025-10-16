package swrenterprise

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrEnterpriseRetentionPolicies_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_enterprise_retention_policies.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrEnterpriseRetentionPolicies_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "policies.#"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.enabled"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.namespace_id"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.namespace_name"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.algorithm"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("namespace_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSwrEnterpriseRetentionPolicies_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_swr_enterprise_retention_policies" "test" {
  depends_on = [huaweicloud_swr_enterprise_retention_policy.test]

  instance_id = huaweicloud_swr_enterprise_instance.test.id
}

data "huaweicloud_swr_enterprise_retention_policies" "filter_by_name" {
  instance_id = huaweicloud_swr_enterprise_instance.test.id
  name        = huaweicloud_swr_enterprise_retention_policy.test.name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_swr_enterprise_retention_policies.filter_by_name.policies) > 0 && alltrue(
	[for v in data.huaweicloud_swr_enterprise_retention_policies.filter_by_name.policies[*].name : 
	  v == huaweicloud_swr_enterprise_retention_policy.test.name]
  )
}

data "huaweicloud_swr_enterprise_retention_policies" "filter_by_namespace_id" {
  instance_id  = huaweicloud_swr_enterprise_instance.test.id
  namespace_id = huaweicloud_swr_enterprise_retention_policy.test.namespace_id
}

output "namespace_id_filter_is_useful" {
  value = length(data.huaweicloud_swr_enterprise_retention_policies.filter_by_namespace_id.policies) > 0 && alltrue(
	[for v in data.huaweicloud_swr_enterprise_retention_policies.filter_by_namespace_id.policies[*].namespace_id : 
	  v == huaweicloud_swr_enterprise_retention_policy.test.namespace_id]
  )
}
`, testAccSwrEnterpriseTrigger_basic(name))
}
