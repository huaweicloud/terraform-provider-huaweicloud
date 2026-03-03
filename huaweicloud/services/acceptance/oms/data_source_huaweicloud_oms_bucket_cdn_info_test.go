package oms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Due to the lack of testing conditions, this test case cannot be tested.
func TestAccDataSourceBucketCdnInfo_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_oms_bucket_cdn_info.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAKAndSK(t)
			acceptance.TestAccPreCheckOmsObsBucketName(t)
			acceptance.TestAccPreCheckOmsDomain(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBucketCdnInfo_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "checked_keys.#"),
				),
			},
		},
	})
}

func testAccDataSourceBucketCdnInfo_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_oms_bucket_cdn_info" "test" {
  cloud_type = "HuaweiCloud"
  ak         = "%[1]s"
  sk         = "%[2]s"
  bucket     = "%[3]s"

  source_cdn {
    authentication_type = "NONE"
    protocol            = "https"
    domain              = "%[4]s"
  }
}
`, acceptance.HW_ACCESS_KEY, acceptance.HW_SECRET_KEY, acceptance.HW_OMS_OBS_BUCKET_NAME, acceptance.HW_OMS_DOMAIN)
}
