package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/dns/v2/ptrrecords"
)

func randomPtrName() string {
	return fmt.Sprintf("acpttest-%s.com.", acctest.RandString(5))
}

func TestAccDNSV2PtrRecord_basic(t *testing.T) {
	var ptrrecord ptrrecords.Ptr
	ptrName := randomPtrName()
	resourceName := "huaweicloud_dns_ptrrecord.ptr_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDNS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSV2PtrRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2PtrRecord_basic(ptrName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2PtrRecordExists(resourceName, &ptrrecord),
					resource.TestCheckResourceAttr(resourceName, "description", "a ptr record"),
				),
			},
			{
				Config: testAccDNSV2PtrRecord_update(ptrName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2PtrRecordExists(resourceName, &ptrrecord),
					resource.TestCheckResourceAttr(resourceName, "description", "ptr record updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
				),
			},
		},
	})
}

func testAccCheckDNSV2PtrRecordDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	dnsClient, err := config.DnsV2Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud DNS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_dns_ptrrecord" {
			continue
		}

		_, err = ptrrecords.Get(dnsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Ptr record still exists")
		}
	}

	return nil
}

func testAccCheckDNSV2PtrRecordExists(n string, ptrrecord *ptrrecords.Ptr) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		dnsClient, err := config.DnsV2Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud DNS client: %s", err)
		}

		found, err := ptrrecords.Get(dnsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Ptr record not found")
		}

		*ptrrecord = *found

		return nil
	}
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
