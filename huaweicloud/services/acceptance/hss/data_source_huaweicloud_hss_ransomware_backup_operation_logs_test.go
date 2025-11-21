package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Due to testing environment limitations, this test case can only test the scenario with empty `data_list`.
func TestAccDataSourceRansomwareBackupOperationLogs_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_ransomware_backup_operation_logs.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRansomwareBackupOperationLogs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.#"),
				),
			},
		},
	})
}

func testDataSourceRansomwareBackupOperationLogs_basic() string {
	return `
data "huaweicloud_hss_ransomware_backup_operation_logs" "test" {
  enterprise_project_id = "0"
}
`
}
