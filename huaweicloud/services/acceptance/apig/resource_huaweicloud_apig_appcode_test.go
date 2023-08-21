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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
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

func testAccApigAppcode_base() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_apig_instance" "test" {
  name                  = "%[2]s"
  edition               = "BASIC"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  enterprise_project_id = "0"

  availability_zones = try(slice(data.huaweicloud_availability_zones.test.names, 0, 1), null)
}

resource "huaweicloud_apig_application" "test" {
  instance_id = huaweicloud_apig_instance.test.id
  name        = "%[2]s"
}
`, common.TestBaseNetwork(name), name)
}

func testAccAppcode_basic() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_appcode" "test" {
  instance_id    = huaweicloud_apig_instance.test.id
  application_id = huaweicloud_apig_application.test.id
}
`, testAccApigAppcode_base())
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
	code := utils.Base64EncodeString(acctest.RandString(64))
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_appcode" "test" {
  instance_id    = huaweicloud_apig_instance.test.id
  application_id = huaweicloud_apig_application.test.id
  value          = "%[2]s"
}
`, testAccApigAppcode_base(), code)
}
