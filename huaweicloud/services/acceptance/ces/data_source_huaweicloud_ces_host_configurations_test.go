package ces

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCesHostDataPoints_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ces_host_configurations.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCesDem0(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCesHostDataPoints_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "datapoints.#"),
					resource.TestCheckResourceAttrSet(dataSource, "datapoints.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "datapoints.0.timestamp"),
					resource.TestCheckResourceAttrSet(dataSource, "datapoints.0.value"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCesHostDataPoints_basic() string {
	currentTime := time.Now()
	tenMinutesAgo := currentTime.Add(-10*time.Minute).Unix() * 1e3
	tenMinutesLater := currentTime.Add(10*time.Minute).Unix() * 1e3
	return fmt.Sprintf(`
data "huaweicloud_ces_host_configurations" "test" {
  namespace = "SYS.ECS"
  type      = "instance_host_info"
  from      = %[1]v
  to        = %[2]v
  dim_0     = "%[3]s"
}
`, tenMinutesAgo, tenMinutesLater, acceptance.HW_CES_DEM_0)
}
