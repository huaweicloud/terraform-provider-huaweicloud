package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDDSSlowLogLinks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dds_slow_log_links.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDDSInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDDSSlowLogLinks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "links.#"),
					resource.TestCheckResourceAttrSet(dataSource, "links.0.file_name"),
					resource.TestCheckResourceAttrSet(dataSource, "links.0.node_name"),
					resource.TestCheckResourceAttrSet(dataSource, "links.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "links.0.file_size"),
					resource.TestCheckResourceAttrSet(dataSource, "links.0.file_link"),
					resource.TestCheckResourceAttrSet(dataSource, "links.0.updated_at"),

					resource.TestCheckOutput("file_name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDDSSlowLogLinks_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dds_slow_log_links" "test" {
  instance_id = "%[1]s"
}

locals {
  file_name = data.huaweicloud_dds_slow_log_links.test.links.0.file_name
}

data "huaweicloud_dds_slow_log_links" "filter" {
  instance_id    = "%[1]s"
  file_name_list = [local.file_name]
}

output "file_name_filter_is_useful" {
  value = length(data.huaweicloud_dds_slow_log_links.filter.links) > 0 && alltrue(
    [for v in data.huaweicloud_dds_slow_log_links.filter.links[*].file_name : v == local.file_name]
  )
}`, acceptance.HW_DDS_INSTANCE_ID)
}
