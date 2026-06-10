package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccTaurusDBHtapStarrocksParametersModify_basic(t *testing.T) {
	resourceName := "huaweicloud_taurusdb_htap_starrocks_parameters_modify.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBInstanceId(t)
			acceptance.TestAccPreCheckTaurusDBHtapInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBHtapStarrocksParametersModify_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "taurusdb_instance_id", acceptance.HW_TAURUSDB_INSTANCE_ID),
					resource.TestCheckResourceAttr(resourceName, "starrocks_instance_id", acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID),
					resource.TestCheckResourceAttr(resourceName, "node_type", "be"),
				),
			},
		},
	})
}

func TestAccTaurusDBHtapStarrocksParametersModify_feParams(t *testing.T) {
	resourceName := "huaweicloud_taurusdb_htap_starrocks_parameters_modify.fe_params"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBInstanceId(t)
			acceptance.TestAccPreCheckTaurusDBHtapInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBHtapStarrocksParametersModify_feParams(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "taurusdb_instance_id", acceptance.HW_TAURUSDB_INSTANCE_ID),
					resource.TestCheckResourceAttr(resourceName, "starrocks_instance_id", acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID),
					resource.TestCheckResourceAttr(resourceName, "node_type", "fe"),
				),
			},
		},
	})
}

func testAccTaurusDBHtapStarrocksParametersModify_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_taurusdb_htap_starrocks_parameters_modify" "test" {
  taurusdb_instance_id  = "%[1]s"
  starrocks_instance_id = "%[2]s"
  node_type             = "be"

  parameter_values = {
    "alter_tablet_worker_count"            = "5"
	"base_compaction_num_threads_per_disk" = "5"
  }
}
`, acceptance.HW_TAURUSDB_INSTANCE_ID, acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID)
}

func testAccTaurusDBHtapStarrocksParametersModify_feParams() string {
	return fmt.Sprintf(`
resource "huaweicloud_taurusdb_htap_starrocks_parameters_modify" "fe_params" {
  taurusdb_instance_id  = "%[1]s"
  starrocks_instance_id = "%[2]s"
  node_type             = "fe"

  parameter_values = {
    "alter_table_timeout_second"     = "21600"
    "bdbje_heartbeat_timeout_second" = "10"
  }
}
`, acceptance.HW_TAURUSDB_INSTANCE_ID, acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID)
}
