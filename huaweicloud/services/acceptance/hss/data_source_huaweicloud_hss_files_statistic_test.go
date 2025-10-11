package hss

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceFilesStatistic_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_files_statistic.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
		startTime  = time.Now().Add(-24*time.Hour).UnixNano() / int64(time.Millisecond)
		endTime    = time.Now().UnixNano() / int64(time.Millisecond)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Setting a host ID with host protection enabled.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceFilesStatistic_basic(startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "host_num"),
					resource.TestCheckResourceAttrSet(dataSource, "change_total_num"),
					resource.TestCheckResourceAttrSet(dataSource, "change_file_num"),
					resource.TestCheckResourceAttrSet(dataSource, "change_registry_num"),
					resource.TestCheckResourceAttrSet(dataSource, "modify_num"),
					resource.TestCheckResourceAttrSet(dataSource, "add_num"),
					resource.TestCheckResourceAttrSet(dataSource, "delete_num"),
				),
			},
		},
	})
}

func testAccDataSourceFilesStatistic_basic(startTime, endTime int64) string {
	return fmt.Sprintf(`
data "huaweicloud_hss_files_statistic" "test" {
  begin_time            = %[1]d
  end_time              = %[2]d
  enterprise_project_id = "0"
}
`, startTime, endTime)
}
