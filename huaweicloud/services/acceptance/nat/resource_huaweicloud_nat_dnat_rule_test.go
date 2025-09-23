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

func getPublicDnatRuleResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("nat", region)
	if err != nil {
		return nil, fmt.Errorf("error creating NAT v2 client: %s", err)
	}

	return nat.GetDnatRule(client, state.Primary.ID)
}

func TestAccPublicDnatRule_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_nat_dnat_rule.test"
		name  = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPublicDnatRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckProjectID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPublicDnatRule_basic_step_1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "nat_gateway_id", "huaweicloud_nat_gateway.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "floating_ip_id", "huaweicloud_vpc_eip.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "private_ip", "huaweicloud_compute_instance.test", "network.0.fixed_ip_v4"),
					resource.TestCheckResourceAttr(rName, "protocol", "udp"),
					resource.TestCheckResourceAttr(rName, "description", "Created by acc test"),
					resource.TestCheckResourceAttr(rName, "internal_service_port", "80"),
					resource.TestCheckResourceAttr(rName, "external_service_port", "8080"),
				),
			},
			{
				Config: testAccPublicDnatRule_basic_step_2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "protocol", "any"),
					resource.TestCheckResourceAttr(rName, "internal_service_port", "0"),
					resource.TestCheckResourceAttr(rName, "external_service_port", "0"),
				),
			},
			{
				Config: testAccPublicDnatRule_basic_step_3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "protocol", "tcp"),
					resource.TestCheckResourceAttr(rName, "internal_service_port", "0"),
					resource.TestCheckResourceAttr(rName, "external_service_port", "0"),
					resource.TestCheckResourceAttr(rName, "internal_service_port_range", "23-823"),
					resource.TestCheckResourceAttr(rName, "external_service_port_range", "8023-8823"),
				),
			},
			{
				Config: testAccPublicDnatRule_basic_step_4(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "global_eip_id", "huaweicloud_global_eip.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "global_eip_address"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
				),
			},
			{
				Config: testAccPublicDnatRule_basic_step_5(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "global_eip_id", "huaweicloud_global_eip.retest", "id"),
					resource.TestCheckResourceAttr(rName, "protocol", "udp"),
					resource.TestCheckResourceAttr(rName, "internal_service_port", "23"),
					resource.TestCheckResourceAttr(rName, "external_service_port", "8023"),
					resource.TestCheckResourceAttrSet(rName, "global_eip_address"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
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

func TestAccPublicDnatRule_withPort(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_nat_dnat_rule.test"
		name  = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPublicDnatRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPublicDnatRule_withPort_step_1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "nat_gateway_id", "huaweicloud_nat_gateway.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "floating_ip_id", "huaweicloud_vpc_eip.test1", "id"),
					resource.TestCheckResourceAttrPair(rName, "port_id", "huaweicloud_compute_instance.test", "network.0.port"),
					resource.TestCheckResourceAttr(rName, "protocol", "udp"),
					resource.TestCheckResourceAttr(rName, "description", "Created by acc test"),
					resource.TestCheckResourceAttr(rName, "internal_service_port", "80"),
					resource.TestCheckResourceAttr(rName, "external_service_port", "8080"),
				),
			},
			{
				Config: testAccPublicDnatRule_withPort_step_2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "nat_gateway_id", "huaweicloud_nat_gateway.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "floating_ip_id", "huaweicloud_vpc_eip.test2", "id"),
					resource.TestCheckResourceAttrPair(rName, "port_id", "huaweicloud_compute_instance.test", "network.0.port"),
					resource.TestCheckResourceAttr(rName, "protocol", "tcp"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "internal_service_port_range", "23-823"),
					resource.TestCheckResourceAttr(rName, "external_service_port_range", "8023-8823"),
				),
			},
		},
	})
}

func testAccPublicDnatRule_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "%[2]s"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_compute_instance" "test" {
  name               = "%[2]s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_nat_gateway" "test" {
  name                  = "%[2]s"
  spec                  = "2"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  enterprise_project_id = "0"
}

data "huaweicloud_global_eip_pools" "all" {}

resource "huaweicloud_global_internet_bandwidth" "test" {
  access_site           = data.huaweicloud_global_eip_pools.all.geip_pools[0].access_site
  charge_mode           = "95peak_guar"
  enterprise_project_id = "0"
  size                  = 300
  isp                   = data.huaweicloud_global_eip_pools.all.geip_pools[0].isp
  name                  = "%[2]s-b1"
  type                  = data.huaweicloud_global_eip_pools.all.geip_pools[0].allowed_bandwidth_types[0].type
}

resource "huaweicloud_global_eip" "test" {
  access_site           = data.huaweicloud_global_eip_pools.all.geip_pools[0].access_site
  enterprise_project_id = "0"
  geip_pool_name        = data.huaweicloud_global_eip_pools.all.geip_pools[0].name
  internet_bandwidth_id = huaweicloud_global_internet_bandwidth.test.id
  name                  = "%[2]s-g1"

  tags = {
    foo = "bar"
  }
}

resource "huaweicloud_global_eip_associate" "test" {
  global_eip_id  = huaweicloud_global_eip.test.id
  is_reserve_gcb = false

  associate_instance {
    region        = huaweicloud_nat_gateway.test.region
    project_id    = "%[3]s"
    instance_type = "NATGW"
    instance_id   = huaweicloud_nat_gateway.test.id
  }

  gc_bandwidth {
    name        = "%[2]s-gc1"
    charge_mode = "bwd"
    size        = 5
  }
}

resource "huaweicloud_global_internet_bandwidth" "retest" {
  access_site           = data.huaweicloud_global_eip_pools.all.geip_pools[0].access_site
  charge_mode           = "95peak_guar"
  enterprise_project_id = "0"
  size                  = 300
  isp                   = data.huaweicloud_global_eip_pools.all.geip_pools[0].isp
  name                  = "%[2]s-b2"
  type                  = data.huaweicloud_global_eip_pools.all.geip_pools[0].allowed_bandwidth_types[0].type
}

resource "huaweicloud_global_eip" "retest" {
  access_site           = data.huaweicloud_global_eip_pools.all.geip_pools[0].access_site
  enterprise_project_id = "0"
  geip_pool_name        = data.huaweicloud_global_eip_pools.all.geip_pools[0].name
  internet_bandwidth_id = huaweicloud_global_internet_bandwidth.retest.id
  name                  = "%[2]s-g2"

  tags = {
    foo = "bar"
  }
}

resource "huaweicloud_global_eip_associate" "retest" {
  global_eip_id  = huaweicloud_global_eip.retest.id
  is_reserve_gcb = false

  associate_instance {
    region        = huaweicloud_nat_gateway.test.region
    project_id    = "%[3]s"
    instance_type = "NATGW"
    instance_id   = huaweicloud_nat_gateway.test.id
  }

  gc_bandwidth {
    name        = "%[2]s-gc2"
    charge_mode = "bwd"
    size        = 5
  }
}
`, common.TestBaseComputeResources(name), name, acceptance.HW_PROJECT_ID)
}

func testAccPublicDnatRule_basic_step_1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_nat_dnat_rule" "test" {
  nat_gateway_id        = huaweicloud_nat_gateway.test.id
  floating_ip_id        = huaweicloud_vpc_eip.test.id
  private_ip            = huaweicloud_compute_instance.test.network[0].fixed_ip_v4
  description           = "Created by acc test"
  protocol              = "udp"
  internal_service_port = 80
  external_service_port = 8080
}
`, testAccPublicDnatRule_base(name))
}

func testAccPublicDnatRule_basic_step_2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_nat_dnat_rule" "test" {
  nat_gateway_id        = huaweicloud_nat_gateway.test.id
  floating_ip_id        = huaweicloud_vpc_eip.test.id
  private_ip            = huaweicloud_compute_instance.test.network[0].fixed_ip_v4
  description           = ""
  protocol              = "any"
  internal_service_port = 0
  external_service_port = 0
}
`, testAccPublicDnatRule_base(name))
}

func testAccPublicDnatRule_basic_step_3(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_nat_dnat_rule" "test" {
  nat_gateway_id              = huaweicloud_nat_gateway.test.id
  floating_ip_id              = huaweicloud_vpc_eip.test.id
  private_ip                  = huaweicloud_compute_instance.test.network[0].fixed_ip_v4
  protocol                    = "tcp"
  internal_service_port_range = "23-823"
  external_service_port_range = "8023-8823"
}
`, testAccPublicDnatRule_base(name))
}

func testAccPublicDnatRule_basic_step_4(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_nat_dnat_rule" "test" {
  depends_on = [
    huaweicloud_global_eip_associate.test
  ]

  nat_gateway_id              = huaweicloud_nat_gateway.test.id
  global_eip_id               = huaweicloud_global_eip.test.id
  private_ip                  = huaweicloud_compute_instance.test.network[0].fixed_ip_v4
  protocol                    = "tcp"
  internal_service_port_range = "23-823"
  external_service_port_range = "8023-8823"
}
`, testAccPublicDnatRule_base(name))
}

func testAccPublicDnatRule_basic_step_5(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_nat_dnat_rule" "test" {
  depends_on = [
    huaweicloud_global_eip_associate.retest
  ]

  nat_gateway_id        = huaweicloud_nat_gateway.test.id
  global_eip_id         = huaweicloud_global_eip.retest.id
  private_ip            = huaweicloud_compute_instance.test.network[0].fixed_ip_v4
  protocol              = "udp"
  internal_service_port = 23
  external_service_port = 8023
}
`, testAccPublicDnatRule_base(name))
}

func testAccPublicDnatRule_withPort_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpc_eip" "test1" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "%[2]s-bw1"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_vpc_eip" "test2" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "%[2]s-bw2"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_compute_instance" "test" {
  name               = "%[2]s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_nat_gateway" "test" {
  name                  = "%[2]s"
  spec                  = "2"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  enterprise_project_id = "0"
}
`, common.TestBaseComputeResources(name), name)
}

func testAccPublicDnatRule_withPort_step_1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_nat_dnat_rule" "test" {
  nat_gateway_id        = huaweicloud_nat_gateway.test.id
  floating_ip_id        = huaweicloud_vpc_eip.test1.id
  port_id               = huaweicloud_compute_instance.test.network[0].port
  description           = "Created by acc test"
  protocol              = "udp"
  internal_service_port = 80
  external_service_port = 8080
}
`, testAccPublicDnatRule_withPort_base(name))
}

func testAccPublicDnatRule_withPort_step_2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_nat_dnat_rule" "test" {
  nat_gateway_id              = huaweicloud_nat_gateway.test.id
  floating_ip_id              = huaweicloud_vpc_eip.test2.id
  port_id                     = huaweicloud_compute_instance.test.network[0].port
  protocol                    = "tcp"
  internal_service_port_range = "23-823"
  external_service_port_range = "8023-8823"
}
`, testAccPublicDnatRule_withPort_base(name))
}
