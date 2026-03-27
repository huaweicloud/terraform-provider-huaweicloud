package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsParametergroupApplyHistories_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_parametergroup_apply_histories.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsParametergroupApplyHistories_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "histories.#"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.instance_name"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.apply_result"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.apply_time"),
				),
			},
		},
	})
}

func testDataSourceRdsParametergroupApplyHistories_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rds_parametergroup_apply_histories" "test" {
  depends_on = [huaweicloud_rds_parametergroup_apply.test]

  config_id = huaweicloud_rds_parametergroup.test.id
}
`, testAccConfigurationApply_apply(name))
}
