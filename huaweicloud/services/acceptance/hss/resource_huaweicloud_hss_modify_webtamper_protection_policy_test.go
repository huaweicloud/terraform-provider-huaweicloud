package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccModifyWebTamperProtectionPolicy_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a host ID with web tamper protection enabled.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testModifyWebTamperProtectionPolicy_basic(),
			},
		},
	})
}

func testModifyWebTamperProtectionPolicy_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_modify_webtamper_protection_policy" "test" {
  host_id = "%s"

  protect_dir_info {
    protect_dir_list {
      protect_dir       = "/test/test1"
      local_backup_dir  = "/test/test2"
      exclude_child_dir = "pro"
      exclude_file_path = "path"
    }

    protect_dir_list {
      protect_dir      = "/test/test3"
      local_backup_dir = "/test/test4"
    }

    exclude_file_type = "log;text"
    protect_mode      = "recovery"
  }

  enable_timing_off = true

  timing_off_config_info {
    week_off_list = [1,7]

    timing_range_list {
      time_range  = "8:00-9:00"
      description = "close time range"
    }

    timing_range_list {
      time_range = "10:00-11:00"
    }
  }

  enable_rasp_protect       = true
  rasp_path                 = "/usr/bin/tomcat/bin"
  enable_privileged_process = true

  privileged_process_info {
    privileged_process_path_list = ["/usr/bin/echo"]
    privileged_child_status      = true
  }

  enterprise_project_id = "0"
}
`, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}
