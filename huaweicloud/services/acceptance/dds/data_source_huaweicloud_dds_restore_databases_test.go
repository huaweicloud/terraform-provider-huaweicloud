package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDdsRestoreDatabases_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dds_restore_databases.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDDSInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDdsRestoreDatabases_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "databases.#"),
				),
			},
		},
	})
}

func testDataSourceDdsRestoreDatabases_basic() string {
	return fmt.Sprintf(`
%[1]s

locals {
  end_time = try(data.huaweicloud_dds_restore_time_ranges.test.restore_times.0.end_time, 0)
}

data "huaweicloud_dds_restore_databases" "test" {
  instance_id  = "%[2]s"
  restore_time = local.end_time
}
`, testDataSourceDdsRestoreTimeRanges_basic(), acceptance.HW_DDS_INSTANCE_ID)
}
