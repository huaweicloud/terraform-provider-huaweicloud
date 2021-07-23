package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccELBV2LoadbalancerDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBV2LoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccELBV2LoadbalancerDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckELBV2LoadbalancerDataSourceID("data.huaweicloud_lb_loadbalancer.test_by_name"),
					testAccCheckELBV2LoadbalancerDataSourceID("data.huaweicloud_lb_loadbalancer.test_by_description"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_lb_loadbalancer.test_by_name", "name", rName),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_lb_loadbalancer.test_by_description", "name", rName),
				),
			},
		},
	})
}

func testAccCheckELBV2LoadbalancerDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find elb load balancer data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("load balancer data source ID not set")
		}

		return nil
	}
}

func testAccELBV2LoadbalancerDataSource_basic(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "huaweicloud_lb_loadbalancer" "test" {
  name          = "%s"
  vip_subnet_id = data.huaweicloud_vpc_subnet.test.subnet_id
  description   = "test for load balancer data source"
}

data "huaweicloud_lb_loadbalancer" "test_by_name" {
  name = huaweicloud_lb_loadbalancer.test.name
}

data "huaweicloud_lb_loadbalancer" "test_by_description" {
  description = huaweicloud_lb_loadbalancer.test.description
}
`, rName)
}
