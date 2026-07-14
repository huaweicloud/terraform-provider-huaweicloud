package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccTaurusDBHtapStarrocksParametersCompare_basic(t *testing.T) {
	var (
		obj          interface{}
		rName        = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_taurusdb_htap_starrocks_parameters_compare.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getHtapStarrocksInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBHtapStarrocksParametersCompare_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "differences.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "differences.0.parameter_name", "alter_tablet_worker_count"),
					resource.TestCheckResourceAttr(resourceName, "differences.0.source_value", "1"),
					resource.TestCheckResourceAttr(resourceName, "differences.0.target_value", "3"),
				),
			},
		},
	})
}

func testAccTaurusDBHtapStarrocksParametersCompare_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  intance_config_id = huaweicloud_taurusdb_htap_starrocks_instance.test.be_configurations[0].configuration_id
}

resource "huaweicloud_taurusdb_htap_starrocks_parameters_compare" "test" {
  source_configuration_id = local.intance_config_id
}
`, testAccTaurusDBHtapStarrocksInstance_basic(rName))
}
