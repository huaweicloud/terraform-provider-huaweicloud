package aom

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/aom"
)

func getRecordingRuleResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("aom", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating AOM client: %s", err)
	}

	return aom.GetRecordingRuleByInstanceId(client, state.Primary.Attributes["instance_id"])
}

func TestAccRecordingRule_basic(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_aom_recording_rule.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getRecordingRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAomPrometheusInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// Avoid CheckDestroy, because there is nothing in the resource destroy method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testRecordingRule_basic_step1(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttrSet(resourceName, "recording_rule"),
				),
			},
			{
				Config: testRecordingRule_basic_step2(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttrSet(resourceName, "recording_rule"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testRecordingRuleImportStateWithInstanceId(resourceName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testRecordingRuleImportStateWithRuleId(resourceName),
			},
		},
	})
}

func testRecordingRule_basic_step1() string {
	return fmt.Sprintf(`
resource "huaweicloud_aom_recording_rule" "test" {
  instance_id    = "%[1]s"
  recording_rule = <<EOF
groups:
  - name: node_basic_aggregation
    interval: 60s
    rules:
      - record: instance:node_cpu_usage:percent_avg5m
        expr: 100 - (avg by(instance) (rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)
EOF
}
`, acceptance.HW_AOM_PROMETHEUS_INSTANCE_ID)
}

func testRecordingRule_basic_step2() string {
	return fmt.Sprintf(`
resource "huaweicloud_aom_recording_rule" "test" {
  instance_id    = "%[1]s"
  recording_rule = <<EOF
groups:
  - name: node_basic_aggregation
    interval: 60s
    rules:
      - record: instance:node_memory_usage:percent
        expr: (node_memory_MemTotal_bytes - node_memory_MemAvailable_bytes) / node_memory_MemTotal_bytes * 100
EOF
}
`, acceptance.HW_AOM_PROMETHEUS_INSTANCE_ID)
}

func testRecordingRuleImportStateWithInstanceId(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}

		instanceId := rs.Primary.Attributes["instance_id"]
		if instanceId == "" {
			return "", fmt.Errorf("attribute (instance_id) of resource (%s) not found", name)
		}

		return instanceId, nil
	}
}

func testRecordingRuleImportStateWithRuleId(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}

		instanceId := rs.Primary.Attributes["instance_id"]
		if instanceId == "" {
			return "", fmt.Errorf("attribute (instance_id) of resource (%s) not found", name)
		}

		ruleId := rs.Primary.ID
		if ruleId == "" {
			return "", fmt.Errorf("attribute ID of resource (%s) not found", name)
		}

		return fmt.Sprintf("%s/%s", instanceId, ruleId), nil
	}
}
