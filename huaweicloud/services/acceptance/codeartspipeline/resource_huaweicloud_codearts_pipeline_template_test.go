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

func getPipelineTemplateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("codearts_pipeline", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	return codeartspipeline.GetPipelineTemplate(client, cfg.DomainID, state.Primary.ID)
}

func TestAccPipelineTemplate_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_pipeline_template.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPipelineTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPipelineTemplate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "language", "none"),
					resource.TestCheckResourceAttr(rName, "is_show_source", "true"),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "is_favorite", "true"),
					resource.TestCheckResourceAttrSet(rName, "definition"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
					resource.TestCheckResourceAttrSet(rName, "update_time"),
					resource.TestCheckResourceAttr(rName, "variables.0.name", "test_var"),
					resource.TestCheckResourceAttr(rName, "variables.0.type", "string"),
					resource.TestCheckResourceAttr(rName, "variables.0.value", "test_value"),
				),
			},
			{
				Config: testPipelineTemplate_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "language", "Java"),
					resource.TestCheckResourceAttr(rName, "is_favorite", "false"),
					resource.TestCheckResourceAttr(rName, "is_show_source", "false"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "variables.#", "1"),
					resource.TestCheckResourceAttrSet(rName, "definition"),
					resource.TestCheckResourceAttr(rName, "variables.0.name", "test_var"),
					resource.TestCheckResourceAttr(rName, "variables.0.type", "string"),
					resource.TestCheckResourceAttr(rName, "variables.0.value", "test_value_update"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"is_favorite"},
			},
		},
	})
}

func testPipelineTemplate_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_codearts_pipeline_template" "test" {
  name           = "%s"
  language       = "none"
  is_show_source = true
  description    = "test description"
  is_favorite    = true
  definition     = jsonencode({
    "stages": [
      {
        "name": "Stage_1",
        "sequence": "0",
        "jobs": [
          {
            "stage_id": 1698069093677,
            "identifier": "16980691185419778b0f5-02cd-4bd3-a8cc-da825644add1",
            "name": "CodeCheck",
            "depends_on": [],
            "timeout": "",
            "timeout_unit": "",
            "steps": [
              {
                "name": "CodeCheck",
                "task": "official_devcloud_codeCheck_template",
                "sequence": 0,
                "inputs": [
                  {
                    "key": "language",
                    "value": "Java"
                  },
                  {
                    "key": "module_or_template_id",
                    "value": "d7dffaefb6d94c63a09cf141668356c7"
                  }
                ],
                "business_type": "Gate",
                "runtime_attribution": "agent",
                "identifier": "16980691136015f94b249-102d-44f2-abb5-02fa3ad9fae0",
                "multi_step_editable": 0,
                "official_task_version": "0.0.1"
              }
            ],
            "resource": "{\"type\":\"system\",\"arch\":\"x86\"}",
            "condition": "$${{ completed() }}",
            "exec_type": "OCTOPUS_JOB",
            "sequence": 0
          }
        ],
        "identifier": "16980690936778259e2a5-a97e-4c66-af8d-1908988c1c21",
        "pre": [
          {
            "task": "official_devcloud_autoTrigger",
            "sequence": 0
          }
        ],
        "post": null,
        "depends_on": [],
        "run_always": false
      }
    ]
  })

  variables {
    name  = "test_var"
    type  = "string"
    value = "test_value"
  }
}
`, name)
}

func testPipelineTemplate_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_codearts_pipeline_template" "test" {
  name           = "%s-update"
  language       = "Java"
  is_show_source = false
  description    = ""
  is_favorite    = false
  definition     = jsonencode({
    "stages": [
      {
        "name": "Stage_1",
        "sequence": "0",
        "jobs": [
          {
            "stage_id": 1698069093677,
            "identifier": "16980691185419778b0f5-02cd-4bd3-a8cc-da825644add1",
            "name": "CodeCheck",
            "depends_on": [],
            "timeout": "",
            "timeout_unit": "",
            "steps": [
              {
                "name": "CodeCheck",
                "task": "official_devcloud_codeCheck_template",
                "sequence": 0,
                "inputs": [
                  {
                    "key": "language",
                    "value": "Java"
                  },
                  {
                    "key": "module_or_template_id",
                    "value": "d7dffaefb6d94c63a09cf141668356c7"
                  }
                ],
                "business_type": "Gate",
                "runtime_attribution": "agent",
                "identifier": "16980691136015f94b249-102d-44f2-abb5-02fa3ad9fae0",
                "multi_step_editable": 0,
                "official_task_version": "0.0.1"
              }
            ],
            "resource": "{\"type\":\"system\",\"arch\":\"x86\"}",
            "condition": "$${{ completed() }}",
            "exec_type": "OCTOPUS_JOB",
            "sequence": 0
          }
        ],
        "identifier": "16980690936778259e2a5-a97e-4c66-af8d-1908988c1c21",
        "pre": [
          {
            "task": "official_devcloud_autoTrigger",
            "sequence": 0
          }
        ],
        "post": null,
        "depends_on": [],
        "run_always": false
      }
    ]
  })
  
  variables {
    name  = "test_var"
    type  = "string"
    value = "test_value_update"
  }
}
`, name)
}
