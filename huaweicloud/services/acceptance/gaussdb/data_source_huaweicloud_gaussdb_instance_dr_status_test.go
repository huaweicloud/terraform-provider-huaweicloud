package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGaussDbInstanceDrStatus_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_gaussdb_instance_dr_status.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDbInstanceDrStatus_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rpo"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rto"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rpo_threshold"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rto_threshold"),
				),
			},
		},
	})
}

func testAccGaussDbInstanceDrStatus_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_instance_dr_status" "test" {
  instance_id   = "%[1]s"
  disaster_type = "stream"
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID)
}
