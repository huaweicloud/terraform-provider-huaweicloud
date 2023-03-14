package dns

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dns/v2/ptrrecords"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func getDNSPtrRecordResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.DnsV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmtp.Errorf("Error creating DNS client : %s", err)
	}
	return ptrrecords.Get(client, state.Primary.ID).Extract()
}

func TestAccDNSV2PtrRecord_basic(t *testing.T) {
	var ptrrecord ptrrecords.Ptr
	resourceName := "huaweicloud_dns_ptrrecord.ptr_1"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&ptrrecord,
		getDNSPtrRecordResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDNS(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2PtrRecord_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "description", "a ptr record"),
				),
			},
			{
				Config: testAccDNSV2PtrRecord_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "description", "ptr record updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
				),
			},
		},
	})
}

func TestAccDNSV2PtrRecord_withEpsId(t *testing.T) {
	var ptrrecord ptrrecords.Ptr
	resourceName := "huaweicloud_dns_ptrrecord.ptr_1"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&ptrrecord,
		getDNSPtrRecordResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDNS(t); acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2PtrRecord_withEpsId(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testAccDNSV2PtrRecord_basic(ptrName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_eip" "eip_1" {
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

resource "huaweicloud_dns_ptrrecord" "ptr_1" {
  name          = "%s"
  description   = "a ptr record"
  floatingip_id = huaweicloud_vpc_eip.eip_1.id
  ttl           = 6000
}
`, ptrName)
}

func testAccDNSV2PtrRecord_update(ptrName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_eip" "eip_1" {
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

resource "huaweicloud_dns_ptrrecord" "ptr_1" {
  name          = "%s"
  description   = "ptr record updated"
  floatingip_id = huaweicloud_vpc_eip.eip_1.id
  ttl           = 6000

  tags = {
    foo = "bar"
  }
}
`, ptrName)
}

func testAccDNSV2PtrRecord_withEpsId(ptrName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_eip" "eip_1" {
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

resource "huaweicloud_dns_ptrrecord" "ptr_1" {
  name                  = "%s"
  description           = "a ptr record"
  floatingip_id         = huaweicloud_vpc_eip.eip_1.id
  ttl                   = 6000
  enterprise_project_id = "%s"
}
`, ptrName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
