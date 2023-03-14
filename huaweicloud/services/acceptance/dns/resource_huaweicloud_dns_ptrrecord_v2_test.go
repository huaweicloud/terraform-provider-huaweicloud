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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func randomPtrName() string {
	return fmt.Sprintf("acpttest-%s.com.", acctest.RandString(5))
}

func TestAccDNSV2PtrRecord_basic(t *testing.T) {
	var ptrrecord ptrrecords.Ptr
	ptrName := randomPtrName()
	resourceName := "huaweicloud_dns_ptrrecord.ptr_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheckDNS(t) },
		Providers:    acceptance.TestAccProviders,
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

func TestAccDNSV2PtrRecord_withEpsId(t *testing.T) {
	var ptrrecord ptrrecords.Ptr
	ptrName := randomPtrName()
	resourceName := "huaweicloud_dns_ptrrecord.ptr_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheckDNS(t); acceptance.TestAccPreCheckEpsID(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: testAccCheckDNSV2PtrRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2PtrRecord_withEpsId(ptrName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2PtrRecordExists(resourceName, &ptrrecord),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testAccCheckDNSV2PtrRecordDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	dnsClient, err := config.DnsV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud DNS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_dns_ptrrecord" {
			continue
		}

		_, err = ptrrecords.Get(dnsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Ptr record still exists")
		}
	}

	return nil
}

func testAccCheckDNSV2PtrRecordExists(n string, ptrrecord *ptrrecords.Ptr) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		dnsClient, err := config.DnsV2Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud DNS client: %s", err)
		}

		found, err := ptrrecords.Get(dnsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("Ptr record not found")
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
