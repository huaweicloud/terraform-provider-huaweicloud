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

func getCustomLine(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	if state.Primary.Attributes["regional"] == "true" {
		cfg.RegionClient = true
	}

	getDNSCustomLineClient, err := cfg.NewServiceClient("dns", state.Primary.Attributes["region"])
	if err != nil {
		return nil, fmt.Errorf("error creating DNS Client: %s", err)
	}

	return dns.GetCustomLineById(getDNSCustomLineClient, state.Primary.ID)
}

func TestAccCustomLine_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		customLine interface{}
		rName      = "huaweicloud_dns_custom_line.test"
		rc         = acceptance.InitResourceCheck(rName, &customLine, getCustomLine)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCustomLine_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "ip_segments.0", "100.100.100.100-100.100.100.100"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
				),
			},
			{
				Config: testCustomLine_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "ip_segments.#", "2"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				// Only ignore regional in acceptance test, it will not be changed in actual use.
				ImportStateVerifyIgnore: []string{"regional"},
			},
		},
	})
}

func testCustomLine_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dns_custom_line" "test" {
  name        = "%s"
  description = "test description"
  ip_segments = ["100.100.100.100-100.100.100.100"]
}
`, name)
}

func testCustomLine_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dns_custom_line" "test" {
  name        = "%s_update"
  ip_segments = ["100.100.100.102-100.100.100.102", "100.100.100.101-100.100.100.101"]
}
`, name)
}

func TestAccCustomLine_regional(t *testing.T) {
	var (
		name    = acceptance.RandomAccResourceName()
		randInt = acctest.RandIntRange(50, 100)

		obj   interface{}
		rName = "huaweicloud_dns_custom_line.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getCustomLine)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckCustomRegion(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCustomLine_regional_step1(name, randInt),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "ip_segments.#", "1"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccCustomLine_regional_step2(name, randInt),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "ip_segments.#", "1"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				// Only ignore regional in acceptance test, it will not be changed in actual use.
				ImportStateVerifyIgnore: []string{"regional"},
				ImportStateIdFunc:       testAccCustomLineImportStateFunc(rName),
			},
		},
	})
}

func testAccCustomLineImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		region := rs.Primary.Attributes["region"]
		customLineId := rs.Primary.ID
		if region == "" || customLineId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<region>/<id>', but got '%s/%s'", region, customLineId)
		}

		return fmt.Sprintf("%s/%s", region, customLineId), nil
	}
}

func testAccCustomLine_regional_step1(name string, randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_dns_custom_line" "test" {
  region      = "%[1]s"

  name        = "%[2]s"
  description = "test description"
  ip_segments = ["100.100.100.%[3]v-100.100.100.%[3]v"]
}
`, acceptance.HW_CUSTOM_REGION_NAME, name, randInt)
}

func testAccCustomLine_regional_step2(name string, randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_dns_custom_line" "test" {
  region      = "%[1]s"

  name        = "%[2]s_update"
  ip_segments = ["100.100.100.%[3]v-100.100.100.%[3]v"]
}
`, acceptance.HW_CUSTOM_REGION_NAME, name, randInt+1)
}
