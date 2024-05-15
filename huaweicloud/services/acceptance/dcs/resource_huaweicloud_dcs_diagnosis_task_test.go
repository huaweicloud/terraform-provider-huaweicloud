package dcs

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

func getDiagnosisTaskResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	var (
		getDiagnosisTaskHttpUrl = "v2/{project_id}/diagnosis/{report_id}"
		getDiagnosisTaskProduct = "dcs"
	)
	getDiagnosisTaskClient, err := cfg.NewServiceClient(getDiagnosisTaskProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DCS client: %s", err)
	}

	getDiagnosisTaskPath := getDiagnosisTaskClient.Endpoint + getDiagnosisTaskHttpUrl
	getDiagnosisTaskPath = strings.ReplaceAll(getDiagnosisTaskPath, "{project_id}", getDiagnosisTaskClient.ProjectID)
	getDiagnosisTaskPath = strings.ReplaceAll(getDiagnosisTaskPath, "{report_id}", state.Primary.ID)
	getDiagnosisTaskOpt := golangsdk.RequestOpts{KeepResponseBody: true}

	getDiagnosisTaskResp, err := getDiagnosisTaskClient.Request("GET", getDiagnosisTaskPath, &getDiagnosisTaskOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving diagnosis report: %s", err)
	}

	getDiagnosisTaskRespBody, err := utils.FlattenResponse(getDiagnosisTaskResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving diagnosis report: %s", err)
	}
	return getDiagnosisTaskRespBody, nil
}

func TestAccDiagnosisTask_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_dcs_diagnosis_task.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDiagnosisTaskResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDCSInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDiagnosisTask_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "instance_id"),
					resource.TestCheckResourceAttrSet(rName, "node_ip_list.#"),
					resource.TestCheckResourceAttrSet(rName, "begin_time"),
					resource.TestCheckResourceAttrSet(rName, "end_time"),
					resource.TestCheckResourceAttrSet(rName, "abnormal_item_sum"),
					resource.TestCheckResourceAttrSet(rName, "failed_item_sum"),
					resource.TestCheckResourceAttrSet(rName, "diagnosis_node_report_list.0.abnormal_sum"),
					resource.TestCheckResourceAttrSet(rName, "diagnosis_node_report_list.0.az_code"),
					resource.TestCheckResourceAttrSet(rName, "diagnosis_node_report_list.0.failed_sum"),
					resource.TestCheckResourceAttrSet(rName, "diagnosis_node_report_list.0.group_name"),
					resource.TestCheckResourceAttrSet(rName, "diagnosis_node_report_list.0.is_faulted"),
					resource.TestCheckResourceAttrSet(rName, "diagnosis_node_report_list.0.node_ip"),
					resource.TestCheckResourceAttrSet(rName, "diagnosis_node_report_list.0.role"),
					resource.TestCheckResourceAttrSet(rName, "diagnosis_node_report_list.0.diagnosis_dimension_list.0.abnormal_num"),
					resource.TestCheckResourceAttrSet(rName, "diagnosis_node_report_list.0.diagnosis_dimension_list.0.failed_num"),
					resource.TestCheckResourceAttrSet(rName, "diagnosis_node_report_list.0.diagnosis_dimension_list.0.name"),
					resource.TestCheckResourceAttrSet(rName,
						"diagnosis_node_report_list.0.diagnosis_dimension_list.0.diagnosis_item_list.0.advice_ids.#"),
					resource.TestCheckResourceAttrSet(rName,
						"diagnosis_node_report_list.0.diagnosis_dimension_list.0.diagnosis_item_list.0.cause_ids.#"),
					resource.TestCheckResourceAttrSet(rName,
						"diagnosis_node_report_list.0.diagnosis_dimension_list.0.diagnosis_item_list.0.impact_ids.#"),
					resource.TestCheckResourceAttr(rName,
						"diagnosis_node_report_list.0.diagnosis_dimension_list.0.diagnosis_item_list.0.error_code", ""),
					resource.TestCheckResourceAttrSet(rName, "diagnosis_node_report_list.0.diagnosis_dimension_list.0.diagnosis_item_list.0.name"),
					resource.TestCheckResourceAttrSet(rName,
						"diagnosis_node_report_list.0.diagnosis_dimension_list.0.diagnosis_item_list.0.result"),
					resource.TestCheckResourceAttr(rName, "diagnosis_node_report_list.0.command_time_taken_list.0.error_code", ""),
					resource.TestCheckResourceAttrSet(rName, "diagnosis_node_report_list.0.command_time_taken_list.0.result"),
					resource.TestCheckResourceAttrSet(rName, "diagnosis_node_report_list.0.command_time_taken_list.0.total_num"),
					resource.TestCheckResourceAttrSet(rName, "diagnosis_node_report_list.0.command_time_taken_list.0.total_usec_sum"),
					resource.TestCheckResourceAttrSet(rName, "diagnosis_node_report_list.0.command_time_taken_list.0.command_list.0.average_usec"),
					resource.TestCheckResourceAttrSet(rName, "diagnosis_node_report_list.0.command_time_taken_list.0.command_list.0.calls_sum"),
					resource.TestCheckResourceAttrSet(rName, "diagnosis_node_report_list.0.command_time_taken_list.0.command_list.0.command_name"),
					resource.TestCheckResourceAttrSet(rName, "diagnosis_node_report_list.0.command_time_taken_list.0.command_list.0.per_usec"),
					resource.TestCheckResourceAttrSet(rName, "diagnosis_node_report_list.0.command_time_taken_list.0.command_list.0.usec_sum"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testDiagnosisReportImportState(rName),
			},
		},
	})
}

func testAccDiagnosisTask_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dcs_diagnosis_task" "test" {
  instance_id = "%s"
  begin_time  = "%s"
  end_time    = "%s"
}`, acceptance.HW_DCS_INSTANCE_ID, acceptance.HW_DCS_BEGIN_TIME, acceptance.HW_DCS_END_TIME)
}

func testDiagnosisReportImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("the resource (%s) not found: %s", name, rs)
		}
		instanceId := rs.Primary.Attributes["instance_id"]
		return fmt.Sprintf("%s/%s", instanceId, rs.Primary.ID), nil
	}
}
