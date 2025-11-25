package cdn

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataTopReferrerStatistics_basic(t *testing.T) {
	var (
		now     = time.Now()
		today   = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		endTime = today.Format(time.RFC3339)

		dcName = "data.huaweicloud_cdn_top_referrer_statistics.test"
		dc     = acceptance.InitDataSourceCheck(dcName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCdnDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataTopReferrerStatistics_basic(endTime),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "statistics.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dcName, "statistics.0.refer"),
					resource.TestCheckResourceAttrSet(dcName, "statistics.0.value"),
				),
			},
		},
	})
}

func testAccDataTopReferrerStatistics_basic(endTime string) string {
	return fmt.Sprintf(`
data "huaweicloud_cdn_top_referrer_statistics" "test" {
  domain_name = "%[1]s"
  start_time  = timeadd("%[2]s", "-120h")
  end_time    = "%[2]s"
  stat_type   = "req_num"
}
`, acceptance.HW_CDN_DOMAIN_NAME, endTime)
}
