package cdn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataStatistics_basic(t *testing.T) {
	rName := "data.huaweicloud_cdn_domain_statistics.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCdnDomainName(t)
			acceptance.TestAccPrecheckCDNAnalytics(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataStatistics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "result", "{}"),
				),
			},
		},
	})
}

func testAccDataStatistics_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cdn_domain_statistics" "test" {
  domain_name           = "%[1]s"
  stat_type             = "%[2]s"
  start_time            = "%[3]s"
  end_time              = "%[4]s"
  action                = "location_detail"
  interval              = 3600
  group_by              = "domain"
  country               = "cn"
  province              = "beijing"
  isp                   = "yidong"
  enterprise_project_id = "0"
}
`, acceptance.HW_CDN_DOMAIN_NAME, acceptance.HW_CDN_STAT_TYPE, acceptance.HW_CDN_START_TIME, acceptance.HW_CDN_END_TIME)
}
