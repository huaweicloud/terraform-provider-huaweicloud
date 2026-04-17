package swr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrSignedImageAssociatedImageTags_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_signed_image_associated_image_tags.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSWRSignedImageAssicatedTags(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSwrSignedImageAssociatedImageTags_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tags.#"),
				),
			},
		},
	})
}

func testAccDataSourceSwrSignedImageAssociatedImageTags_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_swr_signed_image_associated_image_tags" "test" {
  namespace  = "%s"
  repository = "%s"
  sig_tag    = "%s"
}
`, acceptance.HW_SWR_ORGANIZATION, acceptance.HW_SWR_REPOSITORY, acceptance.HW_SWR_SIG_TAG)
}
