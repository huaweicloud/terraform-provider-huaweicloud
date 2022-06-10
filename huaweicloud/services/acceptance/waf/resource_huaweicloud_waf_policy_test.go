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
		_, err := policies.Get(wafClient, rs.Primary.ID).Extract()
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

		found, err := policies.Get(wafClient, rs.Primary.ID).Extract()
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
