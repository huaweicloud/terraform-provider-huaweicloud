package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGeminiDBParameterTemplateCopy_basic(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_geminidb_parameter_template_copy.test"
	rName := acceptance.RandomAccResourceName()
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getParameterTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDBParameterTemplateCopy_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName+"_copy"),
					resource.TestCheckResourceAttr(resourceName, "description", "This is a test parameter template copy"),
				),
			},
			{
				Config: testAccGeminiDBParameterTemplateCopy_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"_copy_update"),
					resource.TestCheckResourceAttr(resourceName, "description", "test parameter template update"),
					resource.TestCheckResourceAttr(resourceName, "values.request_timeout_in_ms", "10000"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"config_id", "values"},
			},
		},
	})
}

func testAccGeminiDBParameterTemplateCopy_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_geminidb_parameter_template_copy" "test" {
  config_id   = huaweicloud_geminidb_parameter_template.test.id
  name        = "%[2]s_copy"
  description = "This is a test parameter template copy"

  values = {
    request_timeout_in_ms = "20000"
  }
}
`, testAccGeminiDBParameterTemplate_basic(rName), rName)
}

func testAccGeminiDBParameterTemplateCopy_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_geminidb_parameter_template_copy" "test" {
  config_id   = huaweicloud_geminidb_parameter_template.test.id
  name        = "%[2]s_copy_update"
  description = "test parameter template update"

  values = {
    request_timeout_in_ms = "10000"
  }
}
`, testAccGeminiDBParameterTemplate_basic(rName), rName)
}
