package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpcSubnetTags_basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr, randGatewayIp := acceptance.RandomCidrAndGatewayIp()
	dataSourceName := "data.huaweicloud_vpc_subnet_tags.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceVpcSubnetTags_basic(randName, randCidr, randGatewayIp),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
				),
			},
		},
	})
}

func testDataSourceVpcSubnetTags_basic(rName, cidr, gatewayIp string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpc_subnet_tags" "test" {
  depends_on = [ huaweicloud_vpc_subnet.test ]
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_vpc_subnet_tags.test.tags) > 0
}
`, testAccVpcSubnetsDataSource_Base(rName, cidr, gatewayIp))
}
