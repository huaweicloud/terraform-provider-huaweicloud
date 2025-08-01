package codeartspipeline

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/codeartspipeline"
)

func getPipelineBasicPluginResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("codearts_pipeline", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	return codeartspipeline.GetPipelineBasicPlugin(client, cfg.DomainID, state.Primary.ID)
}

func TestAccPipelineBasicPlugin_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_pipeline_basic_plugin.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPipelineBasicPluginResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPipelineBasicPlugin_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "plugin_name", name),
					resource.TestCheckResourceAttr(rName, "display_name", name),
					resource.TestCheckResourceAttr(rName, "business_type", "Build"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttrSet(rName, "runtime_attribution"),
					resource.TestCheckResourceAttrSet(rName, "plugin_composition_type"),
					resource.TestCheckResourceAttrSet(rName, "maintainers"),
				),
			},
			{
				Config: testPipelineBasicPlugin_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "plugin_name", name),
					resource.TestCheckResourceAttr(rName, "display_name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "business_type", "Build"),
					resource.TestCheckResourceAttr(rName, "description", "test"),
					resource.TestCheckResourceAttrSet(rName, "runtime_attribution"),
					resource.TestCheckResourceAttrSet(rName, "plugin_composition_type"),
					resource.TestCheckResourceAttrSet(rName, "maintainers"),
					resource.TestCheckResourceAttrSet(rName, "unique_id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testPipelineBasicPlugin_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_codearts_pipeline_basic_plugin" "test" {
  plugin_name   = "%[1]s"
  display_name  = "%[1]s"
  business_type = "Build"
}
`, name)
}

func testPipelineBasicPlugin_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_codearts_pipeline_basic_plugin" "test" {
  plugin_name   = "%[1]s"
  display_name  = "%[1]s-update"
  business_type = "Build"
  description   = "test"
}
`, name)
}
