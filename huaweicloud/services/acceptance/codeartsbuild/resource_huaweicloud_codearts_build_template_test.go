package codeartsbuild

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/codeartsbuild"
)

func getBuildTemplateResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("codearts_build", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CodeArts Build client: %s", err)
	}

	return codeartsbuild.GetBuildTemplate(client, state.Primary.ID)
}

func TestAccBuildTemplate_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_build_template.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getBuildTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testBuildTemplate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "demo"),
					resource.TestCheckResourceAttrSet(rName, "parameters.#"),
					resource.TestCheckResourceAttrSet(rName, "steps.#"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"steps.1.properties"},
			},
		},
	})
}

func testBuildTemplate_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_codearts_build_template" "test" {
  name        = "%[1]s"
  description = "demo"

  parameters {
    name = "hudson.model.StringParameterDefinition"

    params {
      name  = "name"
      value = "test"
    }
    params {
      name  = "type"
      value = "customizeparam"
    }
    params {
      name  = "defaultValue"
      value = "cs"
    }
    params {
      name  = "staticVar"
      value = "false"
    }
    params {
      name  = "sensitiveVar"
      value = "false"
    }
    params {
      name  = "deletion"
      value = "false"
    }
    params {
      name  = "defaults"
      value = "false"
    }
  }

  steps {
    enable    = true
    module_id = "devcloud2018.codeci_action_20035.action"
    name      = "Docker Command"
  }

  steps {
    enable     = true
    module_id  = "devcloud2018.codeci_action_20057.action"
    name       = "update OBS"
    properties = {
      objectKey          = jsonencode("./")
      backetName         = jsonencode("test")
      uploadDirectory    = jsonencode(true)
      artifactSourcePath = jsonencode("bin/*")
      authorizationUser  = jsonencode({
        "displayName": "current user",
        "value": "build" 
      })
      obsHeaders = jsonencode([
        {
          "headerKey": "1",
          "headerValue": "1"
        }
      ])
    }
  }
}
`, name)
}
