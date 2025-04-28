package cci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cci"
)

func getV2DeploymentResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("cci", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CCI client: %s", err)
	}
	return cci.GetV2Deployment(client, state.Primary.Attributes["namespace"], state.Primary.Attributes["name"])
}

func TestAccV2Deployment_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_cciv2_deployment.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getV2DeploymentResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV2Deployment_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "namespace", rName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"template.0.metadata.0.annotations"},
				ImportStateIdFunc:       testAccV2DeploymentImportStateFunc(resourceName),
			},
		},
	})
}

func testAccV2DeploymentImportStateFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["namespace"] == "" || rs.Primary.Attributes["name"] == "" {
			return "", fmt.Errorf("the namespace (%s) or name(%s) is nil",
				rs.Primary.Attributes["namespace"], rs.Primary.Attributes["name"])
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["namespace"], rs.Primary.Attributes["name"]), nil
	}
}

func testAccV2Deployment_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cciv2_deployment" "test" {
  namespace = huaweicloud_cciv2_namespace.test.name
  name      = "%[2]s"

  selector {
    match_labels = {
      app = "template1"
    }
  }

  template {
    metadata {
      labels = {
        app = "template1"
      }

      annotations = {
        "resource.cci.io/instance-type" = "general-computing"
      }
    }

    spec {
      containers {
        name  = "c1"
        image = "alpine:latest"

        resources {
          limits = {
            cpu    = "1"
            memory = "2G"
          }

          requests = {
            cpu    = "1"
            memory = "2G"
          }
        }
      }

      image_pull_secrets {
        name = "imagepull-secret"
      }
    }
  }

  lifecycle {
    ignore_changes = [
      template.0.metadata.0.annotations,
    ]
  }
}
`, testAccV2Namespace_basic(rName), rName)
}
