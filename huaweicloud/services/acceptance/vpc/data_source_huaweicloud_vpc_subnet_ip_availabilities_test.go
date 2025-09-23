package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpcSubnetIpAvailabilities_basic(t *testing.T) {
	dataSource := "data.huaweicloud_vpc_subnet_ip_availabilities.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceVpcSubnetIpAvailabilities_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceVpcSubnetIpAvailabilities_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_vpc_subnet_ip_availabilities" "test" {
  network_id = huaweicloud_vpc_subnet.test.id
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_vpc_subnet_ip_availabilities.test.network_ip_availability) > 0
}
`, testAccVpcSubnetV1_basic(name))
}
