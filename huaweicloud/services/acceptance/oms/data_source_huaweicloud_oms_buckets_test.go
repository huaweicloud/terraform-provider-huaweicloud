package oms

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceObjectstorageBuckets_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_oms_buckets.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAKAndSK(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceObjectstorageBuckets_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "buckets.#", regexp.MustCompile(`^[1-9][0-9]*$`)),
				),
			},
		},
	})
}

func testAccDataSourceObjectstorageBuckets_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_oms_buckets" "test" {
  cloud_type = "HuaweiCloud"
  ak         = "%[1]s"
  sk         = "%[2]s"
}
`, acceptance.HW_ACCESS_KEY, acceptance.HW_SECRET_KEY)
}
