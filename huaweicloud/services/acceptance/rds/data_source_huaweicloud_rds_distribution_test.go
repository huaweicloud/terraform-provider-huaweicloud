package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsDistribution_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_distribution.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsDistribution_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "status"),
					resource.TestCheckResourceAttrSet(dataSource, "distributor_instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "distributor_instance_name"),
				),
			},
		},
	})
}

func testDataSourceRdsDistribution_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rds_distribution" "test" {
  instance_id = "%s"
}
`, acceptance.HW_RDS_INSTANCE_ID)
}
