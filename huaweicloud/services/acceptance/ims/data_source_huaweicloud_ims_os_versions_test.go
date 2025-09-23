package ims

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOsVersions_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_ims_os_versions.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOsVersion_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "os_versions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "os_versions.0.platform"),
					resource.TestCheckResourceAttrSet(dataSource, "os_versions.0.versions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "os_versions.0.versions.0.platform"),
					resource.TestCheckResourceAttrSet(dataSource, "os_versions.0.versions.0.os_version_key"),
					resource.TestCheckResourceAttrSet(dataSource, "os_versions.0.versions.0.os_version"),
					resource.TestCheckResourceAttrSet(dataSource, "os_versions.0.versions.0.os_bit"),
					resource.TestCheckResourceAttrSet(dataSource, "os_versions.0.versions.0.os_type"),
					resource.TestCheckOutput("is_test1_filter_successful", "true"),
					resource.TestCheckOutput("is_test2_filter_successful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceOsVersion_basic = `
data "huaweicloud_ims_os_versions" "test" {
}

data "huaweicloud_ims_os_versions" "test1" {
  tag = "x86"
}

output "is_test1_filter_successful" {
  value = length(data.huaweicloud_ims_os_versions.test1.os_versions) < length(data.huaweicloud_ims_os_versions.test.os_versions) 
}

data "huaweicloud_ims_os_versions" "test2" {
  tag = "x86,bms"
}

output "is_test2_filter_successful" {
  value = length(data.huaweicloud_ims_os_versions.test2.os_versions) < length(data.huaweicloud_ims_os_versions.test1.os_versions) 
}
`
