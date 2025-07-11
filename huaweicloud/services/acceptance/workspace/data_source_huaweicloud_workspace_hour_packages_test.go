package workspace

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccHourPackagesDataSource_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_workspace_hour_packages.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccHourPackagesDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					// resource.TestCheckResourceAttr(dataSourceName, "hour_packages.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "hour_packages.0.cloud_service_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "hour_packages.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "hour_packages.0.resource_spec_code"),
					resource.TestCheckResourceAttrSet(dataSourceName, "hour_packages.0.desktop_resource_spec_code"),
					resource.TestCheckResourceAttrSet(dataSourceName, "hour_packages.0.descriptions.0.zh_cn"),
					resource.TestCheckResourceAttrSet(dataSourceName, "hour_packages.0.descriptions.0.en_us"),
					resource.TestCheckResourceAttrSet(dataSourceName, "hour_packages.0.package_duration"),
					resource.TestCheckResourceAttrSet(dataSourceName, "hour_packages.0.status"),
				),
			},
		},
	})
}

func testAccHourPackagesDataSource_basic() string {
	return `data "huaweicloud_workspace_hour_packages" "test" {}`
}
