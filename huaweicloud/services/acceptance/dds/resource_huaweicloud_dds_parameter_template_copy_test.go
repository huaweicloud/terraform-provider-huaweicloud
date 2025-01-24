package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDDSParameterTemplateCopy_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dds_parameter_template_copy.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDdsParameterTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testParameterTemplateCopy_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"_copy"),
					resource.TestCheckResourceAttr(rName, "description", "test_copy"),
					resource.TestCheckResourceAttr(rName, "node_version", "4.0"),
					resource.TestCheckResourceAttr(rName, "parameters.0.name", "connPoolMaxConnsPerHost"),
					resource.TestCheckResourceAttr(rName, "parameters.0.value", "800"),
					resource.TestCheckResourceAttr(rName, "parameters.1.name", "connPoolMaxShardedConnsPerHost"),
					resource.TestCheckResourceAttr(rName, "parameters.1.value", "800"),
				),
			},
			{
				Config: testParameterTemplateCopy_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"_copy_update"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "node_version", "4.0"),
					resource.TestCheckResourceAttr(rName, "parameters.0.name", "connPoolMaxConnsPerHost"),
					resource.TestCheckResourceAttr(rName, "parameters.0.value", "500"),
					resource.TestCheckResourceAttr(rName, "parameters.1.name", "connPoolMaxShardedConnsPerHost"),
					resource.TestCheckResourceAttr(rName, "parameters.1.value", "500"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parameter_values", "configuration_id"},
			},
		},
	})
}

func testParameterTemplateCopy_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dds_parameter_template_copy" "test" {
  configuration_id = huaweicloud_dds_parameter_template.test.id
  name             = "%[2]s_copy"
  description      = "test_copy"
}
`, testDdsParameterTemplate_basic(name), name)
}

func testParameterTemplateCopy_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dds_parameter_template_copy" "test" {
  configuration_id = huaweicloud_dds_parameter_template.test.id
  name             = "%[2]s_copy_update"
  description      = ""

  parameter_values = {
    connPoolMaxConnsPerHost        = 500
    connPoolMaxShardedConnsPerHost = 500
  }
}
`, testDdsParameterTemplate_basic(name), name)
}
