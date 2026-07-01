package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGaussDBInstanceAutoEnlargePolicy_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_instance_auto_enlarge_policy.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBInstanceAutoEnlargePolicy_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "switch_option"),
					resource.TestCheckResourceAttrSet(dataSource, "limit_volume_size"),
					resource.TestCheckResourceAttrSet(dataSource, "min_volume_size"),
					resource.TestCheckResourceAttrSet(dataSource, "max_volume_size"),
					resource.TestCheckResourceAttrSet(dataSource, "trigger_available_percent"),
					resource.TestCheckResourceAttrSet(dataSource, "percents.#"),
					resource.TestCheckResourceAttrSet(dataSource, "step_size"),
				),
			},
		},
	})
}

func testAccGaussDBInstanceAutoEnlargePolicy_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_instance_auto_enlarge_policy" "test" {
  instance_id = "%s"
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID)
}
