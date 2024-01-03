package swr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceImageTags_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_swr_image_tags.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSwrRepository(t)
			acceptance.TestAccPreCheckSwrOrigination(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceImageTags_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "image_tags.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "image_tags.0.digest"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("digest_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceImageTags_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_swr_image_tags" "test" {
  organization = "%[1]s"
  repository   = "%[2]s"
}

data "huaweicloud_swr_image_tags" "filter_by_name" {
  organization = "%[1]s"
  repository   = "%[2]s"
  name         = data.huaweicloud_swr_image_tags.test.image_tags[0].name
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_swr_image_tags.filter_by_name.image_tags) == 1
}

locals {
  digest = data.huaweicloud_swr_image_tags.test.image_tags[0].digest
}
data "huaweicloud_swr_image_tags" "filter_by_digest" {
  organization   = "%[1]s"
  repository     = "%[2]s"
  digest         = local.digest
}
output "digest_filter_is_useful" {
  value = length(data.huaweicloud_swr_image_tags.filter_by_digest.image_tags) > 0 && alltrue(
	[for v in data.huaweicloud_swr_image_tags.filter_by_digest.image_tags[*].digest : v == local.digest]
  )
}
`, acceptance.HW_SWR_ORGANIZATION, acceptance.HW_SWR_REPOSITORY)
}
