package das

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccSqlLimitRule_basic(t *testing.T) {
	var (
		rName = "huaweicloud_das_sql_limit_rule.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDasInstanceIds(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlLimitRule_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "sql_type", "SELECT"),
					resource.TestCheckResourceAttr(rName, "pattern", "select~test"),
					resource.TestCheckResourceAttr(rName, "max_concurrency", "100"),
				),
			},
			{
				Config: testAccSqlLimitRule_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "sql_type", "UPDATE"),
					resource.TestCheckResourceAttr(rName, "pattern", "update~test"),
					resource.TestCheckResourceAttr(rName, "max_concurrency", "120"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enable_force_new"},
				ImportStateIdFunc:       testAccSqlLimitRuleImportStateFunc(rName),
			},
		},
	})
}

func testAccSqlLimitRule_base() string {
	return fmt.Sprintf(`
locals {
  instance_ids = split(",", "%[1]s")
}
`, acceptance.HW_DAS_INSTANCE_IDS)
}

func testAccSqlLimitRule_basic() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_das_sql_limit_rule" "test" {
  instance_id     = local.instance_ids[0]
  sql_type        = "SELECT"
  pattern         = "select~test"
  max_concurrency = 100
}
`, testAccSqlLimitRule_base())
}

func testAccSqlLimitRule_basic_update() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_das_sql_limit_rule" "test" {
  instance_id     = local.instance_ids[0]
  sql_type        = "UPDATE"
  pattern         = "update~test"
  max_concurrency = 120

  enable_force_new = "true"
}
`, testAccSqlLimitRule_base())
}

func testAccSqlLimitRuleImportStateFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}

		instanceId := rs.Primary.Attributes["instance_id"]
		if instanceId == "" {
			return "", errors.New("import ID is missing, want '<instance_id>/<rule_id>'")
		}
		ruleId := rs.Primary.ID
		if ruleId == "" {
			return "", errors.New("resource ID is missing")
		}
		return fmt.Sprintf("%s/%s", instanceId, ruleId), nil
	}
}
