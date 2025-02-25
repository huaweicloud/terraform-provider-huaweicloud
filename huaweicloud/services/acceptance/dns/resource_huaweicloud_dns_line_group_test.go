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

func getLineGroup(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("dns", region)
	if err != nil {
		return nil, fmt.Errorf("error creating DNS client: %s", err)
	}

	return dns.GetLineGroupById(client, state.Primary.ID)
}

func TestAccLineGroup_basic(t *testing.T) {
	var (
		lineGroup interface{}
		rName     = "huaweicloud_dns_line_group.test"
		rc        = acceptance.InitResourceCheck(rName, &lineGroup, getLineGroup)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLineGroup_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "lines.#", "2"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccLineGroup_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "lines.#", "3"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
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

func testAccLineGroup_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dns_line_group" "test" {
  name        = "%s"
  description = "test description"
  lines       = ["Dianxin_Tianjin", "Dianxin_Jilin"]
}
`, name)
}

func testAccLineGroup_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dns_line_group" "test" {
  name  = "%s_update"
  lines = ["Dianxin_Beijing", "Dianxin_Jilin", "Dianxin_Zhejiang"]
}
`, name)
}
