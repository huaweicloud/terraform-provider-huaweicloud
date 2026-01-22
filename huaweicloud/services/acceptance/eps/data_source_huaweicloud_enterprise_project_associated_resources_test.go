package eps

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataEnterpriseProjectAssociatedResources_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_enterprise_project_associated_resources.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckECSID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataEnterpriseProjectAssociatedResources_basic_test(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "associated_resources.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "associated_resources.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "associated_resources.0.resource_type"),
				),
			},
		},
	})
}

func testAccDataEnterpriseProjectAssociatedResources_basic_test() string {
	return fmt.Sprintf(`
data "huaweicloud_enterprise_project_associated_resources" "test" {
  resource_id   = "%s"
  resource_type = "ecs"
}
`, acceptance.HW_ECS_ID)
}
