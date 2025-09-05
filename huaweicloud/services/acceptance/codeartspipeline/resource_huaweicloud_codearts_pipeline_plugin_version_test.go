package codeartspipeline

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getPipelinePluginVersionResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("codearts_pipeline", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	httpUrl := "v1/{domain_id}/agent-plugin/detail?version={version}&plugin_name={plugin_name}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{domain_id}", cfg.DomainID)
	getPath = strings.ReplaceAll(getPath, "{version}", state.Primary.Attributes["version"])
	getPath = strings.ReplaceAll(getPath, "{plugin_name}", state.Primary.Attributes["plugin_name"])
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	return getRespBody, nil
}

func TestAccPipelinePluginVersion_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_pipeline_plugin_version.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPipelinePluginVersionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPipelinePluginVersion_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "plugin_name", name),
					resource.TestCheckResourceAttr(rName, "display_name", name),
					resource.TestCheckResourceAttr(rName, "version", "0.0.1"),
					resource.TestCheckResourceAttr(rName, "version_description", "version demo"),
					resource.TestCheckResourceAttr(rName, "description", "plugin demo"),
					resource.TestCheckResourceAttr(rName, "runtime_attribution", "agent"),
					resource.TestCheckResourceAttr(rName, "input_info.0.name", "key"),
				),
			},
			{
				Config: testPipelinePluginVersion_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "plugin_name", name),
					resource.TestCheckResourceAttr(rName, "display_name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "version", "0.0.1"),
					resource.TestCheckResourceAttr(rName, "version_description", ""),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "runtime_attribution", "agent"),
				),
			},
			{
				Config: testPipelinePluginVersion_update_published(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "plugin_name", name),
					resource.TestCheckResourceAttr(rName, "display_name", name+"-update-again"),
					resource.TestCheckResourceAttr(rName, "version", "0.0.1"),
					resource.TestCheckResourceAttr(rName, "version_description", ""),
					resource.TestCheckResourceAttr(rName, "description", "test"),
					resource.TestCheckResourceAttr(rName, "runtime_attribution", "agent"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"execution_info"},
			},
		},
	})
}

func TestAccPipelinePluginVersion_createFormal(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_pipeline_plugin_version.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPipelinePluginVersionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPipelinePluginVersion_update_published(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "plugin_name", name),
					resource.TestCheckResourceAttr(rName, "display_name", name+"-update-again"),
					resource.TestCheckResourceAttr(rName, "version", "0.0.1"),
					resource.TestCheckResourceAttr(rName, "version_description", ""),
					resource.TestCheckResourceAttr(rName, "description", "test"),
					resource.TestCheckResourceAttr(rName, "runtime_attribution", "agent"),
				),
			},
		},
	})
}

func TestAccPipelinePluginVersion_deleteDraft(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_pipeline_plugin_version.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPipelinePluginVersionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPipelinePluginVersion_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "plugin_name", name),
					resource.TestCheckResourceAttr(rName, "display_name", name),
					resource.TestCheckResourceAttr(rName, "version", "0.0.1"),
					resource.TestCheckResourceAttr(rName, "version_description", "version demo"),
					resource.TestCheckResourceAttr(rName, "description", "plugin demo"),
					resource.TestCheckResourceAttr(rName, "runtime_attribution", "agent"),
					resource.TestCheckResourceAttr(rName, "input_info.0.name", "key"),
				),
			},
		},
	})
}

//nolint:revive
func testPipelinePluginVersion_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_codearts_pipeline_plugin_version" "test" {
  plugin_name         = "%[1]s"
  display_name        = "%[1]s"
  business_type       = "Build"
  version             = "0.0.1"
  version_description = "version demo"
  description         = "plugin demo"
  runtime_attribution = "agent"

  execution_info {
    inner_execution_info = jsonencode({
      "execution_type": "COMPOSITE",
      "steps": [{
        "identifier": "1756954499618a927e16a-7a8e-48b8-b732-76751313ea53"
        "task": "official_shell_plugin",
        "name": "执行Shell",
        "variables": {
          "OFFICIAL_SHELL_SCRIPT_INPUT": "(decode)ZWNobyAiaGVsbG8i"
        }
      }]
    })
  }

  input_info {
    layout_content = "[{\"componentName\":\"d-form\",\"id\":1756968435834,\"category\":\"form\",\"props\":{\"layout\":\"vertical\",\"ref\":\"formRef\"},\"children\":[{\"componentName\":\"d-form-item\",\"id\":1756968441813,\"category\":\"form\",\"props\":{\"label\":\"num\",\"field\":\"key_1756968441812\",\"help-tips\":\"\",\"visible\":{\"type\":\"fixed\",\"conditions\":[],\"value\":true},\"disabled\":false},\"children\":[{\"componentName\":\"d-input-number\",\"id\":1756968441812,\"props\":{\"showPlaceholder\":\"\",\"type\":true,\"precision\":0,\"disabled\":{\"type\":\"fixed\",\"conditions\":[],\"value\":false}},\"category\":\"form\"}]}]}]"
    name           = "key"
    type           = "d-input-number"
  }
}
`, name)
}

func testPipelinePluginVersion_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_codearts_pipeline_plugin_version" "test" {
  plugin_name         = "%[1]s"
  display_name        = "%[1]s-update"
  business_type       = "Build"
  version             = "0.0.1"
  runtime_attribution = "agent"
  is_formal           = true

  execution_info {
    inner_execution_info = jsonencode({
      "execution_type": "COMPOSITE",
      "steps": [{
        "identifier": "1756954499618a927e16a-7a8e-48b8-b732-76751313ea53"
        "task": "official_shell_plugin",
        "name": "execute shell",
        "variables": {
          "OFFICIAL_SHELL_SCRIPT_INPUT": "(decode)ZWNobyAiaGVsbG8i"
        }
      }]
    })
  }
}
`, name)
}

func testPipelinePluginVersion_update_published(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_codearts_pipeline_plugin_version" "test" {
  plugin_name         = "%[1]s"
  display_name        = "%[1]s-update-again"
  business_type       = "Build"
  version             = "0.0.1"
  description         = "test"
  runtime_attribution = "agent"
  is_formal           = true

  execution_info {
    inner_execution_info = jsonencode({
      "execution_type": "COMPOSITE",
      "steps": [{
        "identifier": "1756954499618a927e16a-7a8e-48b8-b732-76751313ea53"
        "task": "official_shell_plugin",
        "name": "execute shell",
        "variables": {
          "OFFICIAL_SHELL_SCRIPT_INPUT": "(decode)ZWNobyAiaGVsbG8i"
        }
      }]
    })
  }
}
`, name)
}
