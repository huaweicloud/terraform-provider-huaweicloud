package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsMergedBinlogFiles_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_merged_binlog_files.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsMergedBinlogFiles_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "pack_log_infos.#"),
					resource.TestCheckResourceAttrSet(dataSource, "pack_log_infos.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "pack_log_infos.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "pack_log_infos.0.size"),
					resource.TestCheckResourceAttrSet(dataSource, "pack_log_infos.0.size_unit"),
					resource.TestCheckResourceAttrSet(dataSource, "pack_log_infos.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "pack_log_infos.0.query_start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "pack_log_infos.0.query_end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "pack_log_infos.0.file_name"),
				),
			},
		},
	})
}

func testDataSourceRdsMergedBinlogFiles_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rds_merged_binlog_files" "test" {
  instance_id = "%s"
}
`, acceptance.HW_RDS_INSTANCE_ID)
}
