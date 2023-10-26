package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceListeners_basic(t *testing.T) {
	rName := "data.huaweicloud_elb_listeners.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceListeners_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "listeners.#"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.name"),
					//resource.TestCheckResourceAttrSet(rName, "listeners.0.id"),
					//resource.TestCheckResourceAttrSet(rName, "listeners.0.ipv4_address"),
					//resource.TestCheckResourceAttrSet(rName, "listeners.0.ipv4_port_id"),
					//resource.TestCheckResourceAttrSet(rName, "listeners.0.l4_flavor_id"),
					//resource.TestCheckResourceAttrSet(rName, "listeners.0.l7_flavor_id"),
					//resource.TestCheckResourceAttrSet(rName, "listeners.0.vpc_id"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					//resource.TestCheckOutput("vpc_id_filter_is_useful", "true"),
					//resource.TestCheckOutput("ipv4_subnet_id_filter_is_useful", "true"),
					//resource.TestCheckOutput("description_filter_is_useful", "true"),
					//resource.TestCheckOutput("l4_flavor_id_filter_is_useful", "true"),
					//resource.TestCheckOutput("l7_flavor_id_filter_is_useful", "true"),
					//resource.TestCheckOutput("type_is_useful", "true"),
				),
			},
		},
	})
}

func testAccElbListenerConfig_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc_subnet" "test" {}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_elb_loadbalancer" "test" {
  name           = "%[2]s"
  vpc_id         = huaweicloud_vpc.test.id
  ipv4_subnet_id = huaweicloud_vpc_subnet.test.ipv4_subnet_id
  
  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]
  backend_subnets = [
    huaweicloud_vpc_subnet.test.id
  ]
  tags = {
    key   = "value"
    owner = "terraform"
  }
}

resource "huaweicloud_elb_listener" "test" {
 name                        = "%[1]s"
 description                 = "test description"
 protocol                    = "HTTP"
 protocol_port               = 8080
 loadbalancer_id             = huaweicloud_elb_loadbalancer.test.id
 advanced_forwarding_enabled = false

 idle_timeout     = 62
 request_timeout  = 63
 response_timeout = 64

 tags = {
   key   = "value"
   owner = "terraform"
 }
}
`, name, name)
}

func testAccDatasourceListeners_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_elb_listeners" "test" {
  depends_on = [huaweicloud_elb_listener.test]
}

data "huaweicloud_elb_listeners" "name_filter" {
  name       = "%[2]s"
  depends_on = [huaweicloud_elb_listener.test]
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_elb_listeners.name_filter.listeners) > 0 && alltrue(
  [for v in data.huaweicloud_elb_listeners.name_filter.listeners[*].name :v == "%[2]s"]
  )  
}

`, testAccElbListenerConfig_basic(name), name)
}
