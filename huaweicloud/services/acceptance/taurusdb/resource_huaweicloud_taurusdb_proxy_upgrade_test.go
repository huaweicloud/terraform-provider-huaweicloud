package taurusdb

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccTaurusDBProxyUpgrade_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBProxyUpgrade_basic(),
				// No available version to upgrade, expect error with message 'Parameter error'
				ExpectError: regexp.MustCompile(`Parameter error`),
			},
		},
	})
}

func testAccTaurusDBProxyUpgrade_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_taurusdb_proxies" "test" {
  instance_id = "%[1]s"
}

locals {
  proxy_id = data.huaweicloud_taurusdb_proxies.test.proxy_list[0].id
}

resource "huaweicloud_taurusdb_proxy_upgrade" "test" {
  instance_id    = "%[1]s"
  proxy_id       = local.proxy_id
  source_version = "2.24.03.000"
  target_version = "2.26.03.000"
}
`, acceptance.HW_TAURUSDB_INSTANCE_ID)
}
