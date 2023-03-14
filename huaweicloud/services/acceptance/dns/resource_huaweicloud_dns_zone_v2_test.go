package dns

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dns/v2/zones"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func getDNSZoneResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	dnsClient, err := c.DnsV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmtp.Errorf("Error creating DNS client: %s", err)
	}
	return zones.Get(dnsClient, state.Primary.ID).Extract()
}

func TestAccDNSV2Zone_basic(t *testing.T) {
	var zone zones.Zone
	resourceName := "huaweicloud_dns_zone.zone_1"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&zone,
		getDNSZoneResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDNS(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2Zone_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "zone_type", "public"),
					resource.TestCheckResourceAttr(resourceName, "description", "a zone"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "300"),
					resource.TestCheckResourceAttr(resourceName, "tags.zone_type", "public"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttrSet(resourceName, "email"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDNSV2Zone_update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "an updated zone"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "600"),
					resource.TestCheckResourceAttr(resourceName, "tags.zone_type", "public"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "tf-acc"),
					resource.TestCheckResourceAttrSet(resourceName, "email"),
				),
			},
		},
	})
}

func TestAccDNSV2Zone_private(t *testing.T) {
	var zone zones.Zone
	resourceName := "huaweicloud_dns_zone.zone_1"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&zone,
		getDNSZoneResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDNS(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2Zone_private(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "zone_type", "private"),
					resource.TestCheckResourceAttr(resourceName, "description", "a private zone"),
					resource.TestCheckResourceAttr(resourceName, "email", "email@example.com"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "300"),
					resource.TestCheckResourceAttr(resourceName, "tags.zone_type", "private"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
		},
	})
}

func TestAccDNSV2Zone_readTTL(t *testing.T) {
	var zone zones.Zone
	resourceName := "huaweicloud_dns_zone.zone_1"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&zone,
		getDNSZoneResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDNS(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2Zone_readTTL(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestMatchResourceAttr(resourceName, "ttl", regexp.MustCompile("^[0-9]+$")),
				),
			},
		},
	})
}

func TestAccDNSV2Zone_withEpsId(t *testing.T) {
	var zone zones.Zone
	resourceName := "huaweicloud_dns_zone.zone_1"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&zone,
		getDNSZoneResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDNS(t); acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2Zone_withEpsId(name),
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

func testAccDNSV2Zone_basic(zoneName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dns_zone" "zone_1" {
  name        = "%s"
  description = "a zone"
  ttl         = 300

  tags = {
    zone_type = "public"
    owner     = "terraform"
  }
}
	`, zoneName)
}

func testAccDNSV2Zone_update(zoneName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dns_zone" "zone_1" {
  name        = "%s"
  description = "an updated zone"
  ttl         = 600

  tags = {
    zone_type = "public"
    owner     = "tf-acc"
  }
}
	`, zoneName)
}

func testAccDNSV2Zone_readTTL(zoneName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dns_zone" "zone_1" {
  name  = "%s"
  email = "email1@example.com"
}
	`, zoneName)
}

func testAccDNSV2Zone_private(zoneName string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc" "default" {
  name = "vpc-default"
}

resource "huaweicloud_dns_zone" "zone_1" {
  name        = "%s"
  email       = "email@example.com"
  description = "a private zone"
  zone_type   = "private"

  router {
    router_id = data.huaweicloud_vpc.default.id
  }
  tags = {
    zone_type = "private"
    owner     = "terraform"
  }
}
	`, zoneName)
}

func testAccDNSV2Zone_withEpsId(zoneName string) string {
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
