package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/apig"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getAppcodeFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("apig", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG client: %s", err)
	}
	return apig.GetAppcode(client, state.Primary.Attributes["instance_id"],
		state.Primary.Attributes["application_id"], state.Primary.ID)
}

func TestAccAppcode_basic(t *testing.T) {
	var (
		obj interface{}

		autoGeneration   = "huaweicloud_apig_appcode.auto_generate"
		rcAutoGeneration = acceptance.InitResourceCheck(autoGeneration, &obj, getAppcodeFunc)

		manuallyConfig   = "huaweicloud_apig_appcode.manually_config"
		rcManuallyConfig = acceptance.InitResourceCheck(manuallyConfig, &obj, getAppcodeFunc)

		name = acceptance.RandomAccResourceName()
		code = utils.Base64EncodeString(acctest.RandString(64))
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcAutoGeneration.CheckResourceDestroy(),
			rcManuallyConfig.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config:      testAccAppcode_withValueAndNonExistentParentResources(code),
				ExpectError: regexp.MustCompile(`error creating APPCODE`),
			},
			{
				Config:      testAccAppcode_withoutValueAndNonExistentParentResources(),
				ExpectError: regexp.MustCompile(`error auto generating APPCODE`),
			},
			{
				Config: testAccAppcode_basic(name, code),
				Check: resource.ComposeTestCheckFunc(
					rcAutoGeneration.CheckResourceExists(),
					resource.TestCheckResourceAttr(autoGeneration, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttrPair(autoGeneration, "application_id", "huaweicloud_apig_application.test", "id"),
					resource.TestMatchResourceAttr(autoGeneration, "value",
						regexp.MustCompile(`^[a-zA-Z0-9/+][a-zA-Z0-9+_!@#$%/=-]{63,179}$`)),
					resource.TestMatchResourceAttr(autoGeneration, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					rcManuallyConfig.CheckResourceExists(),
					resource.TestCheckResourceAttr(manuallyConfig, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttrPair(manuallyConfig, "application_id", "huaweicloud_apig_application.test", "id"),
					resource.TestCheckResourceAttr(manuallyConfig, "value", code),
					resource.TestMatchResourceAttr(manuallyConfig, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:      autoGeneration,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAppcodeImportIdFunc(autoGeneration),
			},
			{
				ResourceName:      manuallyConfig,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAppcodeImportIdFunc(manuallyConfig),
			},
		},
	})
}

func testAccAppcode_withValueAndNonExistentParentResources(code string) string {
	randomUUID, _ := uuid.NewRandom()

	return fmt.Sprintf(`
resource "huaweicloud_apig_appcode" "test" {
  instance_id    = "%[1]s"
  application_id = "%[1]s"
  value          = "%[2]s"
}
`, randomUUID.String(), code)
}

func testAccAppcode_withoutValueAndNonExistentParentResources() string {
	randomUUID, _ := uuid.NewRandom()

	return fmt.Sprintf(`
resource "huaweicloud_apig_appcode" "test" {
  instance_id    = "%[1]s"
  application_id = "%[1]s"
}
`, randomUUID.String())
}

func testAccAppcode_basic(name, code string) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_application" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
}

resource "huaweicloud_apig_appcode" "auto_generate" {
  instance_id    = "%[1]s"
  application_id = huaweicloud_apig_application.test.id
}

resource "huaweicloud_apig_appcode" "manually_config" {
  instance_id    = "%[1]s"
  application_id = huaweicloud_apig_application.test.id
  value          = "%[3]s"
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name, code)
}

func testAccAppcodeImportIdFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}
		instanceId := rs.Primary.Attributes["instance_id"]
		appId := rs.Primary.Attributes["application_id"]
		appCodeId := rs.Primary.ID
		if instanceId == "" || appId == "" || appCodeId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<application_id>/<id>', "+
				"but got '%s/%s/%s'", instanceId, appId, appCodeId)
		}
		return fmt.Sprintf("%s/%s/%s", instanceId, appId, appCodeId), nil
	}
}
