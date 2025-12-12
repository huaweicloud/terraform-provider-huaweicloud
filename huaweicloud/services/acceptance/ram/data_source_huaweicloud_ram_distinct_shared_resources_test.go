package ram

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceRAMDistinctSharedResources_basic(t *testing.T) {
	var (
		rName = "data.huaweicloud_ram_distinct_shared_resources.test"
		dc    = acceptance.InitDataSourceCheck(rName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceRAMDistinctSharedReources_base(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "distinct_shared_resources.0.resource_urn"),
					resource.TestCheckResourceAttrSet(rName, "distinct_shared_resources.0.resource_type"),
					resource.TestCheckResourceAttrSet(rName, "distinct_shared_resources.0.updated_at"),
				),
			},
		},
	})
}

func testAccDatasourceRAMDistinctSharedReources_base() string {
	return `
data "huaweicloud_ram_distinct_shared_resources" "test" {
	resource_owner = "self"
}
`
}
