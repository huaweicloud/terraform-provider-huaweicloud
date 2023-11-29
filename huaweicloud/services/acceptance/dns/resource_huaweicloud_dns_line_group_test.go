package dns

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getDNSLineGroupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// Query DNS line group.
	var (
		getDNSLineGroupHttpUrl = "v2.1/linegroups/{linegroup_id}"
		getDNSLineGroupProduct = "dns"
	)
	getDNSLineGroupClient, err := cfg.NewServiceClient(getDNSLineGroupProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DNS client: %s", err)
	}

	getDNSLineGroupPath := getDNSLineGroupClient.Endpoint + getDNSLineGroupHttpUrl
	getDNSLineGroupPath = strings.ReplaceAll(getDNSLineGroupPath, "{linegroup_id}", state.Primary.ID)

	getDNSLineGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getDNSLineGroupResp, err := getDNSLineGroupClient.Request("GET", getDNSLineGroupPath, &getDNSLineGroupOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DNS line group: %s", err)
	}

	getDNSLineGroupRespBody, err := utils.FlattenResponse(getDNSLineGroupResp)
	if err != nil {
		return nil, fmt.Errorf("error flatten DNS line group response: %s", err)
	}
	return getDNSLineGroupRespBody, nil
}

func TestAccDNSLineGroup_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dns_line_group.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDNSLineGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDNSLineGroup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "lines.#", "2"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
				),
			},
			{
				Config: testDNSLineGroup_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "description", "test description update"),
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

func testDNSLineGroup_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dns_line_group" "test" {
  name        = "%s"
  description = "test description"
  lines       = ["Dianxin_Tianjin", "Dianxin_Jilin"]
}
`, name)
}

func testDNSLineGroup_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dns_line_group" "test" {
  name        = "%s_update"
  description = "test description update"
  lines       = ["Dianxin_Beijing", "Dianxin_Jilin", "Dianxin_Zhejiang"]
}
`, name)
}
