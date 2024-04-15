package cae

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccComponentDeployment_basic(t *testing.T) {
	baseConfig := testAccComponentConfiguration_base()
	// Avoid CheckDestroy because this resource is a one-time action resource and there is nothing in the destroy
	// method.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCaeEnvironment(t)
			acceptance.TestAccPreCheckCaeApplication(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				// One-time action resource do not need to be checked and no processing is performed in the Read method.
				Config: testAccComponentDeployment_basic_step1(baseConfig),
			},
			{
				Config: testAccComponentDeployment_basic_step2(baseConfig),
			},
			{
				Config: testAccComponentDeployment_basic_step3(baseConfig),
			},
			{
				Config: testAccComponentDeployment_basic_step4(baseConfig),
			},
		},
	})
}

func testAccComponentDeployment_basic_step1(baseConfig string) string {
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
}

resource "huaweicloud_cae_component_deployment" "test" {
  environment_id = "%[2]s"
  application_id = "%[3]s"
  component_id   = huaweicloud_cae_component.test.id

  metadata {
    name = "deploy"
  }

  lifecycle {
    replace_triggered_by = [
      huaweicloud_cae_component_configurations.test.items
    ]
  }
}
`, baseConfig, acceptance.HW_CAE_ENVIRONMENT_ID, acceptance.HW_CAE_APPLICATION_ID)
}

func testAccComponentDeployment_basic_step2(baseConfig string) string {
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

resource "huaweicloud_cae_component_deployment" "test" {
  environment_id = "%[2]s"
  application_id = "%[3]s"
  component_id   = huaweicloud_cae_component.test.id

  metadata {
    name = "configure"
  }

  lifecycle {
    replace_triggered_by = [
      huaweicloud_cae_component_configurations.test.items
    ]
  }
}
`, baseConfig, acceptance.HW_CAE_ENVIRONMENT_ID, acceptance.HW_CAE_APPLICATION_ID)
}

func testAccComponentDeployment_basic_step3(baseConfig string) string {
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
            "foo": "baar"
        }
      }
    })
  }
}

resource "huaweicloud_cae_component_deployment" "test" {
  environment_id = "%[2]s"
  application_id = "%[3]s"
  component_id   = huaweicloud_cae_component.test.id

  metadata {
    name = "configure"
  }

  lifecycle {
    replace_triggered_by = [
      huaweicloud_cae_component_configurations.test.items
    ]
  }
}
`, baseConfig, acceptance.HW_CAE_ENVIRONMENT_ID, acceptance.HW_CAE_APPLICATION_ID)
}

func testAccComponentDeployment_basic_step4(baseConfig string) string {
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
            "foo": "baar"
        }
      }
    })
  }
}

resource "huaweicloud_cae_component_deployment" "test" {
  environment_id = "%[2]s"
  application_id = "%[3]s"
  component_id   = huaweicloud_cae_component.test.id

  metadata {
    name = "stop"
  }

  lifecycle {
    replace_triggered_by = [
      huaweicloud_cae_component_configurations.test.items
    ]
  }
}
`, baseConfig, acceptance.HW_CAE_ENVIRONMENT_ID, acceptance.HW_CAE_APPLICATION_ID)
}
