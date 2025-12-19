package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceImageSensitiveInformationDetail_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_image_sensitive_information_detail.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHSSImageId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceImageSensitiveInformationDetail_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.sensitive_info_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.severity"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.position"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.file_path"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.content"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.latest_scan_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.handle_status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.operate_accept"),

					resource.TestCheckOutput("is_severity_filter_useful", "true"),
					resource.TestCheckOutput("is_handle_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceImageSensitiveInformationDetail_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_hss_image_sensitive_information_detail" "test" {
  image_id   = "%[1]s"
  image_type = "private_image"
}

# Filter using severity.
locals {
  severity = data.huaweicloud_hss_image_sensitive_information_detail.test.data_list[0].severity
}

data "huaweicloud_hss_image_sensitive_information_detail" "severity_filter" {
  image_id   = "%[1]s"
  image_type = "private_image"
  severity   = local.severity
}

output "is_severity_filter_useful" {
  value = length(data.huaweicloud_hss_image_sensitive_information_detail.severity_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_image_sensitive_information_detail.severity_filter.data_list[*].severity : v == local.severity]
  )
}

# Filter using handle_status.
locals {
  handle_status = data.huaweicloud_hss_image_sensitive_information_detail.test.data_list[0].handle_status
}

data "huaweicloud_hss_image_sensitive_information_detail" "handle_status_filter" {
  image_id      = "%[1]s"
  image_type    = "private_image"
  handle_status = local.handle_status
}

output "is_handle_status_filter_useful" {
  value = length(data.huaweicloud_hss_image_sensitive_information_detail.handle_status_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_image_sensitive_information_detail.handle_status_filter.data_list[*].handle_status : v == local.handle_status]
  )
}
`, acceptance.HW_HSS_IMAGE_ID)
}
