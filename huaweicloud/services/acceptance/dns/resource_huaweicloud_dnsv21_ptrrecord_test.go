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

func getV21PtrRecord(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	if state.Primary.Attributes["regional"] == "true" {
		c.RegionClient = true
	}

	client, err := c.NewServiceClient("dns", state.Primary.Attributes["region"])
	if err != nil {
		return nil, fmt.Errorf("error creating DNS client : %s", err)
	}
	return dns.GetDNSV21PtrRecord(client, state.Primary.ID)
}

// Some regions require the region parameter can not be specified, such as 'sa-brazil-1'.
func TestAccV21PtrRecord_basic(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_dnsv21_ptrrecord.test"
		rc           = acceptance.InitResourceCheck(
			resourceName,
			&obj,
			getV21PtrRecord,
		)

		name = fmt.Sprintf("acpttest-ptr-%s.com", acctest.RandString(5))
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV21PtrRecord_basic(name),
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
				Config: testAccV21PtrRecord_update(name),
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
				// Only ignore regional in acceptance test, it will not be changed in actual use.
				ImportStateVerifyIgnore: []string{"regional"},
			},
		},
	})
}

func testAccV21PtrRecord_basic(ptrName string) string {
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

func testAccV21PtrRecord_update(ptrName string) string {
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

// In Chinese website, the region can be specified or not, such as: `cn-north-4`
func TestAccV21PtrRecord_regional(t *testing.T) {
	var (
		name       = acceptance.RandomAccResourceNameWithDash()
		randString = acctest.RandString(5)

		obj   interface{}
		rName = "huaweicloud_dnsv21_ptrrecord.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getV21PtrRecord)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckCustomRegion(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV21PtrRecord_regional_step1(name, randString),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "names.#", "1"),
					resource.TestCheckResourceAttr(rName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(rName, "ttl", "6000"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttrPair(rName, "publicip_id", "huaweicloud_vpc_eip.test", "id"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttrSet(rName, "address"),
					resource.TestCheckResourceAttrSet(rName, "status"),
				),
			},
			{
				Config: testAccV21PtrRecord_regional_step2(name, randString),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "names.#"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "ttl", "7000"),
					resource.TestCheckResourceAttr(rName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttrSet(rName, "address"),
					resource.TestCheckResourceAttrSet(rName, "status"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				// Only ignore regional in acceptance test, it will not be changed in actual use.
				ImportStateVerifyIgnore: []string{"regional"},
				ImportStateIdFunc:       testAccV21PtrRecordImportStateFunc(rName),
			},
		},
	})
}

func testAccV21PtrRecordImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		region := rs.Primary.Attributes["region"]
		ptrRecordId := rs.Primary.ID
		if region == "" || ptrRecordId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<region>/<id>', but got '%s/%s'", region, ptrRecordId)
		}

		return fmt.Sprintf("%s/%s", region, ptrRecordId), nil
	}
}

func testAccV21PtrRecord_regional_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_eip" "test" {
  region = "%[1]s"

  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "%[2]s"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}
`, acceptance.HW_CUSTOM_REGION_NAME, name)
}

func testAccV21PtrRecord_regional_step1(name, randString string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dnsv21_ptrrecord" "test" {
  region = "%[2]s"

  names                 = ["%[3]s-%[4]s.com"]
  description           = "Created by terraform script"
  publicip_id           = huaweicloud_vpc_eip.test.id
  ttl                   = 6000
  enterprise_project_id = "%[5]s"

  tags = {
    foo = "bar"
  }
}
`, testAccV21PtrRecord_regional_base(name), acceptance.HW_CUSTOM_REGION_NAME, name, randString, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccV21PtrRecord_regional_step2(name, randString string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dnsv21_ptrrecord" "test" {
  region = "%[2]s"

  names                 = ["%[3]s-%[4]s-update.com"]
  publicip_id           = huaweicloud_vpc_eip.test.id
  ttl                   = 7000
  enterprise_project_id = "%[5]s"

  tags = {
    owner = "terraform"
  }
}
`, testAccV21PtrRecord_regional_base(name), acceptance.HW_CUSTOM_REGION_NAME, name, randString, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
