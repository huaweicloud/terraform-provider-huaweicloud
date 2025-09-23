package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpcSubnetPrivateIps_basic(t *testing.T) {
	dataSource := "data.huaweicloud_vpc_subnet_private_ips.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceVpcSubnetPrivateIps_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceVpcSubnetPrivateIps_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpc_subnet_private_ips" "test" {
  subnet_id = huaweicloud_vpc_subnet.test.id

  depends_on = [huaweicloud_vpc_subnet_private_ip.test]
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_vpc_subnet_private_ips.test.private_ips) > 0
}
`, testAccSubnetPrivateIP_basic(name))
}
