package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceWebtamperRaspPath_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_webtamper_rasp_path.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Setting a host ID with web tamper protection enabled.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceWebtamperRaspPath_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "rasp_path"),
				),
			},
		},
	})
}

func testAccDataSourceWebtamperRaspPath_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_hss_webtamper_rasp_path" "test" {
  host_id               = "%[1]s"
  enterprise_project_id = "0"
}
`, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}
