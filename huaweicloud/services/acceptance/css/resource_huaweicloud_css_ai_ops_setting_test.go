package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/css"
)

func getResourceAiOpsSettingFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("css", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CSS client: %s", err)
	}

	return css.GetAiOpsSettingInfo(client, state.Primary.ID)
}

func TestAccResourceAiOpsSetting_basic(t *testing.T) {
	var (
		rName  = "huaweicloud_css_ai_ops_setting.test"
		object interface{}
		rc     = acceptance.InitResourceCheck(
			rName,
			&object,
			getResourceAiOpsSettingFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCSSClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAiOpsSetting_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "cluster_id", acceptance.HW_CSS_CLUSTER_ID),
					resource.TestCheckResourceAttr(rName, "check_type", "full_detection"),
					resource.TestCheckResourceAttr(rName, "period", "15:00 GMT+08:00"),
					resource.TestCheckResourceAttrSet(rName, "check_items.#"),
				),
			},
			{
				Config: testAccAiOpsSetting_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "check_type", "partial_detection"),
					resource.TestCheckResourceAttr(rName, "period", "16:00 GMT+08:00"),
					resource.TestCheckResourceAttr(rName, "check_items.#", "2"),
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

func testAccAiOpsSetting_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_css_ai_ops_setting" "test" {
  cluster_id = "%s"
  check_type = "full_detection"
  period     = "15:00 GMT+08:00"
}
`, acceptance.HW_CSS_CLUSTER_ID)
}

func testAccAiOpsSetting_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_css_ai_ops_setting" "test" {
  cluster_id  = "%s"
  check_type  = "partial_detection"
  period      = "16:00 GMT+08:00"
  check_items = ["nodeDiskCheck","dataNodeDiskBalanceCheck"]
}
`, acceptance.HW_CSS_CLUSTER_ID)
}
