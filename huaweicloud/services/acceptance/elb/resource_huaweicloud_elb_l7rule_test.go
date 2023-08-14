package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/elb/v3/l7policies"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getELBl7RuleResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	lbClient, err := cfg.ElbV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ELB client: %s", err)
	}

	l7policyID := state.Primary.Attributes["l7policy_id"]
	return l7policies.GetRule(lbClient, l7policyID, state.Primary.ID).Extract()
}

func TestAccElbV3L7Rule_basic(t *testing.T) {
	var l7rule l7policies.Rule
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_elb_l7rule.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&l7rule,
		getELBl7RuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckElbV3L7RuleConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "PATH"),
					resource.TestCheckResourceAttr(resourceName, "compare_type", "EQUAL_TO"),
					resource.TestCheckResourceAttr(resourceName, "value", "/api"),
					resource.TestCheckResourceAttrPair(resourceName, "l7policy_id",
						"huaweicloud_elb_l7policy.test", "id"),
				),
			},
			{
				Config: testAccCheckElbV3L7RuleConfig_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "PATH"),
					resource.TestCheckResourceAttr(resourceName, "compare_type", "STARTS_WITH"),
					resource.TestCheckResourceAttr(resourceName, "value", "/images"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccELBL7RuleImportStateIdFunc(resourceName),
			},
		},
	})
}

func TestAccElbV3L7Rule_basic_with_conditions(t *testing.T) {
	var l7rule l7policies.Rule
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_elb_l7rule.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&l7rule,
		getELBl7RuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckElbV3L7RuleConfig_basic_with_conditions(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "QUERY_STRING"),
					resource.TestCheckResourceAttr(resourceName, "compare_type", "EQUAL_TO"),
					resource.TestCheckResourceAttrPair(resourceName, "l7policy_id",
						"huaweicloud_elb_l7policy.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.key", "key"),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.value", "value"),
				),
			},
			{
				Config: testAccCheckElbV3L7RuleConfig_update_with_conditions(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "QUERY_STRING"),
					resource.TestCheckResourceAttr(resourceName, "compare_type", "EQUAL_TO"),
					resource.TestCheckResourceAttr(resourceName, "conditions.#", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccELBL7RuleImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccELBL7RuleImportStateIdFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		l7PolicyID := rs.Primary.Attributes["l7policy_id"]
		if l7PolicyID == "" {
			return "", fmt.Errorf("attribute (l7policy_id) of Resource (%s) not found: %s", name, rs)
		}
		return fmt.Sprintf("%s/%s", l7PolicyID, rs.Primary.ID), nil
	}
}

func testAccCheckElbV3L7RuleConfig_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_elb_l7rule" "test" {
  l7policy_id  = huaweicloud_elb_l7policy.test.id
  type         = "PATH"
  compare_type = "EQUAL_TO"
  value        = "/api"
}
`, testAccCheckElbV3L7PolicyConfig_basic(rName))
}

func testAccCheckElbV3L7RuleConfig_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_elb_l7rule" "test" {
  l7policy_id  = huaweicloud_elb_l7policy.test.id
  type         = "PATH"
  compare_type = "STARTS_WITH"
  value        = "/images"
}
`, testAccCheckElbV3L7PolicyConfig_basic(rName))
}

func testAccCheckElbV3L7RuleConfig_basic_with_conditions(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_elb_l7rule" "test" {
  l7policy_id  = huaweicloud_elb_l7policy.test.id
  type         = "QUERY_STRING"
  compare_type = "EQUAL_TO"

  conditions {
    key   = "key"
    value = "value"
  }
}
`, testAccCheckElbV3L7PolicyConfig_basic(rName))
}

func testAccCheckElbV3L7RuleConfig_update_with_conditions(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_elb_l7rule" "test" {
  l7policy_id  = huaweicloud_elb_l7policy.test.id
  type         = "QUERY_STRING"
  compare_type = "EQUAL_TO"

  conditions {
    key   = "key_update"
    value = "value_update1"
  }

  conditions {
    key   = "key_update"
    value = "value_update2"
  }
}
`, testAccCheckElbV3L7PolicyConfig_basic(rName))
}
