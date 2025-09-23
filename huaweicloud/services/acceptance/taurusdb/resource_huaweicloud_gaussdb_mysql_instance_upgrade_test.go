package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGaussDBMysqlInstanceUpgrade_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBMysqlInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBMysqlInstanceUpgrade_basic(),
			},
		},
	})
}

func testAccGaussDBMysqlInstanceUpgrade_basic() string {
	return fmt.Sprintf(`


resource "huaweicloud_gaussdb_mysql_instance_upgrade" "test" {
  instance_id = "%[1]s"
  delay       = true
}`, acceptance.HW_GAUSSDB_MYSQL_INSTANCE_ID)
}
