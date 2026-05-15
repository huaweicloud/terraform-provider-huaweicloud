package geminidb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGeminiDBRecyclingInstances_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_geminidb_recycling_instances.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGeminiDBRecyclingInstances_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.mode"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.data_store.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.data_store.0.version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.charge_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.backup_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.deleted_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.retained_until"),
				),
			},
		},
	})
}

func testAccDataSourceGeminiDBRecyclingInstances_basic() string {
	return `
data "huaweicloud_geminidb_recycling_instances" "test" {}
`
}
