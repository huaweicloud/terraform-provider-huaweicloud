package ims

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceImsImageSharedMembers_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ims_image_shared_members.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSourceImage(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceImsImageSharedMembers_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "members.#"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.member_id"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.member_type"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.image_id"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.schema"),
					resource.TestCheckResourceAttrSet(dataSource, "schema"),
				),
			},
		},
	})
}

func testDataSourceImsImageSharedMembers_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_ims_image_shared_members" "test" {
  image_id = "%[1]s"
}
`, acceptance.HW_IMAGE_SHARE_SOURCE_IMAGE_ID)
}
