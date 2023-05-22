package waf

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/policies"
)

func TestAccWafPolicyV1_basic(t *testing.T) {
	var policy policies.Policy
	randName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_waf_policy.policy_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafPolicyV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafPolicyV1_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafPolicyV1Exists(resourceName, &policy),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "level", "1"),
					resource.TestCheckResourceAttr(resourceName, "full_detection", "false"),
				),
			},
			{
				Config: testAccWafPolicyV1_update(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafPolicyV1Exists(resourceName, &policy),
					resource.TestCheckResourceAttr(resourceName, "name", randName+"_updated"),
					resource.TestCheckResourceAttr(resourceName, "protection_mode", "block"),
					resource.TestCheckResourceAttr(resourceName, "level", "3"),
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

func TestAccWafPolicyV1_withEpsID(t *testing.T) {
	var policy policies.Policy
	randName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_waf_policy.policy_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafPolicyV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafPolicyV1_basic_withEpsID(randName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafPolicyV1Exists(resourceName, &policy),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "level", "1"),
					resource.TestCheckResourceAttr(resourceName, "full_detection", "false"),
				),
			},
			{
				Config: testAccWafPolicyV1_update_withEpsID(randName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafPolicyV1Exists(resourceName, &policy),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "name", randName+"_updated"),
					resource.TestCheckResourceAttr(resourceName, "protection_mode", "block"),
					resource.TestCheckResourceAttr(resourceName, "level", "3"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testWAFResourceImportState(resourceName),
			},
		},
	})
}

func testAccCheckWafPolicyV1Destroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	wafClient, err := config.WafV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_waf_policy" {
			continue
		}
		_, err := policies.GetWithEpsID(wafClient, rs.Primary.ID, rs.Primary.Attributes["enterprise_project_id"]).Extract()
		if err == nil {
			return fmtp.Errorf("Waf policy still exists")
		}
	}
	return nil
}

func testAccCheckWafPolicyV1Exists(n string, policy *policies.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		wafClient, err := config.WafV1Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("error creating huaweicloud WAF client: %s", err)
		}

		found, err := policies.GetWithEpsID(wafClient, rs.Primary.ID, rs.Primary.Attributes["enterprise_project_id"]).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmtp.Errorf("Waf policy not found")
		}

		*policy = *found
		return nil
	}
}

func testAccWafPolicyV1_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_policy" "policy_1" {
  name  = "%s"
  level = 1

  depends_on = [
    huaweicloud_waf_dedicated_instance.instance_1
  ]
}
`, testAccWafDedicatedInstanceV1_conf(name), name)
}

func testAccWafPolicyV1_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_policy" "policy_1" {
  name            = "%s_updated"
  protection_mode = "block"
  level           = 3

  depends_on = [
    huaweicloud_waf_dedicated_instance.instance_1
  ]
}
`, testAccWafDedicatedInstanceV1_conf(name), name)
}

func testAccWafPolicyV1_basic_withEpsID(name, epsID string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_policy" "policy_1" {
  name                  = "%s"
  level                 = 1
  enterprise_project_id = "%s"

  depends_on = [
    huaweicloud_waf_dedicated_instance.instance_1
  ]
}
`, testAccWafDedicatedInstance_epsId(name, epsID), name, epsID)
}

func testAccWafPolicyV1_update_withEpsID(name, epsID string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_policy" "policy_1" {
  name                  = "%s_updated"
  protection_mode       = "block"
  level                 = 3
  enterprise_project_id = "%s"

  depends_on = [
    huaweicloud_waf_dedicated_instance.instance_1
  ]
}
`, testAccWafDedicatedInstance_epsId(name, epsID), name, epsID)
}
