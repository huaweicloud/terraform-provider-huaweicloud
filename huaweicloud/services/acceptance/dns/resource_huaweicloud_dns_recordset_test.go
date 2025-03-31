package dns

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dns/v2/zones"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getDNSRecordsetResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	dnsProduct := "dns"
	if state.Primary.Attributes["zone_type"] != "public" {
		dnsProduct = "dns_region"
	}

	// getDNSRecordset: Query DNS recordset
	getDNSRecordsetClient, err := cfg.NewServiceClient(dnsProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DNS Client: %s", err)
	}

	parts := strings.SplitN(state.Primary.ID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <zone_id>/<recordset_id>")
	}
	zoneID := parts[0]
	recordsetID := parts[1]

	zoneInfo, err := zones.Get(getDNSRecordsetClient, zoneID).Extract()
	if err != nil {
		return "", fmt.Errorf("error getting zone: %s", err)
	}
	zoneType := zoneInfo.ZoneType
	version := "v2.1"
	if zoneType == "private" {
		version = "v2"
	}
	getDNSRecordsetHttpUrl := fmt.Sprintf("%s/zones/{zone_id}/recordsets/{recordset_id}", version)

	getDNSRecordsetPath := getDNSRecordsetClient.Endpoint + getDNSRecordsetHttpUrl
	getDNSRecordsetPath = strings.ReplaceAll(getDNSRecordsetPath, "{zone_id}", zoneID)
	getDNSRecordsetPath = strings.ReplaceAll(getDNSRecordsetPath, "{recordset_id}", recordsetID)

	getDNSRecordsetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getDNSRecordsetResp, err := getDNSRecordsetClient.Request("GET", getDNSRecordsetPath, &getDNSRecordsetOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DNS recordset: %s", err)
	}
	return utils.FlattenResponse(getDNSRecordsetResp)
}

func TestAccDNSRecordset_basic(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_dns_recordset.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getDNSRecordsetResourceFunc)

		name              = fmt.Sprintf("acpttest-recordset-%s.com", acctest.RandString(5))
		nameWithDotSuffix = fmt.Sprintf("%s.", name)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDNSRecordset_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", nameWithDotSuffix),
					resource.TestCheckResourceAttr(rName, "type", "A"),
					resource.TestCheckResourceAttr(rName, "description", "a recordset description"),
					resource.TestCheckResourceAttr(rName, "status", "ENABLE"),
					resource.TestCheckResourceAttr(rName, "ttl", "300"),
					resource.TestCheckResourceAttr(rName, "records.#", "2"),
					resource.TestCheckResourceAttr(rName, "line_id", "Dianxin_Shanxi"),
					resource.TestCheckResourceAttr(rName, "weight", "3"),
					resource.TestCheckResourceAttr(rName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(rName, "tags.key2", "value2"),
					resource.TestCheckResourceAttr(rName, "zone_type", "public"),
					resource.TestCheckResourceAttrSet(rName, "zone_name"),
				),
			},
			{
				Config: testDNSRecordset_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("update.%s", nameWithDotSuffix)),
					resource.TestCheckResourceAttr(rName, "type", "TXT"),
					resource.TestCheckResourceAttr(rName, "description", "a recordset description update"),
					resource.TestCheckResourceAttr(rName, "status", "DISABLE"),
					resource.TestCheckResourceAttr(rName, "ttl", "600"),
					resource.TestCheckResourceAttr(rName, "records.0", "\"test records\""),
					resource.TestCheckResourceAttr(rName, "weight", "5"),
					resource.TestCheckResourceAttr(rName, "tags.key1", "value1_update"),
					resource.TestCheckResourceAttr(rName, "tags.key2", "value2_update"),
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

func TestAccDNSRecordset_publicZone(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_dns_recordset.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getDNSRecordsetResourceFunc)

		name              = fmt.Sprintf("acpttest-recordset-%s.com", acctest.RandString(5))
		nameWithDotSuffix = fmt.Sprintf("%s.", name)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDNSRecordset_publicZone(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", nameWithDotSuffix),
					resource.TestCheckResourceAttr(rName, "type", "A"),
					resource.TestCheckResourceAttr(rName, "description", "a record set"),
					resource.TestCheckResourceAttr(rName, "status", "ENABLE"),
					resource.TestCheckResourceAttr(rName, "ttl", "3000"),
					resource.TestCheckResourceAttr(rName, "records.0", "10.1.0.0"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttr(rName, "zone_type", "public"),
					resource.TestCheckResourceAttrSet(rName, "zone_name"),
				),
			},
			{
				Config: testDNSRecordset_publicZone_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("update.%s", nameWithDotSuffix)),
					resource.TestCheckResourceAttr(rName, "type", "TXT"),
					resource.TestCheckResourceAttr(rName, "description", "an updated record set"),
					resource.TestCheckResourceAttr(rName, "status", "DISABLE"),
					resource.TestCheckResourceAttr(rName, "ttl", "6000"),
					resource.TestCheckResourceAttr(rName, "records.0", "\"test records\""),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value_updated"),
				),
			},
		},
	})
}

func TestAccDNSRecordset_privateZone(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_dns_recordset.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getDNSRecordsetResourceFunc)

		name              = fmt.Sprintf("acpttest-recordset-%s.com", acctest.RandString(5))
		nameWithDotSuffix = fmt.Sprintf("%s.", name)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDNSRecordset_privateZone(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", nameWithDotSuffix),
					resource.TestCheckResourceAttr(rName, "type", "A"),
					resource.TestCheckResourceAttr(rName, "description", "a private record set"),
					resource.TestCheckResourceAttr(rName, "status", "DISABLE"),
					resource.TestCheckResourceAttr(rName, "ttl", "600"),
					resource.TestCheckResourceAttr(rName, "records.0", "10.1.0.3"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar_private"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value_private"),
					resource.TestCheckResourceAttr(rName, "zone_type", "private"),
					resource.TestCheckResourceAttrSet(rName, "zone_name"),
				),
			},
			{
				Config: testDNSRecordset_privateZone_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("update.%s", nameWithDotSuffix)),
					resource.TestCheckResourceAttr(rName, "type", "TXT"),
					resource.TestCheckResourceAttr(rName, "description", "a private record set update"),
					resource.TestCheckResourceAttr(rName, "status", "ENABLE"),
					resource.TestCheckResourceAttr(rName, "ttl", "900"),
					resource.TestCheckResourceAttr(rName, "records.0", "\"test records\""),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar_private_update"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value_private_update"),
				),
			},
			{
				Config:      testDNSRecordset_privateZone_updateWeight(name),
				ExpectError: regexp.MustCompile(`private zone do not support.`),
			},
			{
				Config:      testDNSRecordset_privateZone_updateLineID(name),
				ExpectError: regexp.MustCompile(`private zone do not support.`),
			},
		},
	})
}

func testAccDNSZone_private(zoneName string) string {
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

func testAccRecordset_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dns_zone" "test" {
  name        = "%s"
  description = "Created by terraform script"
  ttl         = 300
}
`, name)
}

func testDNSRecordset_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dns_recordset" "test" {
  zone_id     = huaweicloud_dns_zone.test.id
  name        = "%s"
  type        = "A"
  description = "a recordset description"
  status      = "ENABLE"
  ttl         = 300
  records     = ["10.1.0.0", "9.1.0.0"]
  line_id     = "Dianxin_Shanxi"
  weight      = 3

  tags = {
    key1 = "value1"
    key2 = "value2"
  }
}
`, testAccRecordset_base(name), name)
}

func testDNSRecordset_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dns_recordset" "test" {
  zone_id     = huaweicloud_dns_zone.test.id
  name        = "update.%s"
  type        = "TXT"
  description = "a recordset description update"
  status      = "DISABLE"
  ttl         = 600
  records     = ["\"test records\""]
  line_id     = "Dianxin_Shanxi"
  weight      = 5

  tags = {
    key1 = "value1_update"
    key2 = "value2_update"
  }
}
`, testAccRecordset_base(name), name)
}

func testDNSRecordset_publicZone(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dns_recordset" "test" {
  zone_id     = huaweicloud_dns_zone.test.id
  name        = "%s"
  type        = "A"
  description = "a record set"
  status      = "ENABLE"
  ttl         = 3000
  records     = ["10.1.0.0"]

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccRecordset_base(name), name)
}

func testDNSRecordset_publicZone_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dns_recordset" "test" {
  zone_id     = huaweicloud_dns_zone.test.id
  name        = "update.%s"
  type        = "TXT"
  description = "an updated record set"
  status      = "DISABLE"
  ttl         = 6000
  records     = ["\"test records\""]

  tags = {
    foo = "bar"
    key = "value_updated"
  }
}
`, testAccRecordset_base(name), name)
}

func testDNSRecordset_privateZone(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dns_recordset" "test" {
  zone_id     = huaweicloud_dns_zone.zone_1.id
  name        = "%s"
  type        = "A"
  description = "a private record set"
  status      = "DISABLE"
  ttl         = 600
  records     = ["10.1.0.3"]

  tags = {
    foo = "bar_private"
    key = "value_private"
  }
}
`, testAccDNSZone_private(name), name)
}

func testDNSRecordset_privateZone_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dns_recordset" "test" {
  zone_id     = huaweicloud_dns_zone.zone_1.id
  name        = "update.%s"
  type        = "TXT"
  description = "a private record set update"
  status      = "ENABLE"
  ttl         = 900
  records     = ["\"test records\""]

  tags = {
    foo = "bar_private_update"
    key = "value_private_update"
  }
}
`, testAccDNSZone_private(name), name)
}

func testDNSRecordset_privateZone_updateWeight(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dns_recordset" "test" {
  zone_id     = huaweicloud_dns_zone.zone_1.id
  name        = "update.%s"
  type        = "TXT"
  description = "a private record set update"
  status      = "ENABLE"
  ttl         = 900
  records     = ["\"test records\""]
  weight      = 3

  tags = {
    foo = "bar_private_update"
    key = "value_private_update"
  }
}
`, testAccDNSZone_private(name), name)
}

func testDNSRecordset_privateZone_updateLineID(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dns_recordset" "test" {
  zone_id     = huaweicloud_dns_zone.zone_1.id
  name        = "update.%s"
  type        = "TXT"
  description = "a private record set update"
  status      = "ENABLE"
  ttl         = 900
  records     = ["\"test records\""]
  line_id     = "Dianxin_Shanxi"

  tags = {
    foo = "bar_private_update"
    key = "value_private_update"
  }
}
`, testAccDNSZone_private(name), name)
}
