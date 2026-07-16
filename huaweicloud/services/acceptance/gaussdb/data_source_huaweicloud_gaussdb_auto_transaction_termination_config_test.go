package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussdbAutoTransactionTerminationConfig_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_auto_transaction_termination_config.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbAutoTransactionTerminationConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "type"),
					resource.TestCheckResourceAttrSet(dataSource, "auto_stop"),
					resource.TestCheckResourceAttrSet(dataSource, "usernames.#"),
					resource.TestCheckResourceAttrSet(dataSource, "threshold"),
					resource.TestCheckResourceAttrSet(dataSource, "database_names.#"),
					resource.TestCheckResourceAttrSet(dataSource, "database_names_select.#"),
				),
			},
		},
	})
}

func testDataSourceGaussdbAutoTransactionTerminationConfig_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_auto_transaction_termination_config" "test" {
  instance_id = "%s"
  type        = "exec_time"
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID)
}
