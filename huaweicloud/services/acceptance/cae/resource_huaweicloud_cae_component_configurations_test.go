package cae

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cae"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getComponentConfigurationsFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		httpUrl     = "v1/{project_id}/cae/applications/{application_id}/components/{component_id}/configurations"
		componentId = state.Primary.Attributes["component_id"]
	)
	client, err := cfg.NewServiceClient("cae", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CAE client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{application_id}", state.Primary.Attributes["application_id"])
	getPath = strings.ReplaceAll(getPath, "{component_id}", componentId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"X-Environment-ID": state.Primary.Attributes["environment_id"],
		},
	}
	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, common.ConvertExpected400ErrInto404Err(err, "error_code", cae.ConfigRelatedResourcesNotFoundCodes...)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving configurations of the specified component (%s): %s", componentId, err)
	}
	items := cae.FilterActivatedConfigurations(utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{}))
	if len(items) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}
	return items, nil
}

func TestAccComponentConfigurations_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_cae_component_configurations.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getComponentConfigurationsFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCaeEnvironment(t)
			acceptance.TestAccPreCheckCaeApplication(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccComponentConfiguration_basic_step1(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "component_id",
						"huaweicloud_cae_component.test", "id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccComponentConfigurationsFunc(rName),
			},
		},
	})
}

func testAccComponentConfigurationsFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var environmentId, applicaitonId, componentId string
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		environmentId = rs.Primary.Attributes["environment_id"]
		applicaitonId = rs.Primary.Attributes["application_id"]
		componentId = rs.Primary.ID
		if environmentId == "" || applicaitonId == "" || componentId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<environment_id>/<application_id>/<component_id>', "+
				"but got '%s/%s/%s'", environmentId, applicaitonId, componentId)
		}
		return fmt.Sprintf("%s/%s/%s", environmentId, applicaitonId, componentId), nil
	}
}

func testAccComponentConfiguration_base() string {
	name := acceptance.RandomAccResourceNameWithDash()

	return fmt.Sprintf(`
resource "huaweicloud_cae_component" "test" {
  environment_id = "%[1]s"
  application_id = "%[2]s"

  metadata {
    name = "%[3]s"

    annotations = {
      version = "1.0.0"
    }
  }

  spec {
    replica = 1
    runtime = "Docker"

    source {
      type = "image"
      url  = "nginx:alpine-perl"
    }

    resource_limit {
      cpu    = "500m"
      memory = "1Gi"
    }
  }
}
`, acceptance.HW_CAE_ENVIRONMENT_ID, acceptance.HW_CAE_APPLICATION_ID, name)
}

func testAccComponentConfiguration_basic_step1() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cae_component_configurations" "test" {
  environment_id = "%[2]s"
  application_id = "%[3]s"
  component_id   = huaweicloud_cae_component.test.id

  items {
    type = "lifecycle"
    data = jsonencode({
      "spec": {
        "postStart": {
          "exec": {
            "command": [
              "/bin/bash",
              "-c",
              "sleep",
              "10",
              "done",
            ]
          }
        }
      }
    })
  }
  items {
    type = "env"
    data = jsonencode({
      "spec": {
        "envs": {
            "key": "value",
            "foo": "bar"
        }
      }
    })
  }
}
`, testAccComponentConfiguration_base(),
		acceptance.HW_CAE_ENVIRONMENT_ID, acceptance.HW_CAE_APPLICATION_ID)
}
