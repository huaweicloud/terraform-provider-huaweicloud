package hss

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceWebtamperStaticProtectHistory_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_webtamper_static_protect_history.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
		startTime  = time.Now().Add(-24*time.Hour).UnixNano() / int64(time.Millisecond)
		endTime    = time.Now().UnixNano() / int64(time.Millisecond)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			// Because there is no available data for testing, the test case is only
			// used to verify that the API can be invoked.
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceWebtamperStaticProtectHistory_basic(startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
				),
			},
		},
	})
}

func testAccDataSourceWebtamperStaticProtectHistory_basic(startTime, endTime int64) string {
	return fmt.Sprintf(`
data "huaweicloud_hss_webtamper_static_protect_history" "test" {
  start_time            = %[1]d
  end_time              = %[2]d
  enterprise_project_id = "0"
}
`, startTime, endTime)
}
