package dns

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dns/v2/recordsets"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func getDNSRecordsetResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.DnsV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmtp.Errorf("Error creating DNS client: %s", err)
	}

	zoneID, recordsetID, err := parseDNSV2RecordSetID(state.Primary.ID)
	if err != nil {
		return nil, err
	}
	return recordsets.Get(client, zoneID, recordsetID).Extract()
}

func parseDNSV2RecordSetID(id string) (string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) != 2 {
		return "", "", fmtp.Errorf("Unable to determine DNS record set ID from raw ID: %s", id)
	}

	zoneID := idParts[0]
	recordsetID := idParts[1]

	return zoneID, recordsetID, nil
}

func TestAccDNSV2RecordSet_basic(t *testing.T) {
	var recordset recordsets.RecordSet
	resourceName := "huaweicloud_dns_recordset.recordset_1"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&recordset,
		getDNSRecordsetResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDNS(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2RecordSet_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "a record set"),
					resource.TestCheckResourceAttr(resourceName, "type", "A"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "3000"),
					resource.TestCheckResourceAttr(resourceName, "records.0", "10.1.0.0"),
				),
			},
			{
				Config: testAccDNSV2RecordSet_tags(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "a record set"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "3000"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccDNSV2RecordSet_update(name),
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
	resourceName := "huaweicloud_dns_recordset.recordset_1"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&recordset,
		getDNSRecordsetResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDNS(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2RecordSet_readTTL(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestMatchResourceAttr(resourceName, "ttl", regexp.MustCompile("^[0-9]+$")),
					resource.TestCheckResourceAttr(resourceName, "records.0", "10.1.0.2"),
				),
			},
		},
	})
}

func TestAccDNSV2RecordSet_private(t *testing.T) {
	var recordset recordsets.RecordSet
	resourceName := "huaweicloud_dns_recordset.recordset_1"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&recordset,
		getDNSRecordsetResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDNS(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2RecordSet_private(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
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
