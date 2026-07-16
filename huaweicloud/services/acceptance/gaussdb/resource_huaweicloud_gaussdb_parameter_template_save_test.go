package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/gaussdb"
)

func getResourceParameterTemplateSaveFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("opengauss", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating GaussDB client: %s", err)
	}

	return gaussdb.GetParameterTemplateSaveInfo(client, state.Primary.ID)
}

func TestAccResourceParameterTemplateSave_basic(t *testing.T) {
	var (
		rName  = "huaweicloud_gaussdb_parameter_template_save.test"
		name   = acceptance.RandomAccResourceName()
		object interface{}
		rc     = acceptance.InitResourceCheck(
			rName,
			&object,
			getResourceParameterTemplateSaveFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This parameter indicates the parameter template ID of the specified instance.
			acceptance.TestAccPreCheckGaussDBParameterTemplateId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccParameterTemplateSave_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "config_id", acceptance.HW_GAUSSDB_PARAMETER_TEMPLATE_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "terraform test"),
				),
			},
			{
				Config: testAccParameterTemplateSave_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "description", ""),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"config_id", "values"},
			},
		},
	})
}

func testAccParameterTemplateSave_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_gaussdb_parameter_template_save" "test" {
  config_id   = "%[1]s"
  name        = "%[2]s"
  description = "terraform test"

  values = {
    asp_retention_days = "5"
    autoanalyze        = "on"
  }
}
`, acceptance.HW_GAUSSDB_PARAMETER_TEMPLATE_ID, name)
}

func testAccParameterTemplateSave_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_gaussdb_parameter_template_save" "test" {
  config_id   = "%[1]s"
  name        = "%[2]s"
  description = ""

  values = {
    asp_retention_days  = "4"
    autoanalyze_timeout = "400"
  }
}
`, acceptance.HW_GAUSSDB_PARAMETER_TEMPLATE_ID, name)
}
