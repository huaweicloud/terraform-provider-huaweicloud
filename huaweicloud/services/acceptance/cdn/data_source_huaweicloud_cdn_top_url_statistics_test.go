// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CDN
// ---------------------------------------------------------------

package cdn

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccCdnTopUrlStatistics_basic(t *testing.T) {
	var (
		// Get today's 0:00 UTC format string
		now       = time.Now().UTC()
		today     = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		startTime = today.Format("2006-01-02T15:04:05Z")

		// Get tomorrow's 0:00 UTC format string
		tomorrow = today.AddDate(0, 0, 1)
		endTime  = tomorrow.Format("2006-01-02T15:04:05Z")

		dcName = "data.huaweicloud_cdn_top_url_statistics.test"
		dc     = acceptance.InitDataSourceCheck(dcName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCDN(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCdnTopUrlStatistics_basic(startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "top_url_summary.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dcName, "top_url_summary.0.url"),
					resource.TestCheckResourceAttrSet(dcName, "top_url_summary.0.value"),
					resource.TestCheckResourceAttrSet(dcName, "top_url_summary.0.start_time"),
					resource.TestCheckResourceAttrSet(dcName, "top_url_summary.0.end_time"),
					resource.TestCheckResourceAttrSet(dcName, "top_url_summary.0.stat_type"),
				),
			},
		},
	})
}

func testAccCdnTopUrlStatistics_basic(startTime, endTime string) string {
	return fmt.Sprintf(`
data "huaweicloud_cdn_top_url_statistics" "test" {
  domain_name = "%[1]s"
  start_time  = "%[2]s"
  end_time    = "%[3]s"
  stat_type   = "req_num"
  service_area = "mainland_china"
}`, acceptance.HW_CDN_DOMAIN_NAME, startTime, endTime)
}
