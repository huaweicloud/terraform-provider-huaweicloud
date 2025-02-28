package dns

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dns"
)

func getCustomLine(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	getDNSCustomLineClient, err := cfg.NewServiceClient("dns", acceptance.HW_REGION_NAME)
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
