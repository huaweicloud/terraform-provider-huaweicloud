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

func getPublicSnatRuleResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("nat", region)
	if err != nil {
		return nil, fmt.Errorf("error creating NAT v2 client: %s", err)
	}

	return nat.GetSnatRule(client, state.Primary.ID)
}

func TestAccPublicSnatRule_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_nat_snat_rule.test"
		name  = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPublicSnatRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPublicSnatRule_basic_step_1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "nat_gateway_id", "huaweicloud_nat_gateway.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "floating_ip_id", "huaweicloud_vpc_eip.test.0", "id"),
					resource.TestCheckResourceAttr(rName, "description", "Created by acc test"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(rName, "freezed_ip_address", ""),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
				),
			},
			{
				Config: testAccPublicSnatRule_basic_step_2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "nat_gateway_id", "huaweicloud_nat_gateway.test", "id"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
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

func testAccPublicSnatRule_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpc_eip" "test" {
  count = 2

  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = format("%[2]s-%%d", count.index)
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_nat_gateway" "test" {
  name                  = "%[2]s"
  spec                  = "2"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  enterprise_project_id = "0"
}
`, common.TestBaseNetwork(name), name)
}

func testAccPublicSnatRule_basic_step_1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_nat_snat_rule" "test" {
  nat_gateway_id = huaweicloud_nat_gateway.test.id
  subnet_id      = huaweicloud_vpc_subnet.test.id
  floating_ip_id = huaweicloud_vpc_eip.test[0].id
  description    = "Created by acc test"
}
`, testAccPublicSnatRule_base(name))
}

func testAccPublicSnatRule_basic_step_2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_nat_snat_rule" "test" {
  nat_gateway_id = huaweicloud_nat_gateway.test.id
  subnet_id      = huaweicloud_vpc_subnet.test.id
  floating_ip_id = join(",", huaweicloud_vpc_eip.test[*].id)
}
`, testAccPublicSnatRule_base(name))
}

func TestAccPublicSnatRule_associatedGlobalEIP(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_nat_snat_rule.test"
		name  = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPublicSnatRuleResourceFunc,
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
				Config: testAccPublicSnatRule_associatedGlobalEIP_step_1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "nat_gateway_id", "huaweicloud_nat_gateway.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "global_eip_id", "huaweicloud_global_eip.test", "id"),
					resource.TestCheckResourceAttr(rName, "description", "Created by terraform"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccPublicSnatRule_associatedGlobalEIP_step_2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "nat_gateway_id", "huaweicloud_nat_gateway.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "global_eip_id", "huaweicloud_global_eip.retest", "id"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
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

func testAccPublicSnatRule_associatedGlobalEIP_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_nat_gateway" "test" {
  name                  = "%[2]s"
  description           = "created by terraform"
  spec                  = "1"
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
`, common.TestBaseNetwork(name), name, acceptance.HW_PROJECT_ID)
}

func testAccPublicSnatRule_associatedGlobalEIP_step_1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_nat_snat_rule" "test" {
  depends_on = [
    huaweicloud_global_eip_associate.test
  ]

  nat_gateway_id = huaweicloud_nat_gateway.test.id
  subnet_id      = huaweicloud_vpc_subnet.test.id
  global_eip_id  = huaweicloud_global_eip.test.id
  description    = "Created by terraform"
}
`, testAccPublicSnatRule_associatedGlobalEIP_base(name))
}

func testAccPublicSnatRule_associatedGlobalEIP_step_2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_nat_snat_rule" "test" {
  depends_on = [
    huaweicloud_global_eip_associate.retest
  ]

  nat_gateway_id = huaweicloud_nat_gateway.test.id
  subnet_id      = huaweicloud_vpc_subnet.test.id
  global_eip_id  = huaweicloud_global_eip.retest.id
  description    = "Created by terraform"
}
`, testAccPublicSnatRule_associatedGlobalEIP_base(name))
}

func TestAccPublicSnatRule_netWorkId(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_nat_snat_rule.test"
		name  = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPublicSnatRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPublicSnatRule_netWorkId(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "nat_gateway_id", "huaweicloud_nat_gateway.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "network_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "floating_ip_id", "huaweicloud_vpc_eip.test.0", "id"),
					resource.TestCheckResourceAttr(rName, "description", "Created by acc test"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
				),
			},
		},
	})
}

func testAccPublicSnatRule_netWorkId(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_nat_snat_rule" "test" {
  nat_gateway_id = huaweicloud_nat_gateway.test.id
  network_id     = huaweicloud_vpc_subnet.test.id
  floating_ip_id = huaweicloud_vpc_eip.test[0].id
  description    = "Created by acc test"
}
`, testAccPublicSnatRule_base(name))
}
