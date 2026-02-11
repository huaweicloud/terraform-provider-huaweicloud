package oms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceObjectstorageBucketRegion_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_oms_bucket_region.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAKAndSK(t)
			acceptance.TestAccPreCheckOmsObsBucketName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceObjectstorageBucketRegion_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "name"),
				),
			},
		},
	})
}

func testAccDataSourceObjectstorageBucketRegion_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_oms_bucket_region" "test" {
  cloud_type  = "HuaweiCloud"
  bucket_name = "%[1]s"
  ak          = "%[2]s"
  sk          = "%[3]s"
}
`, acceptance.HW_OMS_OBS_BUCKET_NAME, acceptance.HW_ACCESS_KEY, acceptance.HW_SECRET_KEY)
}
