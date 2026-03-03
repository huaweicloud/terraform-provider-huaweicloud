package oms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCloudVenderRegions_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_oms_cloud_vender_regions.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCloudVenderRegions_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "region_info.#"),
					resource.TestCheckResourceAttrSet(dataSource, "region_info.0.service_name"),
					resource.TestCheckResourceAttrSet(dataSource, "region_info.0.region_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "region_info.0.region_list.0.cloud_type"),
					resource.TestCheckResourceAttrSet(dataSource, "region_info.0.region_list.0.value"),
					resource.TestCheckResourceAttrSet(dataSource, "region_info.0.region_list.0.description"),
				),
			},
		},
	})
}

const testAccDataSourceCloudVenderRegions_basic = `data "huaweicloud_oms_cloud_vender_regions" "test" {}`
