package swrenterprise

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrEnterpriseImmutableTagRules_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_enterprise_immutable_tag_rules.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrEnterpriseImmutableTagRules_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "rules.#"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.namespace_id"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.namespace_name"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.priority"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.disabled"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.action"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.template"),

					resource.TestCheckOutput("namespace_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSwrEnterpriseImmutableTagRules_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_swr_enterprise_immutable_tag_rules" "test" {
  depends_on = [huaweicloud_swr_enterprise_immutable_tag_rule.test]

  instance_id = huaweicloud_swr_enterprise_instance.test.id
}


data "huaweicloud_swr_enterprise_immutable_tag_rules" "filter_by_namespace_id" {
  instance_id  = huaweicloud_swr_enterprise_instance.test.id
  namespace_id = huaweicloud_swr_enterprise_immutable_tag_rule.test.namespace_id
}

output "namespace_id_filter_is_useful" {
  value = length(data.huaweicloud_swr_enterprise_immutable_tag_rules.filter_by_namespace_id.rules) > 0 && alltrue(
	[for v in data.huaweicloud_swr_enterprise_immutable_tag_rules.filter_by_namespace_id.rules[*].namespace_id : 
	  v == huaweicloud_swr_enterprise_immutable_tag_rule.test.namespace_id]
  )
}
`, testAccSwrEnterpriseImmutableTagRule_basic(name))
}
