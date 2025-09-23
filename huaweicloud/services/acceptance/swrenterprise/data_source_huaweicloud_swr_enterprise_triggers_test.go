package swrenterprise

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrEnterpriseTriggers_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_enterprise_triggers.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrEnterpriseTriggers_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "triggers.#"),
					resource.TestCheckResourceAttrSet(dataSource, "triggers.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "triggers.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "triggers.0.enabled"),
					resource.TestCheckResourceAttrSet(dataSource, "triggers.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "triggers.0.event_types.#"),
					resource.TestCheckResourceAttrSet(dataSource, "triggers.0.namespace_id"),
					resource.TestCheckResourceAttrSet(dataSource, "triggers.0.namespace_name"),
					resource.TestCheckResourceAttrSet(dataSource, "triggers.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "triggers.0.updated_at"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("namespace_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSwrEnterpriseTriggers_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_swr_enterprise_triggers" "test" {
  depends_on = [huaweicloud_swr_enterprise_trigger.test]

  instance_id = huaweicloud_swr_enterprise_instance.test.id
}

data "huaweicloud_swr_enterprise_triggers" "filter_by_name" {
  instance_id = huaweicloud_swr_enterprise_instance.test.id
  name        = huaweicloud_swr_enterprise_trigger.test.name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_swr_enterprise_triggers.filter_by_name.triggers) > 0 && alltrue(
	[for v in data.huaweicloud_swr_enterprise_triggers.filter_by_name.triggers[*].name : v == huaweicloud_swr_enterprise_trigger.test.name]
  )
}

data "huaweicloud_swr_enterprise_triggers" "filter_by_namespace_id" {
  instance_id  = huaweicloud_swr_enterprise_instance.test.id
  namespace_id = huaweicloud_swr_enterprise_trigger.test.namespace_id
}

output "namespace_id_filter_is_useful" {
  value = length(data.huaweicloud_swr_enterprise_triggers.filter_by_namespace_id.triggers) > 0 && alltrue(
	[for v in data.huaweicloud_swr_enterprise_triggers.filter_by_namespace_id.triggers[*].namespace_id : 
	  v == huaweicloud_swr_enterprise_trigger.test.namespace_id]
  )
}
`, testAccSwrEnterpriseTrigger_basic(name))
}
