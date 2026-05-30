package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/hss"
)

func getCustomRuleResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("hss", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HSS client: %s", err)
	}

	return hss.GetCustomRuleByList(client, state.Primary.ID)
}

func TestAccCustomRule_basic(t *testing.T) {
	var (
		host       interface{}
		rName      = "huaweicloud_hss_custom_rule.test"
		name       = acceptance.RandomAccResourceName()
		agentId1   = "af08124fd77581dcd2a5f7cdaa208c285af0a1480f7ca9b5258193c021ca2637" // mock data
		agentId2   = "af08124fd77581dcd2a5f7cdaa208c285af0a1480f7ca9b5258193c021ca2638" // mock data
		ruleValue1 = "08a7baa28dd268f8a12bc1f6fd95869321fe51144c5bf3321a6f6305edcd5245" // mock data
		ruleValue2 = "08a7baa28dd268f8a12bc1f6fd95869321fe51144c5bf3321a6f6305edcd5246" // mock data
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&host,
		getCustomRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCustomRule_basic(name, agentId1, ruleValue1),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "rule_name", name),
					resource.TestCheckResourceAttr(rName, "is_all_host", "true"),
					resource.TestCheckResourceAttr(rName, "custom_rule_value_info.0.auto_block", "1"),
					resource.TestCheckResourceAttr(rName, "custom_rule_value_info.0.hash_type", "sha1"),
					resource.TestCheckResourceAttr(rName, "custom_rule_value_info.0.rule_type", "black_hash"),
					resource.TestCheckResourceAttr(rName, "custom_rule_value_info.0.rule_values.0", ruleValue1),
					resource.TestCheckResourceAttr(rName, "rule_status", "0"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
					resource.TestCheckResourceAttrSet(rName, "update_time"),
					resource.TestCheckResourceAttrSet(rName, "host_num"),
				),
			},
			{
				Config: testAccCustomRule_basic_update(name, agentId2, ruleValue2),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "rule_name", name),
					resource.TestCheckResourceAttr(rName, "is_all_host", "false"),
					resource.TestCheckResourceAttr(rName, "custom_rule_value_info.0.auto_block", "0"),
					resource.TestCheckResourceAttr(rName, "custom_rule_value_info.0.hash_type", "md5"),
					resource.TestCheckResourceAttr(rName, "custom_rule_value_info.0.rule_type", "black_hash"),
					resource.TestCheckResourceAttr(rName, "custom_rule_value_info.0.rule_values.0", ruleValue2),
					resource.TestCheckResourceAttr(rName, "rule_status", "1"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"agent_ids"},
			},
		},
	})
}

func testAccCustomRule_basic(name, agentId, ruleValue string) string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_custom_rule" "test" {
  rule_name   = "%[1]s"
  is_all_host = true
  agent_ids   = ["%[2]s"]
  rule_status = 0

  custom_rule_value_info {
    auto_block  = 1
    hash_type   = "sha1"
    rule_type   = "black_hash"
    rule_values = ["%[3]s"]
  }
}
`, name, agentId, ruleValue)
}

func testAccCustomRule_basic_update(name, agentId, ruleValue string) string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_custom_rule" "test" {
  rule_name   = "%[1]s"
  is_all_host = false
  agent_ids   = ["%[2]s"]
  rule_status = 1

  custom_rule_value_info {
    auto_block  = 0
    hash_type   = "md5"
    rule_type   = "black_hash"
    rule_values = ["%[3]s"]
  }
}
`, name, agentId, ruleValue)
}
