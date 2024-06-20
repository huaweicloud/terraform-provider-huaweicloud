package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCssLogBackupRecords_basic(t *testing.T) {
	dataSource := "data.huaweicloud_css_log_backup_records.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCssLogBackupRecords_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.status"),

					resource.TestCheckOutput("id_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCssLogBackupRecords_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_css_log_backup_records" "test" {
  depends_on = [huaweicloud_css_manual_log_backup.test2]

  cluster_id = huaweicloud_css_cluster.test.id
}

output "records" {
  value = data.huaweicloud_css_log_backup_records.test.records
}

locals {
  records = data.huaweicloud_css_log_backup_records.test.records
  job_id  = local.records[0].id
  type    = local.records[0].type
  status  = local.records[0].status
}

data "huaweicloud_css_log_backup_records" "filter_by_id" {
  cluster_id = huaweicloud_css_cluster.test.id
  job_id     = local.job_id
}

data "huaweicloud_css_log_backup_records" "filter_by_type" {
  cluster_id = huaweicloud_css_cluster.test.id
  type       = local.type
}

data "huaweicloud_css_log_backup_records" "filter_by_status" {
  cluster_id = huaweicloud_css_cluster.test.id
  status     = local.status
}

locals {
  list_by_id     = data.huaweicloud_css_log_backup_records.filter_by_id.records
  list_by_type   = data.huaweicloud_css_log_backup_records.filter_by_type.records
  list_by_status = data.huaweicloud_css_log_backup_records.filter_by_status.records
}

output "id_filter_is_useful" {
  value = length(local.list_by_id) > 0 && alltrue(
    [for v in local.list_by_id[*].id : v == local.job_id]
  )
}

output "type_filter_is_useful" {
  value = length(local.list_by_type) > 0 && alltrue(
    [for v in local.list_by_type[*].type : v == local.type]
  )
}

output "status_filter_is_useful" {
  value = length(local.list_by_status) > 0 && alltrue(
    [for v in local.list_by_status[*].status : v == local.status]
  )
}
`, testCssLogBackupRecords_dataBasic(name))
}

func testCssLogBackupRecords_dataBasic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_css_manual_log_backup" "test" {
  depends_on = [huaweicloud_css_log_setting.test]
  
  cluster_id = huaweicloud_css_cluster.test.id
}

resource "huaweicloud_css_manual_log_backup" "test1" {
  depends_on = [huaweicloud_css_manual_log_backup.test]
  
  cluster_id = huaweicloud_css_cluster.test.id
}

resource "huaweicloud_css_manual_log_backup" "test2" {
  depends_on = [huaweicloud_css_manual_log_backup.test1]
  
  cluster_id = huaweicloud_css_cluster.test.id
}
`, testLogSetting_elastic(name))
}
