package rfs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceStackResources_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_rfs_stack_resources.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRfsStackName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceStackResources_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "stack_resources.#"),
					resource.TestCheckResourceAttrSet(dataSource, "stack_resources.0.physical_resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "stack_resources.0.physical_resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "stack_resources.0.logical_resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "stack_resources.0.logical_resource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "stack_resources.0.resource_status"),
					resource.TestCheckResourceAttrSet(dataSource, "stack_resources.0.resource_attributes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "stack_resources.0.resource_attributes.0.key"),
					resource.TestCheckResourceAttrSet(dataSource, "stack_resources.0.resource_attributes.0.value"),
				),
			},
		},
	})
}

func testDataSourceStackResources_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rfs_stack_resources" "test" {
  stack_name = "%[1]s"
}
`, acceptance.HW_RFS_STACK_NAME)
}
