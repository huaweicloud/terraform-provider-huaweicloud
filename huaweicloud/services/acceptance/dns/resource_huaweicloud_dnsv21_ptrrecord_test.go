package dns

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dns"
)

func getDNSV21PtrRecord(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.NewServiceClient("dns_region", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DNS client : %s", err)
	}
	return dns.GetDNSV21PtrRecord(client, state.Primary.ID)
}

func TestAccDNSV21PtrRecord_basic(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_dnsv21_ptrrecord.test"
		rc           = acceptance.InitResourceCheck(
			resourceName,
			&obj,
			getDNSV21PtrRecord,
		)

		name = fmt.Sprintf("acpttest-ptr-%s.com", acctest.RandString(5))
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV21PtrRecord_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "names.#"),
					resource.TestCheckResourceAttr(resourceName, "description", "a ptr record"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "6000"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttrPair(resourceName, "publicip_id", "huaweicloud_vpc_eip.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "address"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				Config: testAccDNSV21PtrRecord_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "names.#"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "ttl", "7000"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttrPair(resourceName, "publicip_id", "huaweicloud_vpc_eip.test", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "address"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
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

func testAccDNSV21PtrRecord_basic(ptrName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dnsv21_ptrrecord" "test" {
  names       = ["1-%[2]s", "2-%[2]s"]
  description = "a ptr record"
  publicip_id = huaweicloud_vpc_eip.test.id
  ttl         = 6000

  tags = {
    key = "value"
  }
}
`, testAccPtrRecord_base(), ptrName)
}

func testAccDNSV21PtrRecord_update(ptrName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dnsv21_ptrrecord" "test" {
  names       = ["3-%[2]s", "4-%[2]s"]
  publicip_id = huaweicloud_vpc_eip.test.id
  ttl         = 7000

  tags = {
    foo = "bar"
  }
}
`, testAccPtrRecord_base(), ptrName)
}
