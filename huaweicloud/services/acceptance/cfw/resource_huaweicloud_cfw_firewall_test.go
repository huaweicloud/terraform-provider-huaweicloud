package cfw

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getFirewallResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getFirewall: Query the List of CFW firewalls
	var (
		getFirewallHttpUrl = "v1/{project_id}/firewall/exist"
		getFirewallProduct = "cfw"
	)
	getFirewallClient, err := cfg.NewServiceClient(getFirewallProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CFW client: %s", err)
	}

	getFirewallPath := getFirewallClient.Endpoint + getFirewallHttpUrl
	getFirewallPath = strings.ReplaceAll(getFirewallPath, "{project_id}", getFirewallClient.ProjectID)
	getFirewallPath += fmt.Sprintf("?offset=0&limit=10&service_type=0&fw_instance_id=%s", state.Primary.ID)

	getFirewallOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getFirewallsResp, err := getFirewallClient.Request("GET", getFirewallPath, &getFirewallOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving firewalls: %s", err)
	}

	getFirewallRespBody, err := utils.FlattenResponse(getFirewallsResp)
	if err != nil {
		return nil, err
	}

	jsonPath := fmt.Sprintf("data.records[?fw_instance_id=='%s']|[0]", state.Primary.ID)
	getFirewallRespBody = utils.PathSearch(jsonPath, getFirewallRespBody, nil)
	if getFirewallRespBody == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return getFirewallRespBody, nil
}

func TestAccFirewall_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cfw_firewall.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getFirewallResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testFirewall_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "charging_mode", "postPaid"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttrSet(rName, "engine_type"),
					resource.TestCheckResourceAttrSet(rName, "ha_type"),
					resource.TestCheckResourceAttrSet(rName, "protect_objects.#"),
					resource.TestCheckResourceAttrSet(rName, "service_type"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "support_ipv6"),
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

func TestAccFirewall_prePaid(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cfw_firewall.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getFirewallResourceFunc,
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
				Config: testFirewall_prePaid(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttrSet(rName, "engine_type"),
					resource.TestCheckResourceAttrSet(rName, "ha_type"),
					resource.TestCheckResourceAttrSet(rName, "protect_objects.#"),
					resource.TestCheckResourceAttrSet(rName, "service_type"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "support_ipv6"),
				),
			},
			{
				Config: testFirewall_prePaid_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "tags.k1", "v1"),
					resource.TestCheckResourceAttr(rName, "tags.k2", "v2"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"period_unit", "period", "auto_renew",
				},
			},
		},
	})
}

func TestAccFirewall_eastWest(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cfw_firewall.test"
	bgpAsNum := acctest.RandIntRange(64512, 65534)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getFirewallResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfwEastWestFirewall(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testFirewall_eastWest(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "charging_mode", "postPaid"),
					resource.TestCheckResourceAttr(rName, "east_west_firewall_inspection_cidr", "172.16.1.0/24"),
					resource.TestCheckResourceAttr(rName, "east_west_firewall_mode", "er"),
					resource.TestCheckResourceAttr(rName, "east_west_firewall_status", "0"),
					resource.TestCheckResourceAttrPair(rName, "east_west_firewall_er_id", "huaweicloud_er_instance.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "engine_type"),
					resource.TestCheckResourceAttrSet(rName, "ha_type"),
					resource.TestCheckResourceAttrSet(rName, "protect_objects.#"),
					resource.TestCheckResourceAttrSet(rName, "service_type"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "support_ipv6"),
					resource.TestCheckResourceAttrSet(rName, "east_west_firewall_er_attachment_id"),
				),
			},
			{
				Config: testFirewall_eastWestUpdate(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "charging_mode", "postPaid"),
					resource.TestCheckResourceAttr(rName, "east_west_firewall_inspection_cidr", "172.16.1.0/24"),
					resource.TestCheckResourceAttr(rName, "east_west_firewall_mode", "er"),
					resource.TestCheckResourceAttr(rName, "east_west_firewall_status", "1"),
					resource.TestCheckResourceAttrPair(rName, "east_west_firewall_er_id", "huaweicloud_er_instance.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "engine_type"),
					resource.TestCheckResourceAttrSet(rName, "ha_type"),
					resource.TestCheckResourceAttrSet(rName, "protect_objects.#"),
					resource.TestCheckResourceAttrSet(rName, "service_type"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "support_ipv6"),
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

func TestAccFirewall_eastWestInExistingFirewall(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cfw_firewall.test"
	bgpAsNum := acctest.RandIntRange(64512, 65534)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getFirewallResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfwEastWestFirewall(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testFirewall_eastWestInExistingFirewall(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "charging_mode", "postPaid"),
					resource.TestCheckResourceAttrSet(rName, "engine_type"),
					resource.TestCheckResourceAttrSet(rName, "ha_type"),
					resource.TestCheckResourceAttrSet(rName, "protect_objects.#"),
					resource.TestCheckResourceAttrSet(rName, "service_type"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "support_ipv6"),
				),
			},
			{
				Config: testFirewall_eastWestInExistingFirewallUpdate(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "charging_mode", "postPaid"),
					resource.TestCheckResourceAttr(rName, "east_west_firewall_inspection_cidr", "172.16.1.0/24"),
					resource.TestCheckResourceAttr(rName, "east_west_firewall_mode", "er"),
					resource.TestCheckResourceAttr(rName, "east_west_firewall_status", "0"),
					resource.TestCheckResourceAttrPair(rName, "east_west_firewall_er_id", "huaweicloud_er_instance.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "engine_type"),
					resource.TestCheckResourceAttrSet(rName, "ha_type"),
					resource.TestCheckResourceAttrSet(rName, "protect_objects.#"),
					resource.TestCheckResourceAttrSet(rName, "service_type"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "support_ipv6"),
					resource.TestCheckResourceAttrSet(rName, "east_west_firewall_er_attachment_id"),
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

func TestAccFirewall_ips(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cfw_firewall.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getFirewallResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testFirewall_ips_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(rName, "ips_switch_status", "1"),
					resource.TestCheckResourceAttr(rName, "ips_protection_mode", "1"),
					resource.TestCheckResourceAttrSet(rName, "engine_type"),
					resource.TestCheckResourceAttrSet(rName, "ha_type"),
					resource.TestCheckResourceAttrSet(rName, "protect_objects.#"),
					resource.TestCheckResourceAttrSet(rName, "service_type"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "support_ipv6"),
				),
			},
			{
				Config: testFirewall_ips_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "ips_switch_status", "0"),
					resource.TestCheckResourceAttr(rName, "ips_protection_mode", "2"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"period_unit", "period", "auto_renew",
				},
			},
		},
	})
}

func TestAccFirewall_flavor(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cfw_firewall.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getFirewallResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testFirewall_flavor(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(rName, "flavor.0.version", "Professional"),
					resource.TestCheckResourceAttr(rName, "flavor.0.extend_eip_count", "2"),
					resource.TestCheckResourceAttr(rName, "flavor.0.extend_bandwidth", "5"),
					resource.TestCheckResourceAttr(rName, "flavor.0.extend_vpc_count", "1"),
					resource.TestCheckResourceAttrSet(rName, "engine_type"),
					resource.TestCheckResourceAttrSet(rName, "flavor.0.eip_count"),
					resource.TestCheckResourceAttrSet(rName, "flavor.0.bandwidth"),
					resource.TestCheckResourceAttrSet(rName, "flavor.0.vpc_count"),
					resource.TestCheckResourceAttrSet(rName, "flavor.0.default_eip_count"),
					resource.TestCheckResourceAttrSet(rName, "flavor.0.default_bandwidth"),
					resource.TestCheckResourceAttrSet(rName, "flavor.0.default_vpc_count"),
					resource.TestCheckResourceAttrSet(rName, "flavor.0.total_rule_count"),
					resource.TestCheckResourceAttrSet(rName, "ha_type"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"period_unit", "period", "auto_renew",
				},
			},
		},
	})
}

func TestAccFirewall_attachmentID(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cfw_firewall.test"
	bgpAsNum := acctest.RandIntRange(64512, 65534)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getFirewallResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfwEastWestFirewall(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFirewall_attachmentID(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(rName, "east_west_firewall_inspection_cidr", "172.16.1.0/24"),
					resource.TestCheckResourceAttr(rName, "east_west_firewall_mode", "er"),
					resource.TestCheckResourceAttrPair(rName, "east_west_firewall_er_id", "huaweicloud_er_instance.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "engine_type"),
					resource.TestCheckResourceAttrSet(rName, "ha_type"),
					resource.TestCheckResourceAttrSet(rName, "protect_objects.#"),
					resource.TestCheckResourceAttrSet(rName, "service_type"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "support_ipv6"),
					resource.TestCheckResourceAttrSet(rName, "east_west_firewall_er_attachment_id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"period_unit", "period", "auto_renew",
				},
			},
		},
	})
}

func testFirewall_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cfw_firewall" "test" {
  name = "%s"

  flavor {
    version = "Professional"
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, name)
}

func testFirewall_prePaid(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cfw_firewall" "test" {
  name = "%s"

  flavor {
    version = "Professional"
  }

  tags = {
    key = "value"
    foo = "bar"
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = false
}
`, name)
}

func testFirewall_prePaid_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cfw_firewall" "test" {
  name = "%s"

  flavor {
    version = "Professional"
  }

  tags = {
    k1 = "v1"
    k2 = "v2"
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = false
}
`, name)
}

func testFirewall_eastWestBase(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
data "huaweicloud_er_availability_zones" "test" {}

resource "huaweicloud_er_instance" "test" {
  availability_zones = slice(data.huaweicloud_er_availability_zones.test.names, 0, 1)

  name = "%[1]s"
  asn  = %[2]d

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name, bgpAsNum)
}

func testFirewall_eastWest(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cfw_firewall" "test" {
  name = "%s"

  east_west_firewall_inspection_cidr = "172.16.1.0/24"
  east_west_firewall_er_id           = huaweicloud_er_instance.test.id
  east_west_firewall_mode            = "er"

  flavor {
    version = "Professional"
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, testFirewall_eastWestBase(name, bgpAsNum), name)
}

func testFirewall_eastWestUpdate(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cfw_firewall" "test" {
  name = "%s"

  east_west_firewall_inspection_cidr = "172.16.1.0/24"
  east_west_firewall_er_id           = huaweicloud_er_instance.test.id
  east_west_firewall_mode            = "er"
  east_west_firewall_status          = 1

  flavor {
    version = "Professional"
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, testFirewall_eastWestBase(name, bgpAsNum), name)
}

func testFirewall_eastWestInExistingFirewall(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cfw_firewall" "test" {
  name = "%s"

  flavor {
    version = "Professional"
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, testFirewall_eastWestBase(name, bgpAsNum), name)
}

func testFirewall_eastWestInExistingFirewallUpdate(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cfw_firewall" "test" {
  name = "%s"

  east_west_firewall_inspection_cidr = "172.16.1.0/24"
  east_west_firewall_er_id           = huaweicloud_er_instance.test.id
  east_west_firewall_mode            = "er"
  east_west_firewall_status          = 0

  flavor {
    version = "Professional"
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, testFirewall_eastWestBase(name, bgpAsNum), name)
}

func testFirewall_ips_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cfw_firewall" "test" {
  name = "%s"

  flavor {
    version = "Professional"
  }

  tags = {
    key = "value"
    foo = "bar"
  }

  charging_mode        = "prePaid"
  period_unit          = "month"
  period               = 1
  auto_renew           = false
  ips_switch_status    = 1
  ips_protection_mode  = 1
}
`, name)
}

func testFirewall_ips_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cfw_firewall" "test" {
  name = "%s"

  flavor {
    version = "Professional"
  }

  tags = {
    key = "value"
    foo = "bar"
  }

  charging_mode        = "prePaid"
  period_unit          = "month"
  period               = 1
  auto_renew           = false
  ips_switch_status    = 0
  ips_protection_mode  = 2
}
`, name)
}

func testFirewall_flavor(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cfw_firewall" "test" {
  name = "%s"

  flavor {
    version          = "Professional"
    extend_eip_count = 2
    extend_bandwidth = 5
    extend_vpc_count = 1
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = false
}
`, name)
}

func testAccFirewall_attachmentID(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cfw_firewall" "test" {
  name = "%s"

  east_west_firewall_inspection_cidr = "172.16.1.0/24"
  east_west_firewall_er_id           = huaweicloud_er_instance.test.id
  east_west_firewall_mode            = "er"

  flavor {
    version = "Professional"
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = false
}
`, testFirewall_eastWestBase(name, bgpAsNum), name)
}
