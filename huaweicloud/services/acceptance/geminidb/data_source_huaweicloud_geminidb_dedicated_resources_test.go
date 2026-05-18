package geminidb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDedicatedResources_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_geminidb_dedicated_resources.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDedicatedResources_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "resources.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.engine_name"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.availability_zone"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.architecture"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.capacity.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.capacity.0.vcpus"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.capacity.0.ram"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.capacity.0.volume"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.status"),
				),
			},
		},
	})
}

const testAccDataSourceDedicatedResources_basic = `data "huaweicloud_geminidb_dedicated_resources" "test" {}`
