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

func getResourceCaeNotificationRuleFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("cae", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CAE client: %s", err)
	}

	return cae.GetEventNotificationRuleById(client, state.Primary.ID)
}

func TestAccResourceNotificationRule_basic(t *testing.T) {
	var (
		rNameWithEnv = "huaweicloud_cae_notification_rule.env"
		rNameWithApp = "huaweicloud_cae_notification_rule.app"
		rNameWithCom = "huaweicloud_cae_notification_rule.com"
		name         = acceptance.RandomAccResourceNameWithDash()
		obj          interface{}
		rcWithEnv    = acceptance.InitResourceCheck(rNameWithEnv, &obj, getResourceCaeNotificationRuleFunc)
		rcWithApp    = acceptance.InitResourceCheck(rNameWithApp, &obj, getResourceCaeNotificationRuleFunc)
		rcWithCom    = acceptance.InitResourceCheck(rNameWithCom, &obj, getResourceCaeNotificationRuleFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCaeEnvironment(t)
			acceptance.TestAccPreCheckCaeApplication(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rcWithEnv.CheckResourceDestroy(),

		Steps: []resource.TestStep{
			{
				Config: testAccResourceNotificationRule_basic_step1(name),
				// For all components in the environment.
				Check: resource.ComposeTestCheckFunc(
					rcWithEnv.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameWithEnv, "name", name),
					resource.TestCheckResourceAttr(rNameWithEnv, "event_name", "Started,Healthy"),
					resource.TestCheckResourceAttr(rNameWithEnv, "scope.0.type", "environments"),
					resource.TestCheckResourceAttr(rNameWithEnv, "scope.0.environments.0", acceptance.HW_CAE_ENVIRONMENT_ID),
					resource.TestCheckResourceAttr(rNameWithEnv, "trigger_policy.0.type", "accumulative"),
					resource.TestCheckResourceAttr(rNameWithEnv, "trigger_policy.0.period", "300"),
					resource.TestCheckResourceAttr(rNameWithEnv, "trigger_policy.0.count", "50"),
					resource.TestCheckResourceAttr(rNameWithEnv, "trigger_policy.0.operator", ">="),
					resource.TestCheckResourceAttr(rNameWithEnv, "notification.0.protocol", "email"),
					resource.TestCheckResourceAttr(rNameWithEnv, "notification.0.endpoint", "terraform@test.com"),
					resource.TestCheckResourceAttr(rNameWithEnv, "notification.0.template", "EN"),
					resource.TestCheckResourceAttr(rNameWithEnv, "enabled", "true"),
					// For all components in the application.
					rcWithApp.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameWithApp, "event_name", "Unhealthy"),
					resource.TestCheckResourceAttr(rNameWithApp, "scope.0.type", "applications"),
					resource.TestCheckResourceAttr(rNameWithApp, "scope.0.applications.0", acceptance.HW_CAE_APPLICATION_ID),
					resource.TestCheckResourceAttr(rNameWithApp, "trigger_policy.0.type", "immediately"),
					resource.TestCheckResourceAttr(rNameWithApp, "notification.0.protocol", "sms"),
					resource.TestCheckResourceAttr(rNameWithApp, "notification.0.endpoint", "12345678987"),
					resource.TestCheckResourceAttr(rNameWithApp, "notification.0.template", "ZH"),
					resource.TestCheckResourceAttr(rNameWithApp, "enabled", "false"),
					// For specified components.
					rcWithCom.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameWithCom, "event_name", "Started,BackOffStart,SuccessfulMountVolume"),
					resource.TestCheckResourceAttr(rNameWithCom, "scope.0.type", "components"),
					resource.TestCheckResourceAttr(rNameWithCom, "scope.0.components.#", "2"),
					resource.TestCheckResourceAttr(rNameWithCom, "trigger_policy.0.type", "accumulative"),
					resource.TestCheckResourceAttr(rNameWithCom, "trigger_policy.0.period", "300"),
					resource.TestCheckResourceAttr(rNameWithCom, "trigger_policy.0.count", "100"),
					resource.TestCheckResourceAttr(rNameWithCom, "trigger_policy.0.operator", ">"),
					resource.TestCheckResourceAttr(rNameWithCom, "notification.0.protocol", "email"),
					resource.TestCheckResourceAttr(rNameWithCom, "notification.0.endpoint", "terraform@test.com"),
					resource.TestCheckResourceAttr(rNameWithCom, "notification.0.template", "ZH"),
				),
			},
			{
				Config: testAccResourceNotificationRule_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					// Updates from all components in the environment to the specified components.
					rcWithEnv.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameWithEnv, "event_name", "FailedPullImage,FailedMount"),
					resource.TestCheckResourceAttr(rNameWithEnv, "scope.0.type", "components"),
					resource.TestCheckResourceAttr(rNameWithEnv, "scope.0.components.#", "2"),
					resource.TestCheckResourceAttr(rNameWithEnv, "trigger_policy.0.type", "immediately"),
					resource.TestCheckResourceAttr(rNameWithEnv, "enabled", "false"),
					// Updates from all components in the application to all components in the environments.
					rcWithApp.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameWithApp, "event_name", "FailedPullImage"),
					resource.TestCheckResourceAttr(rNameWithApp, "scope.0.type", "environments"),
					resource.TestCheckResourceAttr(rNameWithApp, "scope.0.environments.0", acceptance.HW_CAE_ENVIRONMENT_ID),
					resource.TestCheckResourceAttr(rNameWithApp, "trigger_policy.0.type", "accumulative"),
					resource.TestCheckResourceAttr(rNameWithApp, "trigger_policy.0.period", "86400"),
					resource.TestCheckResourceAttr(rNameWithApp, "trigger_policy.0.count", "10"),
					resource.TestCheckResourceAttr(rNameWithApp, "trigger_policy.0.operator", ">"),
					resource.TestCheckResourceAttr(rNameWithApp, "enabled", "true"),
					// Updates from the specified components to all components in the application.
					rcWithCom.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameWithCom, "event_name", "FailedMount"),
					resource.TestCheckResourceAttr(rNameWithCom, "scope.0.type", "applications"),
					resource.TestCheckResourceAttr(rNameWithCom, "scope.0.applications.0", acceptance.HW_CAE_APPLICATION_ID),
					resource.TestCheckResourceAttr(rNameWithCom, "trigger_policy.0.type", "immediately"),
				),
			},
			{
				ResourceName:      rNameWithEnv,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccNotificationRuleImportFunc(rNameWithEnv),
			},
			{
				ResourceName:      rNameWithApp,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccNotificationRuleImportFunc(rNameWithApp),
			},
			{
				ResourceName:      rNameWithCom,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccNotificationRuleImportFunc(rNameWithCom),
			},
		},
	})
}

func testAccResourceNotificationRule_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

# For all components in the environment.
resource "huaweicloud_cae_notification_rule" "env" {
  name       = "%[2]s"
  event_name = "Started,Healthy"

  scope {
    type         = "environments"
    environments = ["%[3]s"]
  }

  trigger_policy {
    type     = "accumulative"
    period   = 300
    count    = 50
    operator = ">="
  }

  notification {
    protocol = "email"
    endpoint = "terraform@test.com"
    template = "EN"
  }

  enabled = true
}

# For all components in the application.
resource "huaweicloud_cae_notification_rule" "app" {
  name       = "%[2]s-app"
  event_name = "Unhealthy"

  scope {
    type         = "applications"
    applications = ["%[4]s"]
  }

  trigger_policy {
    type = "immediately"
  }

  notification {
    protocol = "sms"
    endpoint = "12345678987"
    template = "ZH"
  }
}

# For specified components.
resource "huaweicloud_cae_notification_rule" "com" {
  name       = "%[2]s-com"
  event_name = "Started,BackOffStart,SuccessfulMountVolume"

  scope {
    type       = "components"
    components = huaweicloud_cae_component.test[*].id
  }

  trigger_policy {
    type     = "accumulative"
    period   = 300
    count    = 100
    operator = ">"
  }

  notification {
    protocol = "email"
    endpoint = "terraform@test.com"
    template = "ZH"
  }
}
`, testAccResourceNotificationRule_base(name), name, acceptance.HW_CAE_ENVIRONMENT_ID, acceptance.HW_CAE_APPLICATION_ID)
}

func testAccResourceNotificationRule_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

# Updates from all components in the environment to the specified components.
 resource "huaweicloud_cae_notification_rule" "env" {
  name       = "%[2]s"
  event_name = "FailedPullImage,FailedMount"

  scope {
    type       = "components"
    components = huaweicloud_cae_component.test[*].id
  }

  trigger_policy {
    type = "immediately"
  }

  notification {
    protocol = "email"
    endpoint = "terraform@test.com"
    template = "EN"
  }

  enabled = false
}

# Updates from all components in the application to all components in the environments.
resource "huaweicloud_cae_notification_rule" "app" {
  name       = "%[2]s-app"
  event_name = "FailedPullImage"

  scope {
    type         = "environments"
    environments = ["%[4]s"]
  }

  trigger_policy {
    type     = "accumulative"
    period   = 86400
    count    = 10
    operator = ">"
  }

  notification {
    protocol = "sms"
    endpoint = "12345678987"
    template = "ZH"
  }

  enabled = true
}

# Updates from the specified components to the application.
resource "huaweicloud_cae_notification_rule" "com" {
  name       = "%[2]s-com"
  event_name = "FailedMount"

  scope {
    type         = "applications"
    applications = ["%[3]s"]
  }

  trigger_policy {
    type = "immediately"
  }

  notification {
    protocol = "email"
    endpoint = "terraform@test.com"
    template = "ZH"
  }
}
`, testAccResourceNotificationRule_base(name), name, acceptance.HW_CAE_APPLICATION_ID, acceptance.HW_CAE_ENVIRONMENT_ID)
}

func testAccResourceNotificationRule_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_swr_repositories" "test" {}

locals {
  swr_repositories = [for v in data.huaweicloud_swr_repositories.test.repositories : v if length(v.tags) > 0][0]
}

resource "huaweicloud_cae_component" "test" {
  count          = 2
  environment_id = "%[1]s"
  application_id = "%[2]s"

  metadata {
    name = "%[3]s${count.index}"

    annotations = {
      version = "1.0.0"
    }
  }

  spec {
    replica = 1
    runtime = "Docker"

    source {
      type = "image"
      url  = format("%%s:%%s", local.swr_repositories.path, local.swr_repositories.tags[0])
    }

    resource_limit {
      cpu    = "1000m"
      memory = "4Gi"
    }
  }
}
`, acceptance.HW_CAE_ENVIRONMENT_ID, acceptance.HW_CAE_APPLICATION_ID, name)
}
func testAccNotificationRuleImportFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		name := rs.Primary.Attributes["name"]
		if name == "" {
			return "", fmt.Errorf("attribute (name) of resource (%s) not found: %s", name, rs)
		}
		return name, nil
	}
}
