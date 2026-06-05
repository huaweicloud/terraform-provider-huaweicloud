package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/apig"
)

func getApplicationAiApiKeyFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("apig", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG client: %s", err)
	}
	return apig.GetApplicationAiApiKey(client, state.Primary.Attributes["instance_id"],
		state.Primary.Attributes["application_id"], state.Primary.ID)
}

func TestAccApplicationAiApiKey_basic(t *testing.T) {
	var (
		obj interface{}

		autoGeneration   = "huaweicloud_apig_application_ai_api_key.auto_generate"
		rcAutoGeneration = acceptance.InitResourceCheck(autoGeneration, &obj, getApplicationAiApiKeyFunc)

		manuallyConfig   = "huaweicloud_apig_application_ai_api_key.manually_config"
		rcManuallyConfig = acceptance.InitResourceCheck(manuallyConfig, &obj, getApplicationAiApiKeyFunc)

		manuallyConfigWithAlias   = "huaweicloud_apig_application_ai_api_key.manually_config_with_alias"
		rcManuallyConfigWithAlias = acceptance.InitResourceCheck(manuallyConfigWithAlias, &obj, getApplicationAiApiKeyFunc)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: "3.0.0",
			},
		},
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcAutoGeneration.CheckResourceDestroy(),
			rcManuallyConfig.CheckResourceDestroy(),
			rcManuallyConfigWithAlias.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationAiApiKey_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// The value of the AI API key is auto generated.
					rcAutoGeneration.CheckResourceExists(),
					resource.TestCheckResourceAttr(autoGeneration, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttrPair(autoGeneration, "application_id", "huaweicloud_apig_application.test", "id"),
					resource.TestMatchResourceAttr(autoGeneration, "value", regexp.MustCompile(`^[a-zA-Z0-9_+/=-]{8,128}$`)),
					resource.TestCheckResourceAttr(autoGeneration, "alias", ""),
					resource.TestMatchResourceAttr(autoGeneration, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// The value of the AI API key is manually configured.
					rcManuallyConfig.CheckResourceExists(),
					resource.TestCheckResourceAttr(manuallyConfig, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttrPair(manuallyConfig, "application_id", "huaweicloud_apig_application.test", "id"),
					resource.TestCheckResourceAttrPair(manuallyConfig, "value", "random_string.test.0", "result"),
					resource.TestMatchResourceAttr(manuallyConfig, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttr(manuallyConfig, "alias", ""),
					// The value of the AI API key is manually configured with alias.
					rcManuallyConfigWithAlias.CheckResourceExists(),
					resource.TestCheckResourceAttr(manuallyConfigWithAlias, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttrPair(manuallyConfigWithAlias, "application_id", "huaweicloud_apig_application.test", "id"),
					resource.TestCheckResourceAttrPair(manuallyConfigWithAlias, "value", "random_string.test.1", "result"),
					resource.TestMatchResourceAttr(manuallyConfigWithAlias, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttr(manuallyConfigWithAlias, "alias", name),
				),
			},
			{
				ResourceName:      autoGeneration,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccApplicationAiApiKeyImportIdFunc(autoGeneration),
			},
			{
				ResourceName:      manuallyConfig,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccApplicationAiApiKeyImportIdFunc(manuallyConfig),
			},
			{
				ResourceName:      manuallyConfigWithAlias,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccApplicationAiApiKeyImportIdFunc(manuallyConfigWithAlias),
			},
		},
	})
}

func testAccApplicationAiApiKey_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_apig_instances" "test" {
  instance_id = "%[1]s"
}

locals {
  instance_id = try(data.huaweicloud_apig_instances.test.instances[0].id, null)
}

resource "huaweicloud_apig_instance_feature" "test" {
  instance_id = local.instance_id
  name        = "ai_api_key_enabled"
  enabled     = true
  config      = "{}"
}

resource "huaweicloud_apig_application" "test" {
  depends_on = [huaweicloud_apig_instance_feature.test]

  instance_id = local.instance_id
  name        = "%[2]s"
}

resource "random_string" "test" {
  count = 2

  length           = 100
  min_numeric      = 1
  min_lower        = 1
  min_upper        = 1
  min_special      = 1
  override_special = "-_+/="
}

resource "huaweicloud_apig_application_ai_api_key" "auto_generate" {
  instance_id    = local.instance_id
  application_id = huaweicloud_apig_application.test.id
}

resource "huaweicloud_apig_application_ai_api_key" "manually_config" {
  instance_id    = local.instance_id
  application_id = huaweicloud_apig_application.test.id
  value          = random_string.test[0].result
}

resource "huaweicloud_apig_application_ai_api_key" "manually_config_with_alias" {
  instance_id    = local.instance_id
  application_id = huaweicloud_apig_application.test.id
  value          = random_string.test[1].result
  alias          = "%[2]s"
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name)
}

func testAccApplicationAiApiKeyImportIdFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}
		instanceId := rs.Primary.Attributes["instance_id"]
		appId := rs.Primary.Attributes["application_id"]
		aiApiKeyId := rs.Primary.ID
		if instanceId == "" || appId == "" || aiApiKeyId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<application_id>/<id>', "+
				"but got '%s/%s/%s'", instanceId, appId, aiApiKeyId)
		}
		return fmt.Sprintf("%s/%s/%s", instanceId, appId, aiApiKeyId), nil
	}
}
