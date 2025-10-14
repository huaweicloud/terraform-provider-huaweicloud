package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccModifyWebtamperRaspPath_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a host ID that has enabled web tamper protection.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testModifyWebtamperRaspPath_basic(),
			},
		},
	})
}

func testModifyWebtamperRaspPath_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_modify_webtamper_rasp_path" "test" {
  host_id               = "%s"
  rasp_path             = "/usr/workspace/apache-tomcat-8.5.15/bin"
  enterprise_project_id = "0"
}
`, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}
