package cae

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cae"
)

func getTimerRuleFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("cae", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CAE client: %s", err)
	}

	return cae.GetTimerRuleById(
		client,
		state.Primary.Attributes["environment_id"],
		state.Primary.ID,
		state.Primary.Attributes["enterprise_project_id"],
	)
}

func getCronPrefix(isOneTime bool) string {
	currentTime := time.Now()
	// The triggered time of the rule must be at least two minutes later than the current time.
	nextTime := currentTime.AddDate(0, 0, 1)
	parsedTime := fmt.Sprintf("%d %d %d", nextTime.Second(), nextTime.Minute(), nextTime.Hour())
	if isOneTime {
		parsedTime = fmt.Sprintf("%s %d %d ? %d", parsedTime, nextTime.Day(), nextTime.Month(), nextTime.Year())
	}

	return parsedTime
}

func TestAccTimerRule_basic(t *testing.T) {
	var (
		obj interface{}

		name       = acceptance.RandomAccResourceNameWithDash()
		updateName = acceptance.RandomAccResourceNameWithDash()

		withEnv   = "huaweicloud_cae_timer_rule.env"
		withApp   = "huaweicloud_cae_timer_rule.app"
		withCom   = "huaweicloud_cae_timer_rule.com"
		withEnvRc = acceptance.InitResourceCheck(withEnv, &obj, getTimerRuleFunc)
		withAppRc = acceptance.InitResourceCheck(withApp, &obj, getTimerRuleFunc)
		withComRc = acceptance.InitResourceCheck(withCom, &obj, getTimerRuleFunc)

		withSpecifiedEpsId   = "huaweicloud_cae_timer_rule.specified_eps"
		withSpecifiedEpsIdRc = acceptance.InitResourceCheck(withSpecifiedEpsId, &obj, getTimerRuleFunc)

		oneTimeCron  = getCronPrefix(true)
		withDayCron  = fmt.Sprintf("%s ? * * *", getCronPrefix(false))
		withWeekCron = fmt.Sprintf("%s ? * 0,1,2,3,4,5,6 *", getCronPrefix(false))
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCaeEnvironmentIds(t, 2)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			withEnvRc.CheckResourceDestroy(),
			withAppRc.CheckResourceDestroy(),
			withComRc.CheckResourceDestroy(),
			withSpecifiedEpsIdRc.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccTimerRule_step1(name, oneTimeCron, withDayCron, withWeekCron),
				Check: resource.ComposeTestCheckFunc(
					// For all components in the environment.
					withEnvRc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(withEnv, "environment_id"),
					resource.TestCheckResourceAttr(withEnv, "name", name+"_env"),
					resource.TestCheckResourceAttr(withEnv, "type", "start"),
					resource.TestCheckResourceAttr(withEnv, "effective_range", "environment"),
					resource.TestCheckResourceAttr(withEnv, "effective_policy", "onetime"),
					resource.TestCheckResourceAttr(withEnv, "cron", oneTimeCron),
					resource.TestCheckResourceAttr(withEnv, "applications.#", "0"),
					resource.TestCheckResourceAttr(withEnv, "components.#", "0"),
					// For all components in the application.
					withAppRc.CheckResourceExists(),
					resource.TestCheckResourceAttr(withApp, "name", name+"_app"),
					resource.TestCheckResourceAttr(withApp, "type", "stop"),
					resource.TestCheckResourceAttr(withApp, "effective_range", "application"),
					resource.TestCheckResourceAttr(withApp, "effective_policy", "periodic"),
					resource.TestCheckResourceAttr(withApp, "cron", withDayCron),
					resource.TestCheckResourceAttr(withApp, "applications.#", "1"),
					resource.TestCheckResourceAttrPair(withApp, "applications.0.id", "huaweicloud_cae_application.test", "id"),
					resource.TestCheckResourceAttrPair(withApp, "applications.0.name", "huaweicloud_cae_application.test", "name"),
					resource.TestCheckResourceAttr(withApp, "components.#", "0"),
					// For specified components.
					withComRc.CheckResourceExists(),
					resource.TestCheckResourceAttr(withCom, "name", name+"_com"),
					resource.TestCheckResourceAttr(withCom, "type", "start"),
					resource.TestCheckResourceAttr(withCom, "effective_range", "component"),
					resource.TestCheckResourceAttr(withCom, "effective_policy", "periodic"),
					resource.TestCheckResourceAttr(withCom, "cron", withWeekCron),
					resource.TestCheckResourceAttr(withCom, "applications.#", "0"),
					resource.TestCheckResourceAttr(withCom, "components.#", "2"),
					resource.TestCheckResourceAttrSet(withCom, "components.0.id"),
					resource.TestCheckResourceAttrSet(withCom, "components.0.name"),
					resource.TestCheckNoResourceAttr(withCom, "enterprise_project_id"),
					// For specified enterprise project.
					withSpecifiedEpsIdRc.CheckResourceExists(),
					resource.TestCheckResourceAttr(withSpecifiedEpsId, "name", name+"_specified_eps"),
					resource.TestCheckResourceAttrSet(withSpecifiedEpsId, "environment_id"),
					resource.TestCheckResourceAttrPair(withSpecifiedEpsId, "enterprise_project_id", "data.huaweicloud_cae_environments.test",
						"environments.0.annotations.enterprise_project_id"),
				),
			},
			{
				Config: testAccTimerRule_step2(name, updateName, oneTimeCron, withDayCron, withWeekCron),
				Check: resource.ComposeTestCheckFunc(
					// Update from all components in the environment to all components in the application.
					withEnvRc.CheckResourceExists(),
					resource.TestCheckResourceAttr(withEnv, "name", updateName+"_app"),
					resource.TestCheckResourceAttr(withEnv, "type", "stop"),
					resource.TestCheckResourceAttr(withEnv, "effective_range", "application"),
					resource.TestCheckResourceAttr(withEnv, "effective_policy", "periodic"),
					resource.TestCheckResourceAttr(withEnv, "cron", withWeekCron),
					resource.TestCheckResourceAttr(withEnv, "applications.#", "1"),
					resource.TestCheckResourceAttrPair(withEnv, "applications.0.id", "huaweicloud_cae_application.test", "id"),
					resource.TestCheckResourceAttr(withEnv, "components.#", "0"),
					// Update from all components in the application to specified components.
					withAppRc.CheckResourceExists(),
					resource.TestCheckResourceAttr(withApp, "name", updateName+"_com"),
					resource.TestCheckResourceAttr(withApp, "type", "stop"),
					resource.TestCheckResourceAttr(withApp, "effective_range", "component"),
					resource.TestCheckResourceAttr(withApp, "effective_policy", "onetime"),
					resource.TestCheckResourceAttr(withApp, "cron", oneTimeCron),
					resource.TestCheckResourceAttr(withApp, "components.#", "1"),
					resource.TestCheckResourceAttrPair(withApp, "components.0.id", "huaweicloud_cae_component.test.0", "id"),
					resource.TestCheckResourceAttr(withApp, "components.0.name", ""),
					resource.TestCheckResourceAttr(withApp, "applications.#", "0"),
					// Update from specified components to all components in the environment.
					withComRc.CheckResourceExists(),
					resource.TestCheckResourceAttr(withCom, "name", updateName+"_env"),
					resource.TestCheckResourceAttr(withCom, "type", "stop"),
					resource.TestCheckResourceAttr(withCom, "effective_range", "environment"),
					resource.TestCheckResourceAttr(withCom, "effective_policy", "periodic"),
					resource.TestCheckResourceAttr(withCom, "cron", withDayCron),
					resource.TestCheckResourceAttr(withCom, "components.#", "0"),
					resource.TestCheckResourceAttr(withCom, "applications.#", "0"),
					// For specified enterprise project.
					withSpecifiedEpsIdRc.CheckResourceExists(),
					resource.TestCheckResourceAttr(withSpecifiedEpsId, "name", updateName+"_specified_eps"),
					resource.TestCheckResourceAttrSet(withSpecifiedEpsId, "environment_id"),
					resource.TestCheckResourceAttrPair(withSpecifiedEpsId, "enterprise_project_id", "data.huaweicloud_cae_environments.test",
						"environments.0.annotations.enterprise_project_id"),
				),
			},
			{
				ResourceName:            withEnv,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccTimerRuleImportStateFunc(withEnv, true),
				ImportStateVerifyIgnore: []string{"status"},
			},
			{
				ResourceName:            withApp,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccTimerRuleImportStateFunc(withApp, true),
				ImportStateVerifyIgnore: []string{"status"},
			},
			{
				ResourceName:            withCom,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccTimerRuleImportStateFunc(withCom, true),
				ImportStateVerifyIgnore: []string{"status"},
			},
			{
				ResourceName:            withSpecifiedEpsId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccTimerRuleImportStateFunc(withSpecifiedEpsId, false),
				ImportStateVerifyIgnore: []string{"status"},
			},
		},
	})
}

func testAccResourceTimerRule_base(name string) string {
	return fmt.Sprintf(`
locals {
  env_id_with_default_eps_id   = split(",", "%[1]s")[0]
  env_id_with_specified_eps_id = split(",", "%[1]s")[1]
}

data "huaweicloud_cae_environments" "test" {
  environment_id = local.env_id_with_specified_eps_id
}

resource "huaweicloud_cae_application" "test" {
  environment_id = local.env_id_with_default_eps_id
  name           = "%[2]s"
}

data "huaweicloud_swr_repositories" "test" {}

locals {
  swr_repositories = [for v in data.huaweicloud_swr_repositories.test.repositories : v if length(v.tags) > 0][0]
}

resource "huaweicloud_cae_component" "test" {
  count          = 2
  environment_id = local.env_id_with_default_eps_id
  application_id = huaweicloud_cae_application.test.id

  metadata {
    name = "%[2]s${count.index}"

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

func testAccTimerRule_step1(name, oneTimeCron, withDayCron, withWeekCron string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cae_timer_rule" "env" {
  environment_id   = local.env_id_with_default_eps_id
  name             = "%[2]s_env"
  type             = "start"
  status           = "on"
  effective_range  = "environment"
  effective_policy = "onetime"
  cron             = "%[3]s"
}

resource "huaweicloud_cae_timer_rule" "app" {
  environment_id   = local.env_id_with_default_eps_id
  name             = "%[2]s_app"
  type             = "stop"
  status           = "off"
  effective_range  = "application"
  effective_policy = "periodic"
  cron             = "%[4]s"

  applications {
    id   = huaweicloud_cae_application.test.id
    name = huaweicloud_cae_application.test.name
  }
}

resource "huaweicloud_cae_timer_rule" "com" {
  environment_id   = local.env_id_with_default_eps_id
  name             = "%[2]s_com"
  type             = "start"
  status           = "off"
  effective_range  = "component"
  effective_policy = "periodic"
  cron             = "%[5]s"

  dynamic "components" {
    for_each = huaweicloud_cae_component.test[*]
    content {
      id   = components.value.id
      name = components.value.metadata[0].name
    }
  }
}

resource "huaweicloud_cae_timer_rule" "specified_eps" {
  environment_id        = local.env_id_with_specified_eps_id
  name                  = "%[2]s_specified_eps"
  type                  = "stop"
  status                = "on"
  effective_range       = "environment"
  effective_policy      = "onetime"
  cron                  = "%[3]s"
  enterprise_project_id = data.huaweicloud_cae_environments.test.environments[0].annotations.enterprise_project_id
}
`, testAccResourceTimerRule_base(name), name, oneTimeCron, withDayCron, withWeekCron)
}

func testAccTimerRule_step2(name, updateName, oneTimeCron, withDayCron, withWeekCron string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cae_timer_rule" "env" {
  environment_id   = local.env_id_with_default_eps_id
  name             = "%[2]s_app"
  type             = "stop"
  status           = "off"
  effective_range  = "application"
  effective_policy = "periodic"
  cron             = "%[5]s"

  applications {
    id = huaweicloud_cae_application.test.id
  }
}

resource "huaweicloud_cae_timer_rule" "app" {
  environment_id   = local.env_id_with_default_eps_id
  name             = "%[2]s_com"
  type             = "stop"
  status           = "on"
  effective_range  = "component"
  effective_policy = "onetime"
  cron             = "%[3]s"

  components {
    id   = huaweicloud_cae_component.test[0].id
  }
}

resource "huaweicloud_cae_timer_rule" "com" {
  environment_id   = local.env_id_with_default_eps_id
  name             = "%[2]s_env"
  type             = "stop"
  status           = "on"
  effective_range  = "environment"
  effective_policy = "periodic"
  cron             = "%[4]s"
}

resource "huaweicloud_cae_timer_rule" "specified_eps" {
  environment_id        = local.env_id_with_specified_eps_id
  name                  = "%[2]s_specified_eps"
  type                  = "stop"
  status                = "on"
  effective_range       = "environment"
  effective_policy      = "periodic"
  cron                  = "%[4]s"
  enterprise_project_id = data.huaweicloud_cae_environments.test.environments[0].annotations.enterprise_project_id
}
`, testAccResourceTimerRule_base(name), updateName, oneTimeCron, withDayCron, withWeekCron)
}

func testAccTimerRuleImportStateFunc(name string, isDefaultEps bool) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		var (
			environmentId = rs.Primary.Attributes["environment_id"]
			timerRuleName = rs.Primary.Attributes["name"]
		)
		if environmentId == "" || timerRuleName == "" {
			return "", fmt.Errorf("some import IDs are missing, want '<environment_id>/<name>' or "+
				"'<environment_id>/<name>/<enterprise_project_id>', but got '%s/%s'",
				environmentId, timerRuleName)
		}

		if isDefaultEps {
			return fmt.Sprintf("%s/%s", environmentId, timerRuleName), nil
		}

		epsId := rs.Primary.Attributes["enterprise_project_id"]
		if epsId == "" {
			return "", fmt.Errorf("some import IDs are missing, want '<environment_id>/<name>/<enterprise_project_id>', but got '%s/%s'",
				environmentId, timerRuleName)
		}

		return fmt.Sprintf("%s/%s/%s", environmentId, timerRuleName, epsId), nil
	}
}
