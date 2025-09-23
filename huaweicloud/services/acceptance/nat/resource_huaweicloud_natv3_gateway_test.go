package nat

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/nat"
)

func getPublicGatewayV3ResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NatGatewayClient(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating NAT v2 client: %s", err)
	}

	return nat.GetPublicGateway(client, state.Primary.ID)
}

func TestAccPublicGatewayV3_basic(t *testing.T) {
	var (
		obj           interface{}
		rName         = "huaweicloud_natv3_gateway.test"
		name          = acceptance.RandomAccResourceNameWithDash()
		updateName    = acceptance.RandomAccResourceNameWithDash()
		relatedConfig = common.TestBaseNetwork(name)
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPublicGatewayV3ResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPublicGatewayV3_basic_step_1(name, relatedConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "spec", "1"),
					resource.TestCheckResourceAttr(rName, "description", "Created by acc test"),
					resource.TestCheckResourceAttr(rName, "ngport_ip_address", "192.168.0.101"),
					resource.TestCheckResourceAttr(rName, "session_conf.0.tcp_session_expire_time", "1000"),
					resource.TestCheckResourceAttr(rName, "session_conf.0.udp_session_expire_time", "400"),
					resource.TestCheckResourceAttr(rName, "session_conf.0.icmp_session_expire_time", "20"),
					resource.TestCheckResourceAttr(rName, "session_conf.0.tcp_time_wait_time", "10"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "dnat_rules_limit"),
					resource.TestCheckResourceAttrSet(rName, "snat_rule_public_ip_limit"),
				),
			},
			{
				Config: testAccPublicGatewayV3_basic_step_2(updateName, relatedConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "spec", "2"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "ngport_ip_address", "192.168.0.101"),
					resource.TestCheckResourceAttr(rName, "session_conf.0.tcp_session_expire_time", "900"),
					resource.TestCheckResourceAttr(rName, "session_conf.0.udp_session_expire_time", "300"),
					resource.TestCheckResourceAttr(rName, "session_conf.0.icmp_session_expire_time", "10"),
					resource.TestCheckResourceAttr(rName, "session_conf.0.tcp_time_wait_time", "5"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "baaar"),
					resource.TestCheckResourceAttr(rName, "tags.newkey", "value"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPublicGatewayV3_basic_step_1(name, relatedConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_natv3_gateway" "test" {
  name              = "%[2]s"
  spec              = "1"
  description       = "Created by acc test"
  ngport_ip_address = "192.168.0.101"
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id

  session_conf {
    tcp_session_expire_time  = 1000
    udp_session_expire_time  = 400
    icmp_session_expire_time = 20
    tcp_time_wait_time       = 10
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, relatedConfig, name)
}

func testAccPublicGatewayV3_basic_step_2(name, relatedConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_natv3_gateway" "test" {
  name              = "%[2]s"
  spec              = "2"
  ngport_ip_address = "192.168.0.101"
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id

  session_conf {
    tcp_session_expire_time  = 900
    udp_session_expire_time  = 300
    icmp_session_expire_time = 10
    tcp_time_wait_time       = 5
  }

  tags = {
    foo    = "baaar"
    newkey = "value"
  }
}
`, relatedConfig, name)
}

func TestAccPublicGatewayV3_prepaid(t *testing.T) {
	var (
		obj           interface{}
		rName         = "huaweicloud_natv3_gateway.test"
		name          = acceptance.RandomAccResourceNameWithDash()
		relatedConfig = common.TestBaseNetwork(name)
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPublicGatewayV3ResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPublicGatewayV3_prepaid_step_1(name, relatedConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "spec", "1"),
					resource.TestCheckResourceAttr(rName, "description", "Created by acc test"),
					resource.TestCheckResourceAttr(rName, "ngport_ip_address", "192.168.0.101"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(rName, "session_conf.0.tcp_session_expire_time", "800"),
					resource.TestCheckResourceAttr(rName, "session_conf.0.udp_session_expire_time", "200"),
					resource.TestCheckResourceAttr(rName, "session_conf.0.icmp_session_expire_time", "100"),
					resource.TestCheckResourceAttr(rName, "session_conf.0.tcp_time_wait_time", "50"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttrSet(rName, "billing_info"),
					resource.TestCheckResourceAttrSet(rName, "pps_max"),
					resource.TestCheckResourceAttrSet(rName, "bps_max"),
				),
			},
			{
				Config: testAccPublicGatewayV3_prepaid_step_2(name, relatedConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "spec", "1"),
					resource.TestCheckResourceAttr(rName, "description", "Created by acc test"),
					resource.TestCheckResourceAttr(rName, "ngport_ip_address", "192.168.0.101"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(rName, "session_conf.0.tcp_session_expire_time", "900"),
					resource.TestCheckResourceAttr(rName, "session_conf.0.udp_session_expire_time", "400"),
					resource.TestCheckResourceAttr(rName, "session_conf.0.icmp_session_expire_time", "30"),
					resource.TestCheckResourceAttr(rName, "session_conf.0.tcp_time_wait_time", "0"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "baaar"),
					resource.TestCheckResourceAttr(rName, "tags.newkey", "value"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"charging_mode",
					"period_unit",
					"period",
					"auto_renew",
				},
			},
		},
	})
}

func testAccPublicGatewayV3_prepaid_step_1(name, relatedConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_natv3_gateway" "test" {
  name                  = "%[2]s"
  spec                  = "1"
  description           = "Created by acc test"
  ngport_ip_address     = "192.168.0.101"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  enterprise_project_id = "0"

  session_conf {
    tcp_session_expire_time  = 800
    udp_session_expire_time  = 200
    icmp_session_expire_time = 100
    tcp_time_wait_time       = 50
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = "1"
  auto_renew    = "true"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, relatedConfig, name)
}

func testAccPublicGatewayV3_prepaid_step_2(name, relatedConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_natv3_gateway" "test" {
  name                  = "%[2]s"
  spec                  = "1"
  description           = "Created by acc test"
  ngport_ip_address     = "192.168.0.101"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  enterprise_project_id = "0"

  session_conf {
    tcp_session_expire_time  = 900
    udp_session_expire_time  = 400
    icmp_session_expire_time = 30
    tcp_time_wait_time       = 0
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = "1"
  auto_renew    = "false"

  tags = {
    foo    = "baaar"
    newkey = "value"
  }
}
`, relatedConfig, name)
}
