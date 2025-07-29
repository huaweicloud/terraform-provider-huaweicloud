package workspace

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/workspace"
)

func getResourceAppRuleRestrictionFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("workspace", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Workspace client: %s", err)
	}

	countStr, ok := state.Primary.Attributes["rule_ids.#"]
	if !ok {
		return nil, errors.New("error retrieving the slice count of the rule IDs")
	}

	count, _ := strconv.Atoi(countStr)
	ruleIds := make([]string, 0)
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("rule_ids.%d", i)
		if ruleId, ok := state.Primary.Attributes[key]; ok {
			ruleIds = append(ruleIds, ruleId)
		}
	}

	return workspace.FilterRestrictedApplicationRuleIds(client, ruleIds)
}

func TestAccAppRuleRestriction_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		resourceName       = "huaweicloud_workspace_app_rule_restriction.test"
		appRuleRestriction interface{}
		rc                 = acceptance.InitResourceCheck(resourceName, &appRuleRestriction, getResourceAppRuleRestrictionFunc)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAppRuleRestriction_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestMatchResourceAttr(resourceName, "rule_ids.#", regexp.MustCompile(`^[0-9]+$`)),
				),
			},
			{
				Config: testAccAppRuleRestriction_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "rule_ids.#", "1"),
				),
			},
		},
	})
}

func testAccAppRuleRestriction_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_app_rule" "with_product_rule" {
  name        = "%[1]s"
  description = "Created by terraform script"

  rule {
    scope = "PRODUCT"

    product_rule {
      identify_condition = "process"
      publisher          = "Microsoft Corporation"
      product_name       = "Microsoft Office"
      process_name       = "WINWORD.EXE"
      support_os         = "Windows"
      version            = "1.0"
      product_version    = "2019"
    }
  }
}

resource "huaweicloud_workspace_app_rule" "with_path_rule" {
  name        = "%[1]s_path"
  description = "Created by terraform script for path rule"

  rule {
    scope = "PATH"

    path_rule {
      path = "C:\\Program Files\\Microsoft Office\\root\\Office16\\WINWORD.EXE"
    }
  }
}	
`, name)
}

func testAccAppRuleRestriction_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_rule_restriction" "test" {
  rule_ids = [
    huaweicloud_workspace_app_rule.with_product_rule.id,
    huaweicloud_workspace_app_rule.with_path_rule.id,
  ]
}
`, testAccAppRuleRestriction_base(name))
}

func testAccAppRuleRestriction_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_rule_restriction" "test" {
  rule_ids = [
    huaweicloud_workspace_app_rule.with_path_rule.id,
  ]
}
`, testAccAppRuleRestriction_base(name))
}
