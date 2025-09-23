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

func getPtrRecord(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.DnsV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DNS client : %s", err)
	}
	return ptrrecords.Get(client, state.Primary.ID).Extract()
}

func TestAccPtrRecord_basic(t *testing.T) {
	var (
		ptrRecord    interface{}
		resourceName = "huaweicloud_dns_ptrrecord.test"
		rc           = acceptance.InitResourceCheck(resourceName, &ptrRecord, getPtrRecord)

		name              = fmt.Sprintf("acpttest-ptr-%s.com", acctest.RandString(5))
		nameWithDotSuffix = fmt.Sprintf("%s.", name)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPtrRecord_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", nameWithDotSuffix),
					resource.TestCheckResourceAttr(resourceName, "description", "a ptr record"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "6000"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttrPair(resourceName, "floatingip_id", "huaweicloud_vpc_eip.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "address"),
				),
			},
			{
				Config: testAccPtrRecord_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("update-%s", nameWithDotSuffix)),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
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

func TestAccPtrRecord_withEpsId(t *testing.T) {
	var (
		ptrRecord    interface{}
		resourceName = "huaweicloud_dns_ptrrecord.test"
		rc           = acceptance.InitResourceCheck(resourceName, &ptrRecord, getPtrRecord)

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
				Config: testAccPtrRecord_withEpsId(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", nameWithDotSuffix),
					resource.TestCheckResourceAttr(resourceName, "description", "a ptr record"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "300"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
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

func testAccPtrRecord_base() string {
	return `
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
}
`
}

func testAccPtrRecord_basic_step1(ptrName string) string {
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
`, testAccPtrRecord_base(), ptrName)
}

func testAccPtrRecord_basic_step2(ptrName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dns_ptrrecord" "test" {
  name          = "update-%s"
  floatingip_id = huaweicloud_vpc_eip.test.id
  ttl           = 7000

  tags = {
    foo = "bar"
  }
}
`, testAccPtrRecord_base(), ptrName)
}

func testAccPtrRecord_withEpsId(ptrName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dns_ptrrecord" "test" {
  name                  = "%s"
  description           = "a ptr record"
  floatingip_id         = huaweicloud_vpc_eip.test.id
  enterprise_project_id = "%s"
}
`, testAccPtrRecord_base(), ptrName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
