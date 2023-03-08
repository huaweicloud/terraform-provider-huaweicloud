package huaweicloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dns/v2/recordsets"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func randomZoneName() string {
	return fmt.Sprintf("acpttest-zone-%s.com.", acctest.RandString(5))
}

func TestAccDNSV2RecordSet_basic(t *testing.T) {
	var recordset recordsets.RecordSet
	zoneName := randomZoneName()
	resourceName := "huaweicloud_dns_recordset.recordset_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDNS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSV2RecordSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2RecordSet_basic(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2RecordSetExists(resourceName, &recordset),
					resource.TestCheckResourceAttr(resourceName, "name", zoneName),
					resource.TestCheckResourceAttr(resourceName, "description", "a record set"),
					resource.TestCheckResourceAttr(resourceName, "type", "A"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "3000"),
					resource.TestCheckResourceAttr(resourceName, "records.0", "10.1.0.0"),
				),
			},
			{
				Config: testAccDNSV2RecordSet_tags(zoneName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", zoneName),
					resource.TestCheckResourceAttr(resourceName, "description", "a record set"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "3000"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccDNSV2RecordSet_update(zoneName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "an updated record set"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "6000"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value_updated"),
					resource.TestCheckResourceAttr(resourceName, "records.0", "10.1.0.1"),
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

func TestAccDNSV2RecordSet_readTTL(t *testing.T) {
	var recordset recordsets.RecordSet
	zoneName := randomZoneName()
	resourceName := "huaweicloud_dns_recordset.recordset_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDNS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSV2RecordSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2RecordSet_readTTL(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2RecordSetExists(resourceName, &recordset),
					resource.TestMatchResourceAttr(resourceName, "ttl", regexp.MustCompile("^[0-9]+$")),
					resource.TestCheckResourceAttr(resourceName, "records.0", "10.1.0.2"),
				),
			},
		},
	})
}

func TestAccDNSV2RecordSet_private(t *testing.T) {
	var recordset recordsets.RecordSet
	zoneName := randomZoneName()
	resourceName := "huaweicloud_dns_recordset.recordset_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDNS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSV2RecordSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2RecordSet_private(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2RecordSetExists(resourceName, &recordset),
					resource.TestCheckResourceAttr(resourceName, "name", zoneName),
					resource.TestCheckResourceAttr(resourceName, "description", "a private record set"),
					resource.TestCheckResourceAttr(resourceName, "type", "A"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "3000"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "records.0", "10.1.0.3"),
				),
			},
		},
	})
}

func testAccCheckDNSV2RecordSetDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	dnsClient, err := config.DnsV2Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud DNS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_dns_recordset" {
			continue
		}

		zoneID, recordsetID, err := parseDNSV2RecordSetID(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = recordsets.Get(dnsClient, zoneID, recordsetID).Extract()
		if err == nil {
			return fmtp.Errorf("Record set still exists")
		}
	}

	return nil
}

func testAccCheckDNSV2RecordSetExists(n string, recordset *recordsets.RecordSet) resource.TestCheckFunc {
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

		zoneID, recordsetID, err := parseDNSV2RecordSetID(rs.Primary.ID)
		if err != nil {
			return err
		}

		found, err := recordsets.Get(dnsClient, zoneID, recordsetID).Extract()
		if err != nil {
			return err
		}

		if found.ID != recordsetID {
			return fmtp.Errorf("Record set not found")
		}

		*recordset = *found

		return nil
	}
}

func testAccDNSV2RecordSet_base(zoneName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dns_zone" "zone_1" {
  name        = "%s"
  email       = "email@example.com"
  description = "a zone for acc test"
  ttl         = 6000
}
`, zoneName)
}

func testAccDNSV2RecordSet_basic(zoneName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dns_recordset" "recordset_1" {
  zone_id     = huaweicloud_dns_zone.zone_1.id
  name        = "%s"
  type        = "A"
  description = "a record set"
  ttl         = 3000
  records     = ["10.1.0.0"]
}
`, testAccDNSV2RecordSet_base(zoneName), zoneName)
}

func testAccDNSV2RecordSet_tags(zoneName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dns_recordset" "recordset_1" {
  zone_id     = huaweicloud_dns_zone.zone_1.id
  name        = "%s"
  type        = "A"
  description = "a record set"
  ttl         = 3000
  records     = ["10.1.0.0"]

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccDNSV2RecordSet_base(zoneName), zoneName)
}

func testAccDNSV2RecordSet_update(zoneName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dns_recordset" "recordset_1" {
  zone_id     = huaweicloud_dns_zone.zone_1.id
  name        = "%s"
  type        = "A"
  description = "an updated record set"
  ttl         = 6000
  records     = ["10.1.0.1"]

  tags = {
    foo = "bar"
    key = "value_updated"
  }
}
`, testAccDNSV2RecordSet_base(zoneName), zoneName)
}

func testAccDNSV2RecordSet_readTTL(zoneName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dns_recordset" "recordset_1" {
  zone_id = huaweicloud_dns_zone.zone_1.id
  name    = "%s"
  type    = "A"
  records = ["10.1.0.2"]
}
`, testAccDNSV2RecordSet_base(zoneName), zoneName)
}

func testAccDNSV2RecordSet_private(zoneName string) string {
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
}

resource "huaweicloud_dns_recordset" "recordset_1" {
  zone_id     = huaweicloud_dns_zone.zone_1.id
  name        = "%s"
  type        = "A"
  description = "a private record set"
  ttl         = 3000
  records     = ["10.1.0.3"]

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, zoneName, zoneName)
}
