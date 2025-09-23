package aom

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/aom"
)

func getServiceDiscoveryRuleResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("aom", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating AOM client: %s", err)
	}

	// wait 30 seconds before listing the rule, to avoid error
	// lintignore:R018
	time.Sleep(30 * time.Second)

	return aom.GetServiceDiscoveryRule(client, state.Primary.ID)
}

func TestAccAOMServiceDiscoveryRule_basic(t *testing.T) {
	var ar interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := rName + "-update"
	resourceName := "huaweicloud_aom_service_discovery_rule.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&ar,
		getServiceDiscoveryRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAOMServiceDiscoveryRule_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "priority", "9999"),
					resource.TestCheckResourceAttr(resourceName, "detect_log_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "discovery_rule_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_default_rule", "true"),
					resource.TestCheckResourceAttr(resourceName, "log_file_suffix.0", "log"),
					resource.TestCheckResourceAttr(resourceName, "service_type", "Python"),
					resource.TestCheckResourceAttr(resourceName, "discovery_rules.0.check_content.0", "python"),
					resource.TestCheckResourceAttr(resourceName, "log_path_rules.0.args.0", "python"),
					resource.TestCheckResourceAttr(
						resourceName, "name_rules.0.service_name_rule.0.args.0", "python"),
					resource.TestCheckResourceAttr(
						resourceName, "name_rules.0.application_name_rule.0.args.0", "python"),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAOMServiceDiscoveryRule_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "priority", "9998"),
					resource.TestCheckResourceAttr(resourceName, "detect_log_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "discovery_rule_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_default_rule", "false"),
					resource.TestCheckResourceAttr(resourceName, "log_file_suffix.0", "trace"),
					resource.TestCheckResourceAttr(resourceName, "service_type", "Java"),
					resource.TestCheckResourceAttr(resourceName, "discovery_rules.0.check_content.0", "java"),
					resource.TestCheckResourceAttr(resourceName, "log_path_rules.0.args.0", "java"),
					resource.TestCheckResourceAttr(
						resourceName, "name_rules.0.service_name_rule.0.args.0", "java"),
					resource.TestCheckResourceAttr(
						resourceName, "name_rules.0.application_name_rule.0.args.0", "java"),
				),
			},
		},
	})
}

func testAOMServiceDiscoveryRule_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_aom_service_discovery_rule" "test" {
  name                   = "%s"
  priority               = 9999
  detect_log_enabled     = true
  discovery_rule_enabled = true
  is_default_rule        = true
  log_file_suffix        = ["log"]
  service_type           = "Python"
  description            = "test"

  discovery_rules {
    check_content = ["python"]
    check_mode    = "contain"
    check_type    = "cmdLine"
  }

  log_path_rules {
    name_type = "cmdLineHash"
    args      = ["python"]
    value     = ["/tmp/log"]
  }

  name_rules {
    service_name_rule {
      name_type = "str"
      args      = ["python"]
    }
    application_name_rule {
      name_type = "str"
      args      = ["python"]
    }
  }
}
`, rName)
}

func testAOMServiceDiscoveryRule_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_aom_service_discovery_rule" "test" {
  name                   = "%s"
  priority               = 9998
  detect_log_enabled     = false
  discovery_rule_enabled = false
  is_default_rule        = false
  log_file_suffix        = ["trace"]
  service_type           = "Java"

  discovery_rules {
    check_content = ["java"]
    check_mode    = "contain"
    check_type    = "cmdLine"
  }

  log_path_rules {
    name_type = "cmdLineHash"
    args      = ["java"]
    value     = ["/tmp/log"]
  }

  name_rules {
    service_name_rule {
      name_type = "str"
      args      = ["java"]
    }
    application_name_rule {
      name_type = "str"
      args      = ["java"]
    }
  }
}
`, rName)
}
