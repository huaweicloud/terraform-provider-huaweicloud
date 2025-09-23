package ddm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceDdmAvailableRdsInstances_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ddm_available_rds_instances.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDdmAvailableRdsInstances_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.mode"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.port"),
				),
			},
		},
	})
}

func testAccDatasourceDdmAvailableRdsInstances_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_ddm_available_rds_instances" "test" {
  instance_id = "%s"
}
`, acceptance.HW_DDM_INSTANCE_ID)
}
