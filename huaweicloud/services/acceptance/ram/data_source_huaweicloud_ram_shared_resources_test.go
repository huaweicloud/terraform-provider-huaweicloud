package ram

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceRAMSharedResources_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_ram_shared_resources.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceRAMSharedReources_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "shared_resources.0.resource_urn"),
					resource.TestCheckResourceAttrSet(rName, "shared_resources.0.resource_type"),
					resource.TestCheckResourceAttrSet(rName, "shared_resources.0.resource_share_id"),
					resource.TestCheckResourceAttrSet(rName, "shared_resources.0.status"),
					resource.TestCheckResourceAttrSet(rName, "shared_resources.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "shared_resources.0.updated_at"),

					resource.TestCheckOutput("resource_type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceRAMSharedReources_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_ram_shared_resources" "test" {
  resource_owner = "self"
}

data "huaweicloud_ram_shared_resources" "resource_type_filter" {
  resource_owner = "self"
  resource_type  = data.huaweicloud_ram_shared_resources.test.shared_resources.0.resource_type
}

locals {
  resource_type = data.huaweicloud_ram_shared_resources.test.shared_resources.0.resource_type
}

output "resource_type_filter_is_useful" {
  value = length(data.huaweicloud_ram_shared_resources.test.shared_resources) > 0 && alltrue(
    [for v in data.huaweicloud_ram_shared_resources.resource_type_filter.shared_resources[*].resource_type : v == local.resource_type]
  )
}
`, testRAMShare_basic(name))
}
