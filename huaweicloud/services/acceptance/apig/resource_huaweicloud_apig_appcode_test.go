package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/applications"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getAppcodeFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	return applications.GetAppCode(client, state.Primary.Attributes["instance_id"],
		state.Primary.Attributes["application_id"], state.Primary.ID).Extract()
}

// Auto generate APPCODE.
func TestAccAppcode_basic(t *testing.T) {
	var (
		appCode applications.AppCode

		rName = "huaweicloud_apig_appcode.test"
		rc    = acceptance.InitResourceCheck(rName, &appCode, getAppcodeFunc)
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
				Config: testAccAppcode_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAppcodeImportIdFunc(),
			},
		},
	})
}

func testAccAppcodeImportIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rName := "huaweicloud_apig_appcode.test"
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("Resource (%s) not found: %s", rName, rs)
		}
		instanceId := rs.Primary.Attributes["instance_id"]
		appId := rs.Primary.Attributes["application_id"]
		appCodeId := rs.Primary.ID
		if instanceId == "" || appId == "" || appCodeId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<application_id>/<id>', but got '%s/%s/%s'",
				instanceId, appId, appCodeId)
		}
		return fmt.Sprintf("%s/%s/%s", instanceId, appId, appCodeId), nil
	}
}

func testAccAppcode_basic() string {
	name := acceptance.RandomAccResourceName()
	return fmt.Sprintf(`
resource "huaweicloud_apig_application" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
}

resource "huaweicloud_apig_appcode" "test" {
  instance_id    = "%[1]s"
  application_id = huaweicloud_apig_application.test.id
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name)
}

// Manually configure APPCODE.
func TestAccAppcode_manuallyConfig(t *testing.T) {
	var (
		appCode applications.AppCode

		rName = "huaweicloud_apig_appcode.test"
		rc    = acceptance.InitResourceCheck(rName, &appCode, getAppcodeFunc)
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
				Config: testAccAppcode_manuallyConfig(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAppcodeImportIdFunc(),
			},
		},
	})
}

func testAccAppcode_manuallyConfig() string {
	var (
		name = acceptance.RandomAccResourceName()
		code = utils.Base64EncodeString(acctest.RandString(64))
	)
	return fmt.Sprintf(`
resource "huaweicloud_apig_application" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
}

resource "huaweicloud_apig_appcode" "test" {
  instance_id    = "%[1]s"
  application_id = huaweicloud_apig_application.test.id
  value          = "%[3]s"
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name, code)
}
