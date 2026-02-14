package oms

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running test, prepare a OBS bucket and create a folder named test in this bucket.
func TestAccDataSourceBucketObjects_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_oms_bucket_objects.test"
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
				Config: testAccDataSourceBucketObjects_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "records.#", regexp.MustCompile(`^[1-9][0-9]*$`)),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.size"),
				),
			},
		},
	})
}

func testAccDataSourceBucketObjects_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_oms_bucket_objects" "test" {
  cloud_type  = "HuaweiCloud"
  ak          = "%[1]s"
  sk          = "%[2]s"
  bucket_name = "%[3]s"
  file_path   = "test/"
}
`, acceptance.HW_ACCESS_KEY, acceptance.HW_SECRET_KEY, acceptance.HW_OMS_OBS_BUCKET_NAME)
}
