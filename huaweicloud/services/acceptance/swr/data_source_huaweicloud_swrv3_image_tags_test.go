package swr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrv3ImageTags_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swrv3_image_tags.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSwrRepository(t)
			acceptance.TestAccPreCheckSwrOrigination(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrv3ImageTags_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tags.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.schema"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.size"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.tag"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.digest"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.path"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.internal_path"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.repo_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.image_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.created"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.updated"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.manifest"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.is_trusted"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.tag_type"),

					resource.TestCheckOutput("tag_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSwrv3ImageTags_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_swrv3_image_tags" "test" {
  organization  = "%[1]s"
  repository    = "%[2]s"
  with_manifest = true
}

data "huaweicloud_swrv3_image_tags" "filter_by_tag" {
  organization = "%[1]s"
  repository   = "%[2]s"
  tag          = data.huaweicloud_swrv3_image_tags.test.tags[0].tag
}

output "tag_filter_is_useful" {
  value = length(data.huaweicloud_swrv3_image_tags.filter_by_tag.tags) == 1
}
`, acceptance.HW_SWR_ORGANIZATION, acceptance.HW_SWR_REPOSITORY)
}
