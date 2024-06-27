package dns

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dns/v2/zones"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getDNSZoneResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	dnsProduct := "dns"
	if state.Primary.Attributes["zone_type"] != "public" {
		dnsProduct = "dns_region"
	}

	dnsClient, err := c.NewServiceClient(dnsProduct, acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DNS client: %s", err)
	}
	return zones.Get(dnsClient, state.Primary.ID).Extract()
}

func TestAccDNSZone_basic(t *testing.T) {
	var zone zones.Zone
	resourceName := "huaweicloud_dns_zone.zone_1"
	name := fmt.Sprintf("acpttest-zone-%s.com.", acctest.RandString(5))

	rc := acceptance.InitResourceCheck(
		resourceName,
		&zone,
		getDNSZoneResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDNSZone_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "zone_type", "public"),
					resource.TestCheckResourceAttr(resourceName, "description", "a zone"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "300"),
					resource.TestCheckResourceAttr(resourceName, "tags.zone_type", "public"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttrSet(resourceName, "email"),
					resource.TestCheckResourceAttr(resourceName, "status", "DISABLE"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDNSZone_update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "an updated zone"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "600"),
					resource.TestCheckResourceAttr(resourceName, "tags.zone_type", "public"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "tf-acc"),
					resource.TestCheckResourceAttrSet(resourceName, "email"),
					resource.TestCheckResourceAttr(resourceName, "status", "ENABLE"),
				),
			},
		},
	})
}

func TestAccDNSZone_private(t *testing.T) {
	var zone zones.Zone
	resourceName := "huaweicloud_dns_zone.test"
	name := fmt.Sprintf("acpttest-zone-%s.com.", acctest.RandString(5))
	vpcName := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&zone,
		getDNSZoneResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDNSZone_private_step1(name, vpcName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "zone_type", "private"),
					resource.TestCheckResourceAttr(resourceName, "description", "a private zone"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "300"),
					resource.TestCheckResourceAttr(resourceName, "email", "email@example.com"),
					resource.TestCheckResourceAttr(resourceName, "router.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.zone_type", "private"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
			{
				Config: testAccDNSZone_private_step2(name, vpcName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckOutput("valid_route_id", "true"),
					resource.TestCheckResourceAttr(resourceName, "router.#", "2"),
				),
			},
		},
	})
}

func TestAccDNSZone_readTTL(t *testing.T) {
	var zone zones.Zone
	resourceName := "huaweicloud_dns_zone.zone_1"
	name := fmt.Sprintf("acpttest-zone-%s.com.", acctest.RandString(5))

	rc := acceptance.InitResourceCheck(
		resourceName,
		&zone,
		getDNSZoneResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDNSZone_readTTL(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestMatchResourceAttr(resourceName, "ttl", regexp.MustCompile("^[0-9]+$")),
				),
			},
		},
	})
}

func TestAccDNSZone_withEpsId(t *testing.T) {
	var zone zones.Zone
	resourceName := "huaweicloud_dns_zone.zone_1"
	name := fmt.Sprintf("acpttest-zone-%s.com.", acctest.RandString(5))

	rc := acceptance.InitResourceCheck(
		resourceName,
		&zone,
		getDNSZoneResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDNSZone_withEpsId(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "zone_type", "private"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testAccDNSZone_basic(zoneName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dns_zone" "zone_1" {
  name        = "%s"
  description = "a zone"
  ttl         = 300
  status      = "DISABLE"

  tags = {
    zone_type = "public"
    owner     = "terraform"
  }
}
`, zoneName)
}

func testAccDNSZone_update(zoneName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dns_zone" "zone_1" {
  name        = "%s"
  description = "an updated zone"
  ttl         = 600
  status      = "ENABLE"

  tags = {
    zone_type = "public"
    owner     = "tf-acc"
  }
}
`, zoneName)
}

func testAccDNSZone_readTTL(zoneName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dns_zone" "zone_1" {
  name  = "%s"
  email = "email1@example.com"
}
`, zoneName)
}

func testAccDNSZone_private_step1(zoneName, vpcName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  count = 3
  name  = "%s_${count.index}"
  cidr  = "192.168.0.0/16"
}

resource "huaweicloud_dns_zone" "test" {
  name        = "%s"
  email       = "email@example.com"
  description = "a private zone"
  zone_type   = "private"

  dynamic "router" {
    for_each = slice(huaweicloud_vpc.test[*].id, 0, 2)

    content {
      router_id = router.value
    }
  }

  tags = {
    zone_type = "private"
    owner     = "terraform"
  }
}
`, vpcName, zoneName)
}

func testAccDNSZone_private_step2(zoneName, vpcName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  count = 3
  name  = "%s_${count.index}"
  cidr  = "192.168.0.0/16"
}

resource "huaweicloud_dns_zone" "test" {
  name        = "%s"
  email       = "email@example.com"
  description = "a private zone"
  zone_type   = "private"

  dynamic "router" {
    for_each = slice(huaweicloud_vpc.test[*].id, 1, 3)

    content {
      router_id = router.value
    }
  }

  tags = {
    zone_type = "private"
    owner     = "terraform"
  }
}

locals {
  router_ids = huaweicloud_dns_zone.test.router[*].router_id
}

output "valid_route_id" {
  value = contains(local.router_ids, huaweicloud_vpc.test[1].id) && contains(local.router_ids, huaweicloud_vpc.test[2].id)
}
`, vpcName, zoneName)
}

func testAccDNSZone_withEpsId(zoneName string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc" "default" {
  name = "vpc-default"
}

resource "huaweicloud_dns_zone" "zone_1" {
  name                  = "%s"
  email                 = "email@example.com"
  description           = "a private zone"
  zone_type             = "private"
  enterprise_project_id = "%s"

  router {
    router_id = data.huaweicloud_vpc.default.id
  }
  tags = {
    zone_type = "private"
    owner     = "terraform"
  }
}
`, zoneName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
