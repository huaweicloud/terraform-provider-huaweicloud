package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceClusterProtectProtectionItems_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_cluster_protect_protection_items.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceClusterProtectProtectionItems_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vuls.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "baselines.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "baselines.0.baseline_desc"),
					resource.TestCheckResourceAttrSet(dataSourceName, "baselines.0.baseline_index"),
					resource.TestCheckResourceAttrSet(dataSourceName, "baselines.0.baseline_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "malwares.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "malwares.0.malware_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.image_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.image_version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "clusters.#"),
				),
			},
		},
	})
}

func testAccDataSourceClusterProtectProtectionItems_basic() string {
	return `
data "huaweicloud_hss_cluster_protect_protection_items" "test" {
  enterprise_project_id = "0"
}
`
}
