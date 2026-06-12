package gaussdb

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussdbClientAuthConfigHistory_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_client_auth_config_history.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbClientAuthConfigHistory_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "hba_histories.#"),
					resource.TestCheckResourceAttrSet(dataSource, "hba_histories.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "hba_histories.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "hba_histories.0.time"),
					resource.TestCheckResourceAttrSet(dataSource, "hba_histories.0.before_confs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "hba_histories.0.after_confs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "hba_histories.0.before_confs.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "hba_histories.0.before_confs.0.database"),
					resource.TestCheckResourceAttrSet(dataSource, "hba_histories.0.before_confs.0.user"),
					resource.TestCheckResourceAttrSet(dataSource, "hba_histories.0.before_confs.0.address"),
					resource.TestCheckResourceAttrSet(dataSource, "hba_histories.0.before_confs.0.method"),
					resource.TestCheckResourceAttrSet(dataSource, "hba_histories.0.after_confs.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "hba_histories.0.after_confs.0.database"),
					resource.TestCheckResourceAttrSet(dataSource, "hba_histories.0.after_confs.0.user"),
					resource.TestCheckResourceAttrSet(dataSource, "hba_histories.0.after_confs.0.address"),
					resource.TestCheckResourceAttrSet(dataSource, "hba_histories.0.after_confs.0.method"),
					resource.TestCheckOutput("time_filter", "true"),
				),
			},
		},
	})
}

func testDataSourceGaussdbClientAuthConfigHistory_basic() string {
	now := time.Now().UTC()
	startTime := now.Format("2006-01-02") + " 00:00:00"
	endTime := now.Format("2006-01-02") + " 23:59:59"

	return fmt.Sprintf(`
data "huaweicloud_gaussdb_client_auth_config_history" "test" {
  instance_id = "%[1]s"
}

data "huaweicloud_gaussdb_client_auth_config_history" "time_filter" {
  instance_id = "%[1]s"
  start_time  = "%[2]s"
  end_time    = "%[3]s"
}
output "time_filter" {
  value = length(data.huaweicloud_gaussdb_client_auth_config_history.time_filter.hba_histories) > 0
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID, startTime, endTime)
}
