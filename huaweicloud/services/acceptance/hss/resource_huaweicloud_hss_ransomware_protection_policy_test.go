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

func getRansomwareProtectionPolicyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region   = acceptance.HW_REGION_NAME
		epsId    = acceptance.HW_ENTERPRISE_PROJECT_ID_TEST
		policyId = state.Primary.ID
		product  = "hss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating HSS client: %s", err)
	}

	return hss.QueryProtectionPolicyByPolicyId(client, policyId, epsId, region)
}

func TestAccRansomwareProtectionPolicy_basic(t *testing.T) {
	var (
		policy interface{}
		name   = acceptance.RandomAccResourceName()
		rName  = "huaweicloud_hss_ransomware_protection_policy.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&policy,
		getRansomwareProtectionPolicyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccRansomwareProtectionPolicy_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "runtime_detection_status", "closed"),
					resource.TestCheckResourceAttr(rName, "protection_type", "rtf,doc,txt"),
					resource.TestCheckResourceAttr(rName, "protection_mode", "alarm_only"),
					resource.TestCheckResourceAttr(rName, "protection_directory", "/root;/home"),
					resource.TestCheckResourceAttr(rName, "policy_name", name),
					resource.TestCheckResourceAttr(rName, "operating_system", "Linux"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(rName, "deploy_mode", "opened"),
					resource.TestCheckResourceAttr(rName, "ai_protection_status", "opened"),
					resource.TestCheckResourceAttr(rName, "bait_protection_status", "opened"),
					resource.TestCheckResourceAttrSet(rName, "count_associated_server"),
					resource.TestCheckResourceAttrSet(rName, "default_policy"),
					resource.TestCheckResourceAttrSet(rName, "process_whitelist_attribute.#"),
				),
			},
			{
				Config: testAccRansomwareProtectionPolicy_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "policy_name", fmt.Sprintf("%s-update", name)),
					resource.TestCheckResourceAttr(rName, "protection_mode", "alarm_and_isolation"),
					resource.TestCheckResourceAttr(rName, "protection_directory", "/data"),
					resource.TestCheckResourceAttr(rName, "protection_type", "xls,ppt"),
					resource.TestCheckResourceAttr(rName, "exclude_directory", "/tmp"),
					resource.TestCheckResourceAttr(rName, "runtime_detection_status", "closed"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccRansomwareProtectionPolicyImportStateIDFunc(rName),
				// The following fields will be ignored during import verification
				ImportStateVerifyIgnore: []string{
					"process_whitelist",
					"agent_id_list",
				},
			},
		},
	})
}

func testAccRansomwareProtectionPolicy_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_ransomware_protection_policy" "test" {
  policy_name              = "%s"
  protection_mode          = "alarm_only"
  protection_directory     = "/root;/home"
  protection_type          = "rtf,doc,txt"
  operating_system         = "Linux"
  enterprise_project_id    = "%s"
  deploy_mode              = "opened"
  runtime_detection_status = "closed"
  ai_protection_status     = "opened"
  bait_protection_status   = "opened"

  process_whitelist {
    path = "/usr/bin/safe_process"
    hash = "a1b2c3d4e5f6"
  }
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccRansomwareProtectionPolicy_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_ransomware_protection_policy" "test" {
  policy_name              = "%s-update"
  protection_mode          = "alarm_and_isolation"
  protection_directory     = "/data"
  protection_type          = "xls,ppt"
  operating_system         = "Linux"
  enterprise_project_id    = "%s"
  deploy_mode              = "opened"
  exclude_directory        = "/tmp"
  runtime_detection_status = "closed"
  ai_protection_status     = "opened"
  bait_protection_status   = "opened"
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccRansomwareProtectionPolicyImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", resourceName)
		}

		epsId := rs.Primary.Attributes["enterprise_project_id"]
		id := rs.Primary.ID
		if epsId == "" || id == "" {
			return "", fmt.Errorf("invalid format specified for import ID, "+
				"want '<enterprise_project_id>/<id>', but got '%s/%s'", epsId, id)
		}
		return fmt.Sprintf("%s/%s", epsId, id), nil
	}
}
