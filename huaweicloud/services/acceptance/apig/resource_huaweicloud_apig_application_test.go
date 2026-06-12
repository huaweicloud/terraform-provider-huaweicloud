package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/apig"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getApplicationFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("apig", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG client: %s", err)
	}
	return apig.GetApplication(client, state.Primary.Attributes["instance_id"], state.Primary.ID)
}

func TestAccApplication_basic(t *testing.T) {
	var (
		obj interface{}

		rName         = "huaweicloud_apig_application.test"
		rc            = acceptance.InitResourceCheck(rName, &obj, getApplicationFunc)
		resertSecret  = "huaweicloud_apig_application.reset_secret"
		rcResetSecret = acceptance.InitResourceCheck(resertSecret, &obj, getApplicationFunc)

		// Only letters, digits and underscores (_) are allowed in the environment name and dedicated instance name.
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		appKey          = acctest.RandString(32)
		updateAppKey    = acctest.RandString(32)
		appSecret       = acctest.RandString(64)
		updateAppSecret = acctest.RandString(64)
		code            = utils.Base64EncodeString(acctest.RandString(64))
		updateCode      = utils.Base64EncodeString(acctest.RandString(64))
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccApplication_basic_step1(name, appKey, appSecret, code),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Created by terraform test"),
					resource.TestCheckResourceAttr(rName, "app_key", appKey),
					resource.TestCheckResourceAttr(rName, "app_secret", appSecret),
					resource.TestCheckResourceAttr(rName, "app_codes.#", "1"),
					resource.TestCheckResourceAttr(rName, "app_codes.0", code),
					resource.TestMatchResourceAttr(rName, "registration_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(rName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					rcResetSecret.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resertSecret, "app_key"),
					resource.TestCheckResourceAttrSet(resertSecret, "app_secret"),
				),
			},
			{
				// update name, description, app_codes, app_key and app_secret.
				Config: testAccApplication_basic_step2(updateName, updateAppKey, updateAppSecret, updateCode),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", "Updated by terraform test"),
					resource.TestCheckResourceAttr(rName, "app_key", updateAppKey),
					resource.TestCheckResourceAttr(rName, "app_secret", updateAppSecret),
					resource.TestCheckResourceAttr(rName, "app_codes.#", "1"),
					resource.TestCheckResourceAttr(rName, "app_codes.0", updateCode),
					rcResetSecret.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resertSecret, "app_secret"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccApplicationImportIdFunc(),
			},
		},
	})
}

func testAccApplicationImportIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rName := "huaweicloud_apig_application.test"
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("Resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.ID == "" || rs.Primary.Attributes["instance_id"] == "" {
			return "", fmt.Errorf("resource not found: %s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.ID), nil
	}
}

func testAccApplication_basic_step1(name, appKey, appSecret, code string) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_application" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  description = "Created by terraform test"
  app_key     = "%[3]s"
  app_secret  = "%[4]s"
  app_codes   = ["%[5]s"]
}

resource "huaweicloud_apig_application" "reset_secret" {
  instance_id = "%[1]s"
  name        = "%[2]s_reset"
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name, appKey, appSecret, code)
}

func testAccApplication_basic_step2(name, appKey, appSecret, code string) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_application" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  description = "Updated by terraform test"
  app_key     = "%[3]s"
  app_secret  = "%[4]s"
  app_codes   = ["%[5]s"]
}

resource "huaweicloud_apig_application" "reset_secret" {
  instance_id   = "%[1]s"
  name          = "%[2]s_reset"
  secret_action = "RESET"
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name, appKey, appSecret, code)
}
