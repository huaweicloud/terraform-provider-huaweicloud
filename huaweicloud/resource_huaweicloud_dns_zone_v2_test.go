package huaweicloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dns/v2/zones"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccDNSV2Zone_basic(t *testing.T) {
	var zone zones.Zone
	var zoneName = fmt.Sprintf("acpttest%s.com.", acctest.RandString(5))
	resourceName := "huaweicloud_dns_zone.zone_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDNS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSV2ZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2Zone_basic(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2ZoneExists(resourceName, &zone),
					resource.TestCheckResourceAttr(resourceName, "name", zoneName),
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
				Config: testAccDNSV2Zone_update(zoneName),
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
	var zoneName = fmt.Sprintf("acpttest%s.com.", acctest.RandString(5))
	resourceName := "huaweicloud_dns_zone.zone_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDNS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSV2ZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2Zone_private(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2ZoneExists(resourceName, &zone),
					resource.TestCheckResourceAttr(resourceName, "name", zoneName),
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
	var zoneName = fmt.Sprintf("acpttest%s.com.", acctest.RandString(5))
	resourceName := "huaweicloud_dns_zone.zone_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDNS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSV2ZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2Zone_readTTL(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2ZoneExists(resourceName, &zone),
					resource.TestMatchResourceAttr(resourceName, "ttl", regexp.MustCompile("^[0-9]+$")),
				),
			},
		},
	})
}

func TestAccDNSV2Zone_withEpsId(t *testing.T) {
	var zone zones.Zone
	var zoneName = fmt.Sprintf("acpttest%s.com.", acctest.RandString(5))
	resourceName := "huaweicloud_dns_zone.zone_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDNS(t); testAccPreCheckEpsID(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSV2ZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2Zone_withEpsId(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2ZoneExists(resourceName, &zone),
					resource.TestCheckResourceAttr(resourceName, "name", zoneName),
					resource.TestCheckResourceAttr(resourceName, "zone_type", "private"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testAccCheckDNSV2ZoneDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	dnsClient, err := config.DnsV2Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud DNS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_dns_zone" {
			continue
		}

		_, err := zones.Get(dnsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Zone still exists")
		}
	}

	return nil
}

func testAccCheckDNSV2ZoneExists(n string, zone *zones.Zone) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		dnsClient, err := config.DnsV2Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud DNS client: %s", err)
		}

		found, err := zones.Get(dnsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("Zone not found")
		}

		*zone = *found

		return nil
	}
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
	`, zoneName, HW_ENTERPRISE_PROJECT_ID_TEST)
}
