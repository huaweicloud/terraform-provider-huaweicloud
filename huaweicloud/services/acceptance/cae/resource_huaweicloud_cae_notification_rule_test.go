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

func getNotificationRuleFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("cae", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CAE client: %s", err)
	}

	return cae.GetEventNotificationRuleById(client, state.Primary.ID, state.Primary.Attributes["enterprise_project_id"])
}

func TestAccNotificationRule_basic(t *testing.T) {
	var (
		obj interface{}

		rNameWithEnv = "huaweicloud_cae_notification_rule.env"
		rcWithEnv    = acceptance.InitResourceCheck(rNameWithEnv, &obj, getNotificationRuleFunc)

		rNameWithApp = "huaweicloud_cae_notification_rule.app"
		rcWithApp    = acceptance.InitResourceCheck(rNameWithApp, &obj, getNotificationRuleFunc)

		rNameWithCom = "huaweicloud_cae_notification_rule.com"
		rcWithCom    = acceptance.InitResourceCheck(rNameWithCom, &obj, getNotificationRuleFunc)

		name = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCaeEnvironmentIds(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcWithEnv.CheckResourceDestroy(),
			rcWithApp.CheckResourceDestroy(),
			rcWithCom.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationRule_basic_step1(name),
				// For all components in the environment.
				Check: resource.ComposeTestCheckFunc(
					rcWithEnv.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameWithEnv, "name", name),
					resource.TestCheckResourceAttr(rNameWithEnv, "event_name", "Started,Healthy"),
					resource.TestCheckResourceAttr(rNameWithEnv, "scope.0.type", "environments"),
					resource.TestCheckResourceAttr(rNameWithEnv, "scope.0.environments.#", "1"),
					resource.TestCheckResourceAttrSet(rNameWithEnv, "scope.0.environments.0"),
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
					resource.TestCheckResourceAttr(rNameWithApp, "scope.0.applications.#", "2"),
					resource.TestCheckResourceAttrPair(rNameWithApp, "scope.0.applications.0", "huaweicloud_cae_application.test.0", "id"),
					resource.TestCheckResourceAttrPair(rNameWithApp, "scope.0.applications.1", "huaweicloud_cae_application.test.1", "id"),
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
				Config: testAccNotificationRule_basic_step2(name),
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
					resource.TestCheckResourceAttr(rNameWithApp, "scope.0.environments.#", "1"),
					resource.TestCheckResourceAttrPair(rNameWithApp, "scope.0.environments.0",
						"data.huaweicloud_cae_environments.test", "environments.0.id"),
					resource.TestCheckResourceAttr(rNameWithApp, "trigger_policy.0.type", "accumulative"),
					resource.TestCheckResourceAttr(rNameWithApp, "trigger_policy.0.period", "86400"),
					resource.TestCheckResourceAttr(rNameWithApp, "trigger_policy.0.count", "10"),
					resource.TestCheckResourceAttr(rNameWithApp, "trigger_policy.0.operator", ">"),
					resource.TestCheckResourceAttr(rNameWithApp, "enabled", "true"),
					// Updates from the specified components to all components in the application.
					rcWithCom.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameWithCom, "event_name", "FailedMount"),
					resource.TestCheckResourceAttr(rNameWithCom, "scope.0.type", "applications"),
					resource.TestCheckResourceAttr(rNameWithCom, "scope.0.applications.#", "2"),
					resource.TestCheckResourceAttrPair(rNameWithCom, "scope.0.applications.0", "huaweicloud_cae_application.test.0", "id"),
					resource.TestCheckResourceAttrPair(rNameWithCom, "scope.0.applications.1", "huaweicloud_cae_application.test.1", "id"),
					resource.TestCheckResourceAttr(rNameWithCom, "trigger_policy.0.type", "immediately"),
				),
			},
			{
				ResourceName:      rNameWithEnv,
				ImportState:       true,
				ImportStateVerify: true,
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

func testAccNotificationRule_base(name string) string {
	return fmt.Sprintf(`
locals {
  env_ids = split(",", "%[1]s")
}

# Query by environment ID under default enterprise project ID.
data "huaweicloud_cae_environments" "test" {
  environment_id = local.env_ids[0]
}

data "huaweicloud_swr_repositories" "test" {}

locals {
  swr_repositories = [for v in data.huaweicloud_swr_repositories.test.repositories : v if length(v.tags) > 0][0]
}

resource "huaweicloud_cae_application" "test" {
  count = 2

  environment_id = try(data.huaweicloud_cae_environments.test.environments[0].id, "NOT_FOUND")
  name           = format("%[2]s-%%d", count.index)
}

resource "huaweicloud_cae_component" "test" {
  count = 2

  environment_id = try(data.huaweicloud_cae_environments.test.environments[0].id, "NOT_FOUND")
  application_id = huaweicloud_cae_application.test[count.index].id

  metadata {
    name = format("%[2]s-%%d", count.index)

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
`, acceptance.HW_CAE_ENVIRONMENT_IDS, name)
}

func testAccNotificationRule_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

# For all components in the environment.
resource "huaweicloud_cae_notification_rule" "env" {
  name       = "%[2]s"
  event_name = "Started,Healthy"

  scope {
    type         = "environments"
    environments = data.huaweicloud_cae_environments.test.environments[*].id
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
    applications = huaweicloud_cae_application.test[*].id
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
`, testAccNotificationRule_base(name), name)
}

func testAccNotificationRule_basic_step2(name string) string {
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
    environments = data.huaweicloud_cae_environments.test.environments[*].id
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
    applications = huaweicloud_cae_application.test[*].id
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
`, testAccNotificationRule_base(name), name)
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

func TestAccNotificationRule_withEpsId(t *testing.T) {
	var (
		obj interface{}

		rNameWithEnvUnderEps = "huaweicloud_cae_notification_rule.env"
		rcWithEnvUnderEps    = acceptance.InitResourceCheck(rNameWithEnvUnderEps, &obj, getNotificationRuleFunc)

		name = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCaeEnvironmentIds(t, 2)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rcWithEnvUnderEps.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationRule_withEpsId_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rcWithEnvUnderEps.CheckResourceExists(),
				),
			},
		},
	})
}

func testAccNotificationRule_withEpsId_base() string {
	return fmt.Sprintf(`
locals {
  env_ids = split(",", "%[1]s")
}

# Query by environment ID under non-default enterprise project ID.
data "huaweicloud_cae_environments" "test" {
  environment_id = local.env_ids[1]
}
`, acceptance.HW_CAE_ENVIRONMENT_IDS)
}

func testAccNotificationRule_withEpsId_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

# For all components in the environment.
resource "huaweicloud_cae_notification_rule" "env" {
  name       = "%[2]s"
  event_name = "Started,Healthy"

  scope {
    type         = "environments"
    environments = data.huaweicloud_cae_environments.test.environments[*].id
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
`, testAccNotificationRule_withEpsId_base(), name)
}
