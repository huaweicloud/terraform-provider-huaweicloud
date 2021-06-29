package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/apigw/v2/applications"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func TestAccApigApplicationV2_basic(t *testing.T) {
	var (
		// Only letters, digits and underscores (_) are allowed in the name.
		rName        = fmt.Sprintf("tf_acc_test_%s", acctest.RandString(5))
		resourceName = "huaweicloud_apig_application.test"
		application  applications.Application
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApigApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApigApplication_basic(rName, acctest.RandString(64)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigApplicationExists(resourceName, &application),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by script"),
				),
			},
			{
				// update name, description and app_code.
				Config: testAccApigApplication_update(rName, acctest.RandString(64)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigApplicationExists(resourceName, &application),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"_update"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated by script"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccApigInstanceSubResourceImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccCheckApigApplicationDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	client, err := config.ApigV2Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_apig_application" {
			continue
		}
		_, err := applications.Get(client, rs.Primary.Attributes["instance_id"], rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("APIG v2 application (%s) is still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckApigApplicationExists(appName string, app *applications.Application) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[appName]
		if !ok {
			return fmt.Errorf("Resource %s not found", appName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No APIG V2 application Id")
		}

		config := testAccProvider.Meta().(*config.Config)
		client, err := config.ApigV2Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
		}
		found, err := applications.Get(client, rs.Primary.Attributes["instance_id"], rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*app = *found
		return nil
	}
}

func testAccApigInstanceSubResourceImportStateIdFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" || rs.Primary.Attributes["instance_id"] == "" {
			return "", fmt.Errorf("resource not found: %s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.ID), nil
	}
}

func testAccApigApplication_basic(rName, code string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_application" "test" {
  name        = "%s"
  instance_id = huaweicloud_apig_instance.test.id
  description = "Created by script"

  app_codes = ["%s"]
}
`, testAccApigInstance_basic(rName), rName, utils.EncodeBase64String(code))
}

func testAccApigApplication_update(rName, code string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_application" "test" {
  name        = "%s_update"
  instance_id = huaweicloud_apig_instance.test.id
  description = "Updated by script"

  app_codes = ["%s"]
}
`, testAccApigInstance_basic(rName), rName, utils.EncodeBase64String(code))
}
