package secmaster

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

func getAlertRuleResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getAlertRule: Query the SecMaster alert rule detail
	var (
		getAlertRuleHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/siem/alert-rules/{id}"
		getAlertRuleProduct = "secmaster"
	)
	getAlertRuleClient, err := cfg.NewServiceClient(getAlertRuleProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	getAlertRulePath := getAlertRuleClient.Endpoint + getAlertRuleHttpUrl
	getAlertRulePath = strings.ReplaceAll(getAlertRulePath, "{project_id}", getAlertRuleClient.ProjectID)
	getAlertRulePath = strings.ReplaceAll(getAlertRulePath, "{workspace_id}", state.Primary.Attributes["workspace_id"])
	getAlertRulePath = strings.ReplaceAll(getAlertRulePath, "{id}", state.Primary.ID)

	getAlertRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getAlertRuleResp, err := getAlertRuleClient.Request("GET", getAlertRulePath, &getAlertRuleOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getAlertRuleResp)
}

func TestAccAlertRule_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_secmaster_alert_rule.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAlertRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMaster(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAlertRule_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "pipeline_id", acceptance.HW_SECMASTER_PIPELINE_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "severity", "TIPS"),
					resource.TestCheckResourceAttr(rName, "description", "this is a test rule created by terraform"),
					resource.TestCheckResourceAttr(rName, "status", "ENABLED"),
					resource.TestCheckResourceAttr(rName, "query_rule", "* | select status, count(*) as count group by status"),
					resource.TestCheckResourceAttr(rName, "query_plan.0.query_interval", "1"),
					resource.TestCheckResourceAttr(rName, "query_plan.0.query_interval_unit", "HOUR"),
					resource.TestCheckResourceAttr(rName, "query_plan.0.time_window", "1"),
					resource.TestCheckResourceAttr(rName, "query_plan.0.time_window_unit", "HOUR"),
					resource.TestCheckResourceAttr(rName, "triggers.0.operator", "GT"),
					resource.TestCheckResourceAttr(rName, "triggers.0.expression", "5"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testAlertRule_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "pipeline_id", acceptance.HW_SECMASTER_PIPELINE_ID),
					resource.TestCheckResourceAttr(rName, "name", name+"_update"),
					resource.TestCheckResourceAttr(rName, "severity", "MEDIUM"),
					resource.TestCheckResourceAttr(rName, "description", "this is a test rule created by terraform update"),
					resource.TestCheckResourceAttr(rName, "status", "DISABLED"),
					resource.TestCheckResourceAttr(rName, "query_rule", "* | select status, count(*) as count group by status"),
					resource.TestCheckResourceAttr(rName, "query_plan.0.query_interval", "5"),
					resource.TestCheckResourceAttr(rName, "query_plan.0.query_interval_unit", "MINUTE"),
					resource.TestCheckResourceAttr(rName, "query_plan.0.time_window", "5"),
					resource.TestCheckResourceAttr(rName, "query_plan.0.time_window_unit", "MINUTE"),
					resource.TestCheckResourceAttr(rName, "triggers.0.operator", "EQ"),
					resource.TestCheckResourceAttr(rName, "triggers.0.expression", "10"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAlertRuleImportState(rName),
				ImportStateVerifyIgnore: []string{
					"type",
				},
			},
		},
	})
}

func testAlertRule_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_alert_rule" "test" {
  workspace_id = "%s"
  pipeline_id  = "%s"
  name         = "%s"
  description  = "this is a test rule created by terraform"
  status       = "ENABLED"
  severity     = "TIPS"

  type = {
    "name"     = "DNS protocol attacks"
    "category" = "DDoS attacks"
  }

  triggers {
    mode              = "COUNT"
    operator          = "GT"
    expression        = 5
    severity          = "MEDIUM"
    accumulated_times = 1
  }

  query_rule = "* | select status, count(*) as count group by status"
  query_type = "SQL"
  query_plan {
    query_interval      = 1
    query_interval_unit = "HOUR"
    time_window         = 1
    time_window_unit    = "HOUR"
  }
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_PIPELINE_ID, name)
}

func testAlertRule_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_alert_rule" "test" {
  workspace_id = "%s"
  pipeline_id  = "%s"
  name         = "%s_update"
  description  = "this is a test rule created by terraform update"
  status       = "DISABLED"
  severity     = "MEDIUM"

  type = {
    "name"     = "DNS protocol attacks"
    "category" = "DDoS attacks"
  }

  triggers {
    mode              = "COUNT"
    operator          = "EQ"
    expression        = 10
    severity          = "MEDIUM"
    accumulated_times = 1
  }

  query_rule = "* | select status, count(*) as count group by status"
  query_type = "SQL"
  query_plan {
    query_interval      = 5
    query_interval_unit = "MINUTE"
    time_window         = 5
    time_window_unit    = "MINUTE"
  }
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_PIPELINE_ID, name)
}

func testAlertRuleImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["workspace_id"] == "" {
			return "", fmt.Errorf("attribute (workspace_id) of resource (%s) not found: %s", name, rs)
		}

		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["workspace_id"], rs.Primary.ID), nil
	}
}
