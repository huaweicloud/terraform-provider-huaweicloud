package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/workspace"
)

func getDesktopNameRule(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("workspace", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Workspace client: %s", err)
	}

	return workspace.GetDesktopNameRule(client, state.Primary.ID)
}

func TestAccDesktopNameRule_basic(t *testing.T) {
	var (
		nameRule         interface{}
		resourceName     = "huaweicloud_workspace_desktop_name_rule.test"
		ruleName         = acceptance.RandomAccResourceName()
		updateRuleName   = acceptance.RandomAccResourceName()
		namePrefix       = "$DomainUser$"
		updataNamePrefix = "test$DomainUser$end"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&nameRule,
		getDesktopNameRule,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDesktopNameRule_step1(ruleName, namePrefix),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", ruleName),
					resource.TestCheckResourceAttr(resourceName, "is_contain_user", "true"),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", namePrefix),
					resource.TestCheckResourceAttr(resourceName, "digit_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "start_number", "2"),
					resource.TestCheckResourceAttr(resourceName, "single_domain_user_increment", "1"),
					resource.TestCheckResourceAttr(resourceName, "is_default_policy", "false"),
				),
			},
			{
				Config: testAccDesktopNameRule_step2(updateRuleName, updataNamePrefix),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateRuleName),
					resource.TestCheckResourceAttr(resourceName, "is_contain_user", "true"),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", updataNamePrefix),
					resource.TestCheckResourceAttr(resourceName, "digit_number", "2"),
					resource.TestCheckResourceAttr(resourceName, "start_number", "10"),
					resource.TestCheckResourceAttr(resourceName, "single_domain_user_increment", "0"),
					resource.TestCheckResourceAttr(resourceName, "is_default_policy", "true"),
				),
			},
			{
				Config: testAccDesktopNameRule_step3(updateRuleName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "is_contain_user", "false"),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", "test"),
					resource.TestCheckResourceAttr(resourceName, "digit_number", "3"),
					resource.TestCheckResourceAttr(resourceName, "start_number", "100"),
					resource.TestCheckResourceAttr(resourceName, "single_domain_user_increment", "0"),
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

func testAccDesktopNameRule_step1(rName string, namePrefix string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_desktop_name_rule" "test" {
  name                         = "%[1]s"
  name_prefix                  = "%[2]s"
  digit_number                 = 1
  start_number                 = 2
  single_domain_user_increment = 1
}
`, rName, namePrefix)
}

func testAccDesktopNameRule_step2(rName string, namePrefix string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_desktop_name_rule" "test" {
  name                         = "%[1]s"
  name_prefix                  = "%[2]s"
  digit_number                 = 2
  start_number                 = 10
  single_domain_user_increment = 0
  is_default_policy            = true 
}
`, rName, namePrefix)
}

func testAccDesktopNameRule_step3(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_desktop_name_rule" "test" {
  name                         = "%[1]s"
  name_prefix                  = "test"
  digit_number                 = 3
  start_number                 = 100
  single_domain_user_increment = 0
}
`, rName)
}
