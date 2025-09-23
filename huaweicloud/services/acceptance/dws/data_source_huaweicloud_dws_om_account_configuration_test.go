package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOmAccountConfiguration_basic(t *testing.T) {
	var (
		dataSource   = "data.huaweicloud_dws_om_account_configuration.test"
		dc           = acceptance.InitDataSourceCheck(dataSource)
		onAccount    = "data.huaweicloud_dws_om_account_configuration.on_account"
		dcOnAccount  = acceptance.InitDataSourceCheck(onAccount)
		offAccount   = "data.huaweicloud_dws_om_account_configuration.off_account"
		dcOffAccount = acceptance.InitDataSourceCheck(offAccount)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccOmAccountConfiguration_clusterNotExist(),
				ExpectError: regexp.MustCompile("Cluster does not exist or has been deleted"),
			},
			{
				Config: testAccOmAccountConfiguration_basic(),
				Check: resource.ComposeTestCheckFunc(
					dcOnAccount.CheckResourceExists(),
					resource.TestCheckResourceAttr(onAccount, "status", "on"),
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_time_equal", "true"),
					dcOffAccount.CheckResourceExists(),
					resource.TestCheckResourceAttr(offAccount, "status", "off"),
				),
			},
		},
	})
}

func testAccOmAccountConfiguration_clusterNotExist() string {
	clusterId, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_dws_om_account_configuration" "test" {
  cluster_id = "%s"
}
`, clusterId)
}

func testAccOmAccountConfiguration_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dws_om_account_action" "test" {
  cluster_id = "%[1]s"
  operation  = "addOmUser"
}

data "huaweicloud_dws_om_account_configuration" "on_account" {
  depends_on = [huaweicloud_dws_om_account_action.test]
  cluster_id = "%[1]s"
}

# Extended 2 times (16 hours).
resource "huaweicloud_dws_om_account_action" "increase_period" {
  count      = 2
  depends_on = [huaweicloud_dws_om_account_action.test]
  cluster_id = "%[1]s"
  operation  = "increaseOmUserPeriod"
}

data "huaweicloud_dws_om_account_configuration" "test" {
  depends_on = [huaweicloud_dws_om_account_action.increase_period]
  cluster_id = "%[1]s"
}

# Get the time when the account switch is turned on and add 16 hours to it.
locals {
  period_time = timeadd(data.huaweicloud_dws_om_account_configuration.on_account.om_user_expires_time, "16h")
}

# Assert that the extended time is equal to the original time plus 16 hours.
output "is_time_equal" {
  value = timecmp(local.period_time, data.huaweicloud_dws_om_account_configuration.test.om_user_expires_time) == 0
}

resource "huaweicloud_dws_om_account_action" "off_account" {
  depends_on = [huaweicloud_dws_om_account_action.increase_period]
  cluster_id = "%[1]s"
  operation  = "deleteOmUser"
}

data "huaweicloud_dws_om_account_configuration" "off_account" {
  depends_on = [huaweicloud_dws_om_account_action.off_account]
  cluster_id = "%[1]s"
}
`, acceptance.HW_DWS_CLUSTER_ID)
}
