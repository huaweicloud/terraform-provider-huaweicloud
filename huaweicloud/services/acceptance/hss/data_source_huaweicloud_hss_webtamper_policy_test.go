package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceWebTamperPolicy_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_webtamper_policy.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a host ID with web tamper protection enabled under the
			// default enterprise project.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceWebTamperPolicy_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "protect_dir_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "protect_dir_info.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "protect_dir_info.0.protect_dir_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "protect_dir_info.0.protect_dir_list.0.protect_dir"),
					resource.TestCheckResourceAttrSet(dataSourceName, "protect_dir_info.0.protect_dir_list.0.local_backup_dir"),
					resource.TestCheckResourceAttrSet(dataSourceName, "protect_dir_info.0.protect_dir_list.0.protect_status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "protect_dir_info.0.exclude_file_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "protect_dir_info.0.protect_mode"),
					resource.TestCheckResourceAttrSet(dataSourceName, "enable_timing_off"),
					resource.TestCheckResourceAttrSet(dataSourceName, "enable_rasp_protect"),
					resource.TestCheckResourceAttrSet(dataSourceName, "enable_privileged_process"),
					resource.TestCheckResourceAttrSet(dataSourceName, "privileged_child_status"),
				),
			},
		},
	})
}

func testAccDataSourceWebTamperPolicy_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_hss_webtamper_policy" "test" {
  host_id               = "%s"
  enterprise_project_id = "0"
}
`, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}
