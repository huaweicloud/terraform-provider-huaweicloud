package swr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrSignedImageAttachments_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_signed_image_attachments.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSWRSignedImageAttachments(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSwrSignedImageAttachments_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "accessories.#"),
					resource.TestCheckResourceAttrSet(dataSource, "accessories.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "accessories.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "accessories.0.namespace_name"),
					resource.TestCheckResourceAttrSet(dataSource, "accessories.0.repo_name"),
					resource.TestCheckResourceAttrSet(dataSource, "accessories.0.sig_tag"),
					resource.TestCheckResourceAttrSet(dataSource, "accessories.0.sig_digest"),
					resource.TestCheckResourceAttrSet(dataSource, "accessories.0.target_digest"),
					resource.TestCheckResourceAttrSet(dataSource, "accessories.0.size"),
					resource.TestCheckResourceAttrSet(dataSource, "accessories.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "accessories.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "accessories.0.updated_at"),
				),
			},
		},
	})
}

func testAccDataSourceSwrSignedImageAttachments_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_swr_signed_image_attachments" "test" {
  namespace  = "%s"
  repository = "%s"
  tag        = "%s"
}
`, acceptance.HW_SWR_ORGANIZATION, acceptance.HW_SWR_REPOSITORY, acceptance.HW_SWR_TAG)
}
