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

func getApplicationFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	return applications.Get(client, state.Primary.Attributes["instance_id"], state.Primary.ID).Extract()
}

func TestAccApplication_basic(t *testing.T) {
	var (
		app applications.Application

		rName = "huaweicloud_apig_application.test"
		// Only letters, digits and underscores (_) are allowed in the environment name and dedicated instance name.
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		description       = "Created by script"
		updateDescription = "Updated by script"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&app,
		getApplicationFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccApplication_basic(name, description),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", description),
					resource.TestCheckResourceAttrSet(rName, "app_key"),
					resource.TestCheckResourceAttrSet(rName, "app_secret"),
				),
			},
			{
				// update name, description and app_code.
				Config: testAccApplication_basic(updateName, updateDescription),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", updateDescription),
					resource.TestCheckResourceAttrSet(rName, "app_key"),
					resource.TestCheckResourceAttrSet(rName, "app_secret"),
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

func testAccApplication_basic(name, description string) string {
	code := utils.Base64EncodeString(acctest.RandString(64))

	return fmt.Sprintf(`
resource "huaweicloud_apig_application" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  description = "%[3]s"

  app_codes = ["%[4]s"]
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name, description, code)
}
