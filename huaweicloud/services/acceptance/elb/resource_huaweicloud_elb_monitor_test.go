package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/elb/v3/monitors"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getELBMonitorResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ElbV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ELB client: %s", err)
	}
	return monitors.Get(client, state.Primary.ID).Extract()
}

func TestAccElbV3Monitor_basic(t *testing.T) {
	var monitor monitors.Monitor
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_elb_monitor.monitor_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&monitor,
		getELBMonitorResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3MonitorConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "interval", "20"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "10"),
					resource.TestCheckResourceAttr(resourceName, "max_retries", "5"),
					resource.TestCheckResourceAttr(resourceName, "max_retries_down", "5"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "url_path", "/aa"),
					resource.TestCheckResourceAttr(resourceName, "domain_name", "www.aa.com"),
					resource.TestCheckResourceAttr(resourceName, "port", "8000"),
					resource.TestCheckResourceAttr(resourceName, "status_code", "200,401-500,502"),
					resource.TestCheckResourceAttr(resourceName, "http_method", "GET"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
			{
				Config: testAccElbV3MonitorConfig_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "interval", "30"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "20"),
					resource.TestCheckResourceAttr(resourceName, "max_retries", "8"),
					resource.TestCheckResourceAttr(resourceName, "max_retries_down", "8"),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "url_path", "/bb"),
					resource.TestCheckResourceAttr(resourceName, "domain_name", "www.bb.com"),
					resource.TestCheckResourceAttr(resourceName, "port", "8888"),
					resource.TestCheckResourceAttr(resourceName, "status_code", "200,301,404-500,504"),
					resource.TestCheckResourceAttr(resourceName, "http_method", "HEAD"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccElbV3Monitor_with_ip_protocol(t *testing.T) {
	var monitor monitors.Monitor
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_elb_monitor.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&monitor,
		getELBMonitorResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckElbGatewayType(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3MonitorConfig_with_ip_protocol(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "port", "6666"),
					resource.TestCheckResourceAttr(resourceName, "interval", "20"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "10"),
					resource.TestCheckResourceAttr(resourceName, "max_retries", "5"),
				),
			},
			{
				Config: testAccElbV3MonitorConfig_with_ip_protocol_update(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "port", "8888"),
					resource.TestCheckResourceAttr(resourceName, "interval", "30"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "20"),
					resource.TestCheckResourceAttr(resourceName, "max_retries", "8"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccElbV3MonitorConfig_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_elb_monitor" "monitor_1" {
  pool_id          = huaweicloud_elb_pool.test.id
  name             = "%s"
  protocol         = "HTTP"
  interval         = 20
  timeout          = 10
  max_retries      = 5
  max_retries_down = 5
  url_path         = "/aa"
  domain_name      = "www.aa.com"
  port             = "8000"
  status_code      = "200,401-500,502"
  enabled          = false
}
`, testAccElbV3PoolConfig_basic(rName), rName)
}

func testAccElbV3MonitorConfig_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_elb_monitor" "monitor_1" {
  pool_id          = huaweicloud_elb_pool.test.id
  name             = "%s-update"
  protocol         = "HTTPS"
  interval         = 30
  timeout          = 20
  max_retries      = 8
  max_retries_down = 8
  url_path         = "/bb"
  domain_name      = "www.bb.com"
  port             = 8888
  status_code      = "200,301,404-500,504"
  http_method      = "HEAD"
  enabled          = true
}
`, testAccElbV3PoolConfig_basic(rName), rName)
}

func testAccElbV3MonitorConfig_with_ip_protocol_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_elb_loadbalancer" "test" {
  name              = "%[2]s"
  vpc_id            = huaweicloud_vpc.test.id
  ipv4_subnet_id    = huaweicloud_vpc_subnet.test.ipv4_subnet_id
  loadbalancer_type = "gateway"
  description       = "test gateway description"

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]
}

resource "huaweicloud_elb_pool" "test" {
  name            = "%[2]s"
  protocol        = "IP"
  lb_method       = "2_TUPLE_HASH"
  loadbalancer_id = huaweicloud_elb_loadbalancer.test.id
}
`, common.TestVpc(rName), rName)
}

func testAccElbV3MonitorConfig_with_ip_protocol(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_elb_monitor" "test" {
  pool_id     = huaweicloud_elb_pool.test.id
  name        = "%[2]s"
  port        = 6666
  protocol    = "TCP"
  interval    = 20
  timeout     = 10
  max_retries = 5
}
`, testAccElbV3MonitorConfig_with_ip_protocol_base(rName), rName)
}

func testAccElbV3MonitorConfig_with_ip_protocol_update(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_elb_monitor" "test" {
  pool_id     = huaweicloud_elb_pool.test.id
  name        = "%[2]s"
  port        = 8888
  protocol    = "TCP"
  interval    = 30
  timeout     = 20
  max_retries = 8
}
`, testAccElbV3MonitorConfig_with_ip_protocol_base(rName), rNameUpdate)
}
