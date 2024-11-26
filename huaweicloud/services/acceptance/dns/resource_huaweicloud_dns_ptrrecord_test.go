package dns

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dns/v2/ptrrecords"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getDNSPtrRecordResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.DnsV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DNS client : %s", err)
	}
	return ptrrecords.Get(client, state.Primary.ID).Extract()
}

func TestAccDNSPtrRecord_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_dns_ptrrecord.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getDNSPtrRecordResourceFunc)

		name              = fmt.Sprintf("acpttest-ptr-%s.com", acctest.RandString(5))
		nameWithDotSuffix = fmt.Sprintf("%s.", name)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDNSPtrRecord_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", nameWithDotSuffix),
					resource.TestCheckResourceAttr(resourceName, "description", "a ptr record"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "6000"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttrPair(resourceName, "floatingip_id", "huaweicloud_vpc_eip.test", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "address"),
				),
			},
			{
				Config: testAccDNSPtrRecord_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("update-%s", nameWithDotSuffix)),
					resource.TestCheckResourceAttr(resourceName, "description", "ptr record updated"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "7000"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttrPair(resourceName, "floatingip_id", "huaweicloud_vpc_eip.test", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "address"),
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

func TestAccDNSPtrRecord_withEpsId(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_dns_ptrrecord.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getDNSPtrRecordResourceFunc)

		name              = fmt.Sprintf("acpttest-ptr-%s.com", acctest.RandString(5))
		nameWithDotSuffix = fmt.Sprintf("%s.", name)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDNSPtrRecord_withEpsId(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", nameWithDotSuffix),
					resource.TestCheckResourceAttr(resourceName, "description", "a ptr record"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "6000"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttrPair(resourceName, "floatingip_id", "huaweicloud_vpc_eip.test", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "address"),
				),
			},
		},
	})
}

func testAccDNSPtrRecord_base() string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "test"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}`)
}

func testAccDNSPtrRecord_basic(ptrName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dns_ptrrecord" "test" {
  name          = "%s"
  description   = "a ptr record"
  floatingip_id = huaweicloud_vpc_eip.test.id
  ttl           = 6000

  tags = {
    key = "value"
  }
}
`, testAccDNSPtrRecord_base(), ptrName)
}

func testAccDNSPtrRecord_update(ptrName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dns_ptrrecord" "test" {
  name          = "update-%s"
  description   = "ptr record updated"
  floatingip_id = huaweicloud_vpc_eip.test.id
  ttl           = 7000

  tags = {
    foo = "bar"
  }
}
`, testAccDNSPtrRecord_base(), ptrName)
}

func testAccDNSPtrRecord_withEpsId(ptrName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dns_ptrrecord" "test" {
  name                  = "%s"
  description           = "a ptr record"
  floatingip_id         = huaweicloud_vpc_eip.test.id
  ttl                   = 6000
  enterprise_project_id = "%s"
}
`, testAccDNSPtrRecord_base(), ptrName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
