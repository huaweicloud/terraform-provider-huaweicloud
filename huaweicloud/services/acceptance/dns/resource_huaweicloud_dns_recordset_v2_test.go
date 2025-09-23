package dns

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dns/v2/recordsets"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getDNSV2RecordsetResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.DnsV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DNS client: %s", err)
	}

	idParts := strings.Split(state.Primary.ID, "/")
	if len(idParts) != 2 {
		return nil, fmt.Errorf("unable to determine DNS record set ID from raw ID: %s", state.Primary.ID)
	}
	zoneID := idParts[0]
	recordsetID := idParts[1]
	return recordsets.Get(client, zoneID, recordsetID).Extract()
}

func TestAccDNSV2RecordSet_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName      = "huaweicloud_dns_recordset_v2.test"
		rc                = acceptance.InitResourceCheck(resourceName, &obj, getDNSV2RecordsetResourceFunc)
		name              = fmt.Sprintf("acpttest-zone-%s.com", acctest.RandString(5))
		nameWithDotSuffix = fmt.Sprintf("%s.", name)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2RecordSet_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", nameWithDotSuffix),
					resource.TestCheckResourceAttr(resourceName, "description", "a record set"),
					resource.TestCheckResourceAttr(resourceName, "type", "A"),
					resource.TestCheckResourceAttrSet(resourceName, "ttl"),
					resource.TestCheckResourceAttr(resourceName, "records.0", "10.1.0.0"),
				),
			},
			{
				Config: testAccDNSV2RecordSet_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", nameWithDotSuffix),
					resource.TestCheckResourceAttr(resourceName, "description", "a record set"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "3000"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccDNSV2RecordSet_basic_step3(name),
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

func TestAccDNSV2RecordSet_private(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_dns_recordset_v2.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getDNSV2RecordsetResourceFunc)

		name              = fmt.Sprintf("acpttest-recordset-%s.com", acctest.RandString(5))
		nameWithDotSuffix = fmt.Sprintf("%s.", name)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2RecordSet_private(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", nameWithDotSuffix),
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
resource "huaweicloud_dns_zone" "test" {
  name        = "%s"
  email       = "email@example.com"
  description = "a zone for acc test"
  ttl         = 6000
}
`, zoneName)
}

func testAccDNSV2RecordSet_basic_step1(zoneName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dns_recordset_v2" "test" {
  zone_id     = huaweicloud_dns_zone.test.id
  name        = "%s"
  type        = "A"
  description = "a record set"
  records     = ["10.1.0.0"]
}
`, testAccDNSV2RecordSet_base(zoneName), zoneName)
}

func testAccDNSV2RecordSet_basic_step2(zoneName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dns_recordset_v2" "test" {
  zone_id     = huaweicloud_dns_zone.test.id
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

func testAccDNSV2RecordSet_basic_step3(zoneName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dns_recordset_v2" "test" {
  zone_id     = huaweicloud_dns_zone.test.id
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

func testAccDNSV2RecordSet_private(zoneName string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc" "default" {
  name = "vpc-default"
}

resource "huaweicloud_dns_zone" "test" {
  name        = "%s"
  email       = "email@example.com"
  description = "a private zone"
  zone_type   = "private"

  router {
    router_id = data.huaweicloud_vpc.default.id
  }
}

resource "huaweicloud_dns_recordset_v2" "test" {
  zone_id     = huaweicloud_dns_zone.test.id
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
