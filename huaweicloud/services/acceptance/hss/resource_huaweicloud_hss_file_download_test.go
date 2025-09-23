package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccFileDownload_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a host ID with host protection enabled,
			// and the host is under the default enterprise project.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testFileDownload_basic(),
			},
		},
	})
}

func testFileDownload_base() string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_vulnerability_information_export" "test" {
  export_size           = 200
  category              = "vul"
  type                  = "linux_vul"
  host_id               = "%[1]s"
  repair_priority       = "Medium"
  handle_status         = "unhandled"
  asset_value           = "test"
  enterprise_project_id = "0"

  export_headers = [
    ["vul_id", "Vulnerability ID"],
    ["vul_name", "Vulnerability Name"]
  ]
}
`, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}

func testFileDownload_basic() string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_hss_file_download" "test" {
  file_id               = huaweicloud_hss_vulnerability_information_export.test.file_id
  enterprise_project_id = "0"
  export_file_name      = "test_export"
}
`, testFileDownload_base())
}
