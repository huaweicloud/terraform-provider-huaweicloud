package cae

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cae"
)

func getComponentFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	environmentId := state.Primary.Attributes["environment_id"]
	applicationId := state.Primary.Attributes["application_id"]
	return cae.GetComponentById(cfg, acceptance.HW_REGION_NAME, environmentId, applicationId, state.Primary.ID)
}

func TestAccComponent_basic(t *testing.T) {
	var (
		obj        interface{}
		rName      = "huaweicloud_cae_component.test"
		name       = acceptance.RandomAccResourceNameWithDash()
		updateName = acceptance.RandomAccResourceNameWithDash()
		rc         = acceptance.InitResourceCheck(
			rName,
			&obj,
			getComponentFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCaeEnvironment(t)
			acceptance.TestAccPreCheckCaeApplication(t)
			acceptance.TestAccPreCheckCaeComponent(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccComponent_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "metadata.0.name", name),
					resource.TestCheckResourceAttr(rName, "environment_id", acceptance.HW_CAE_ENVIRONMENT_ID),
					resource.TestCheckResourceAttr(rName, "application_id", acceptance.HW_CAE_APPLICATION_ID),
					resource.TestCheckResourceAttr(rName, "metadata.0.annotations.version", "1.0.0"),
					resource.TestCheckResourceAttr(rName, "spec.0.replica", "2"),
					resource.TestCheckResourceAttr(rName, "spec.0.runtime", "Docker"),
					resource.TestCheckResourceAttr(rName, "spec.0.resource_limit.0.cpu", "1000m"),
					resource.TestCheckResourceAttr(rName, "spec.0.resource_limit.0.memory", "4Gi"),
					resource.TestCheckResourceAttr(rName, "spec.0.source.0.type", "image"),
					resource.TestCheckResourceAttr(rName, "spec.0.source.0.url", acceptance.HW_CAE_IMAGE_URL),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testAccComponent_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "metadata.0.name", updateName),
					resource.TestCheckResourceAttr(rName, "spec.0.replica", "1"),
					resource.TestCheckResourceAttr(rName, "spec.0.%", "5"),
					resource.TestCheckResourceAttr(rName, "spec.0.runtime", "Java17"),
					resource.TestCheckResourceAttr(rName, "spec.0.resource_limit.0.cpu", "500m"),
					resource.TestCheckResourceAttr(rName, "spec.0.resource_limit.0.memory", "1Gi"),
					resource.TestCheckResourceAttr(rName, "spec.0.source.0.type", "code"),
					resource.TestCheckResourceAttr(rName, "spec.0.source.0.sub_type", "GitHub"),
					resource.TestCheckResourceAttr(rName, "spec.0.source.0.url", acceptance.HW_CAE_CODE_URL),
					resource.TestCheckResourceAttr(rName, "spec.0.source.0.code.0.%", "3"),
					resource.TestCheckResourceAttr(rName, "spec.0.source.0.code.0.auth_name", acceptance.HW_CAE_CODE_AUTH_NAME),
					resource.TestCheckResourceAttr(rName, "spec.0.source.0.code.0.branch", acceptance.HW_CAE_CODE_BRANCH),
					resource.TestCheckResourceAttr(rName, "spec.0.source.0.code.0.namespace", acceptance.HW_CAE_CODE_NAMESPACE),
					resource.TestCheckResourceAttr(rName, "spec.0.build.0.archive.0.artifact_namespace", acceptance.HW_CAE_ARTIFACT_NAMESPACE),
					resource.TestCheckResourceAttr(rName, "spec.0.build.0.parameters.base_image", acceptance.HW_CAE_BUILD_BASE_IMAGE),
					resource.TestCheckResourceAttr(rName, "spec.0.build.0.parameters.dockerfile_path", "./Dockerfile"),
					resource.TestCheckResourceAttr(rName, "spec.0.build.0.parameters.build_cmd", "echo test"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metadata.0.annotations", "spec.0.build.0.parameters"},
				ImportStateIdFunc:       testAccComponentImportStateFunc(rName),
			},
		},
	})
}

func testAccComponent_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cae_component" "test" {
  environment_id = "%s"
  application_id = "%s"
  
  metadata {
    name = "%s"
    
    annotations = {
      version = "1.0.0"
    }
  }
  
  spec {
    replica = 2
    runtime = "Docker"

    source {
      type = "image"
      url  = "%s"
    }
  
    resource_limit {
      cpu    = "1000m"
      memory = "4Gi"
    }
  }
}
`, acceptance.HW_CAE_ENVIRONMENT_ID, acceptance.HW_CAE_APPLICATION_ID, name, acceptance.HW_CAE_IMAGE_URL)
}

func testAccComponent_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cae_component" "test" {
  environment_id = "%s"
  application_id = "%s"
  
  metadata {
    name = "%s"
    
    annotations = {
      version = "1.0.0"
    }
  }
  
  spec {
    replica = 1
    runtime = "Java17"

    source {
      type     = "code"
      sub_type = "GitHub"
      url      = "%s"
  
      code {
        auth_name = "%s"
        branch    = "%s"
        namespace = "%s"
      }
    }
  
    resource_limit {
      cpu    = "500m"
      memory = "1Gi"
    }
  
    build {
      archive {
        artifact_namespace = "%s"
      }

      parameters = {
        base_image      = "%s"
        dockerfile_path = "./Dockerfile"
        build_cmd       = "echo test"
      }
    }
  }
}
`, acceptance.HW_CAE_ENVIRONMENT_ID, acceptance.HW_CAE_APPLICATION_ID, name, acceptance.HW_CAE_CODE_URL,
		acceptance.HW_CAE_CODE_AUTH_NAME, acceptance.HW_CAE_CODE_BRANCH, acceptance.HW_CAE_CODE_NAMESPACE,
		acceptance.HW_CAE_ARTIFACT_NAMESPACE, acceptance.HW_CAE_BUILD_BASE_IMAGE)
}

func testAccComponentImportStateFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		var (
			environmentId = rs.Primary.Attributes["environment_id"]
			applicationId = rs.Primary.Attributes["application_id"]
			componentId   = rs.Primary.ID
		)
		if environmentId == "" || applicationId == "" || componentId == "" {
			return "", fmt.Errorf("some import IDs are missing, want '<environment_id>/<application_id>/<id>', but got '%s/%s/%s'",
				environmentId, applicationId, componentId)
		}

		return fmt.Sprintf("%s/%s/%s", environmentId, applicationId, componentId), nil
	}
}
