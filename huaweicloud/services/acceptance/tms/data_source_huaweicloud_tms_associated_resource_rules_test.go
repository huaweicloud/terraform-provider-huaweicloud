package tms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTmsAssociatedResourceRules_basic(t *testing.T) {
	dataSource := "data.huaweicloud_tms_associated_resource_rules.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTmsAssociatedResourceRules_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "rules.#"),
				),
			},
		},
	})
}

func TestAccDataSourceTmsAssociatedResourceRules_filter(t *testing.T) {
	dataSource := "data.huaweicloud_tms_associated_resource_rules.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTmsAssociatedResourceRules_filter(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "rules.#"),
				),
			},
		},
	})
}

func testAccDataSourceTmsAssociatedResourceRules_basic() string {
	return `
data "huaweicloud_tms_associated_resource_rules" "test" {}
`
}

func testAccDataSourceTmsAssociatedResourceRules_filter() string {
	return `
data "huaweicloud_tms_associated_resource_rules" "test" {
  region_id = "cn-north-4"
}
`
}
