package dcs

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDcsSlowLogs_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_dcs_slow_logs.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)
	name := acceptance.RandomAccResourceName()

	now := time.Now().UnixNano() / 1e6
	startTime := now - 86400000
	endTime := now

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDcsSlowLogs_basic(name, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "slowlogs.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "slowlogs.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "slowlogs.0.command"),
					resource.TestCheckResourceAttrSet(dataSourceName, "slowlogs.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "slowlogs.0.duration"),
					resource.TestCheckResourceAttrSet(dataSourceName, "slowlogs.0.shard_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "slowlogs.0.database_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "slowlogs.0.username"),
					resource.TestCheckResourceAttrSet(dataSourceName, "slowlogs.0.node_role"),
					resource.TestCheckResourceAttrSet(dataSourceName, "slowlogs.0.client_ip"),
				),
			},
		},
	})
}

func testAccDataSourceDcsSlowLogs_basic(name string, startTime, endTime int64) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dcs_slow_logs" "test" {
  depends_on  = [data.huaweicloud_dcs_redis_run_logs.test]

  instance_id = huaweicloud_dcs_instance.test.id
  sort_key    = "duration"
  sort_dir    = "desc"
  start_time  = "%d"
  end_time    = "%d"
}
`, testDataSourceRedisRunLogs_basic(name), startTime, endTime)
}
