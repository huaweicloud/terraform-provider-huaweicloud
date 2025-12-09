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

func getResourceRaspProtectionPolicyFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("hss", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HSS client: %s", err)
	}

	return hss.GetRaspProtectionPolicy(client, state.Primary.ID, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func TestAccResourceRaspProtectionPolicy_basic(t *testing.T) {
	var (
		resourceName = "huaweicloud_hss_rasp_protection_policy.test"
		name         = acceptance.RandomAccResourceName()
		updateName   = acceptance.RandomAccResourceName()

		object interface{}
		rc     = acceptance.InitResourceCheck(
			resourceName,
			&object,
			getResourceRaspProtectionPolicyFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccRaspProtectionPolicy_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "policy_name", name),
					resource.TestCheckResourceAttr(resourceName, "os_type", "Linux"),
					resource.TestCheckResourceAttr(resourceName, "feature_list.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "feature_list.0.chk_feature_id", "1"),
					resource.TestCheckResourceAttr(resourceName, "feature_list.0.protective_action", "1"),
					resource.TestCheckResourceAttr(resourceName, "feature_list.0.enabled", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "rule_list.#"),
					resource.TestCheckResourceAttrSet(resourceName, "rule_list.0.chk_feature_name"),
					resource.TestCheckResourceAttrSet(resourceName, "rule_list.0.feature_configure"),
				),
			},
			{
				Config: testAccRaspProtectionPolicy_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "policy_name", updateName),
					resource.TestCheckResourceAttr(resourceName, "os_type", "Linux"),
					resource.TestCheckResourceAttr(resourceName, "feature_list.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "feature_list.0.chk_feature_id", "3"),
					resource.TestCheckResourceAttr(resourceName, "feature_list.0.protective_action", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "rule_list.#"),
					resource.TestCheckResourceAttrSet(resourceName, "rule_list.0.chk_feature_name"),
					resource.TestCheckResourceAttrSet(resourceName, "rule_list.0.feature_configure"),
					resource.TestCheckResourceAttrSet(resourceName, "rule_list.0.enabled"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccRaspProtectionPolicyImportStateFunc(resourceName),
				ImportStateVerifyIgnore: []string{
					"feature_list",
					"enterprise_project_id",
				},
			},
		},
	})
}

func testAccRaspProtectionPolicy_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_rasp_protection_policy" "test" {
  policy_name           = "%[1]s"
  os_type               = "Linux"
  enterprise_project_id = "%[2]s"

  feature_list {
    chk_feature_id    = 1
    protective_action = 1
    enabled           = 0
    feature_configure = "/guiserver/rule/create"
  }

  feature_list {
    chk_feature_id    = 2
    protective_action = 1
    enabled           = 1
    feature_configure = "/guiserver/rule/update"
  }
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccRaspProtectionPolicy_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_rasp_protection_policy" "test" {
  policy_name           = "%[1]s"
  os_type               = "Linux"
  enterprise_project_id = "%[2]s"

  feature_list {
    chk_feature_id    = 3
    protective_action = 1
    enabled           = 1
    feature_configure = "/hips/cloudsoa-api/rasp/v1/alarmList"
  }
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccRaspProtectionPolicyImportStateFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", resourceName)
		}

		epsId := rs.Primary.Attributes["enterprise_project_id"]
		policyId := rs.Primary.ID
		if epsId == "" || policyId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, "+
				"want '<enterprise_project_id>/<id>', but got '%s/%s'", epsId, policyId)
		}
		return fmt.Sprintf("%s/%s", epsId, policyId), nil
	}
}
