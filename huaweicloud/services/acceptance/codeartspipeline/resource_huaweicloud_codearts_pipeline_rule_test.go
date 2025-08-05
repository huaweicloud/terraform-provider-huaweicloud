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

func getPipelineRuleResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("codearts_pipeline", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	httpUrl := "v2/{domain_id}/rules/{rule_id}/detail"
	return codeartspipeline.GetPipelineRuleItems(client, httpUrl, cfg.DomainID, state.Primary.ID)
}

func TestAccPipelineRule_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_pipeline_rule.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPipelineRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPipelineRule_basic("0.0.3", name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "Gate"),
					resource.TestCheckResourceAttr(rName, "plugin_name", "official_devcloud_codeCheck"),
				),
			},
			{
				Config: testPipelineRule_basic("0.0.1", name+"-update"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "type", "Gate"),
					resource.TestCheckResourceAttr(rName, "plugin_name", "official_devcloud_codeCheck"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"layout_content"},
			},
		},
	})
}

func testPipelineRule_basic(version, name string) string {
	return fmt.Sprintf(`
data "huaweicloud_codearts_pipeline_plugin_metrics" "test" {
  plugin_name = "official_devcloud_codeCheck"
  version     = "%[1]s"
}

locals {
  output_value = jsondecode(data.huaweicloud_codearts_pipeline_plugin_metrics.test.metrics[0].output_value)
}

resource "huaweicloud_codearts_pipeline_rule" "test" {
  name           = "%[2]s"
  type           = "Gate"
  layout_content = "layout_content"
  plugin_id      = data.huaweicloud_codearts_pipeline_plugin_metrics.test.plugin_name
  plugin_name    = data.huaweicloud_codearts_pipeline_plugin_metrics.test.plugin_name
  plugin_version = data.huaweicloud_codearts_pipeline_plugin_metrics.test.version

  dynamic "content" {
    for_each = local.output_value

    content {
      can_modify_when_inherit = true
      group_name              = content.value.group_name

      dynamic "properties" {
        for_each = content.value.properties

        content {
          key        = properties.value.key
          name       = properties.value.desc
          operator   = "="
          type       = "judge"
          value      = properties.value.value
          value_type = "float"
          is_valid   = true
        }
      }
    }
  }
}
`, version, name)
}
