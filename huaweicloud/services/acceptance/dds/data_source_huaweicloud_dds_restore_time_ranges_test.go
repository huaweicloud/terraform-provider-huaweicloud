package dds

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDdsRestoreTimeRanges_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dds_restore_time_ranges.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDDSInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDdsRestoreTimeRanges_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "restore_times.#"),
					resource.TestCheckResourceAttrSet(dataSource, "restore_times.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "restore_times.0.end_time"),
				),
			},
		},
	})
}

func testDataSourceDdsRestoreTimeRanges_basic() string {
	date := strings.Split(time.Now().Format("2006-01-02T15:04:05Z"), "T")[0]

	return fmt.Sprintf(`
data "huaweicloud_dds_restore_time_ranges" "test" {
  instance_id = "%s"
  date        = "%s"
}
`, acceptance.HW_DDS_INSTANCE_ID, date)
}
