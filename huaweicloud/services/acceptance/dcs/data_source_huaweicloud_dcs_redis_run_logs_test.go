package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRedisRunLogs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dcs_redis_run_logs.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRedisRunLogs_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "file_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "file_list.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "file_list.0.file_name"),
					resource.TestCheckResourceAttrSet(dataSource, "file_list.0.group_name"),
					resource.TestCheckResourceAttrSet(dataSource, "file_list.0.replication_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "file_list.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "file_list.0.time"),
					resource.TestCheckResourceAttrSet(dataSource, "file_list.0.backup_id"),
				),
			},
		},
	})
}

func testDataSourceRedisRunLogs_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dcs_redis_run_logs" "test" {
  depends_on = [huaweicloud_dcs_redis_run_log_collect.test]

  instance_id = huaweicloud_dcs_instance.test.id
  log_type    = "run"
}
`, testRedisRunLogCollect_basic(name))
}
