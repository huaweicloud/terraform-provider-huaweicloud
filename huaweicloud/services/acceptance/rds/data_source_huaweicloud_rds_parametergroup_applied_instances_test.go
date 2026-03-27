package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsParametergroupAppliedInstances_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_parametergroup_applied_instances.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsParametergroupAppliedInstances_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "entities.#"),
					resource.TestCheckResourceAttrSet(dataSource, "entities.0.entity_id"),
					resource.TestCheckResourceAttrSet(dataSource, "entities.0.entity_name"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_count_limit"),
				),
			},
		},
	})
}

func testDataSourceRdsParametergroupAppliedInstances_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rds_parametergroup_applied_instances" "test" {
  depends_on = [huaweicloud_rds_parametergroup_apply.test]

  config_id = huaweicloud_rds_parametergroup.test.id
}
`, testAccConfigurationApply_apply(name))
}
