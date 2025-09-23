package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/waf"
)

func getResourceIpIntelligenceRuleFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("waf", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating WAF client: %s", err)
	}

	return waf.GetIpIntelligenceRuleInfo(client, state.Primary.Attributes["policy_id"], state.Primary.ID)
}

func TestAccResourceIpIntelligenceRule_basic(t *testing.T) {
	var (
		rName = "huaweicloud_waf_ip_intelligence_rule.test"
		name1 = acceptance.RandomAccResourceName()
		name2 = acceptance.RandomAccResourceName()

		object interface{}
		rc     = acceptance.InitResourceCheck(
			rName,
			&object,
			getResourceIpIntelligenceRuleFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running test, prepare a WAF professional edition cloud instance or a WAF dedicated instance
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIpIntelligenceRule_basic(name1, name2),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "policy_id", "huaweicloud_waf_policy.test", "id"),
					resource.TestCheckResourceAttr(rName, "type", "idc"),
					resource.TestCheckResourceAttr(rName, "tags.#", "1"),
					resource.TestCheckResourceAttr(rName, "name", name1),
					resource.TestCheckResourceAttr(rName, "policyname", name2),
					resource.TestCheckResourceAttr(rName, "description", "terraform test"),
					resource.TestCheckResourceAttr(rName, "action.0.category", "log"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "policyid"),
					resource.TestCheckResourceAttrSet(rName, "timestamp"),
				),
			},
			{
				Config: testAccIpIntelligenceRule_update(name1, name2),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "tags.#", "2"),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name1)),
					resource.TestCheckResourceAttr(rName, "policyname", fmt.Sprintf("%s_update", name2)),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "action.0.category", "block"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"policyname", "enterprise_project_id"},
				ImportStateIdFunc:       testAccIpIntelligenceRuleImportStateFunc(rName),
			},
		},
	})
}

func testAccIpIntelligenceRule_basic(name1, name2 string) string {
	return fmt.Sprintf(`
resource "huaweicloud_waf_policy" "test" {
  name                  = "%[1]s"
  level                 = 1
  enterprise_project_id = "0"
}

resource "huaweicloud_waf_ip_intelligence_rule" "test" {
  policy_id             = huaweicloud_waf_policy.test.id
  type                  = "idc"
  tags                  = ["HW"]
  name                  = "%[1]s"
  policyname            = "%[2]s"
  description           = "terraform test"
  enterprise_project_id = "0"

  action {
    category = "log"
  }
}
`, name1, name2)
}

func testAccIpIntelligenceRule_update(name1, name2 string) string {
	return fmt.Sprintf(`
resource "huaweicloud_waf_policy" "test" {
  name                  = "%[1]s"
  level                 = 1
  enterprise_project_id = "0"
}

resource "huaweicloud_waf_ip_intelligence_rule" "test" {
  policy_id             = huaweicloud_waf_policy.test.id
  type                  = "idc"
  tags                  = ["HW","Tencent"]
  name                  = "%[1]s_update"
  policyname            = "%[2]s_update"
  description           = ""
  enterprise_project_id = "0"

  action {
    category = "block"
  }
}
`, name1, name2)
}

func testAccIpIntelligenceRuleImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var policyId, ruleId string
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		policyId = rs.Primary.Attributes["policy_id"]
		ruleId = rs.Primary.ID

		if policyId == "" || ruleId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<policy_id>/<id>', but got '%s/%s'",
				policyId, ruleId)
		}

		return fmt.Sprintf("%s/%s", policyId, ruleId), nil
	}
}
