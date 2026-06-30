package modelarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/modelarts"
)

func getV2WorkflowScheduleResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region     = acceptance.HW_REGION_NAME
		workflowId = state.Primary.Attributes["workflow_id"]
		scheduleId = state.Primary.ID
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return nil, fmt.Errorf("error creating ModelArts client: %s", err)
	}

	return modelarts.GetV2WorkflowScheduleById(client, workflowId, scheduleId)
}

func TestAccV2WorkflowSchedule_basic(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_modelartsv2_workflow_schedule.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getV2WorkflowScheduleResourceFunc)
	)

	// Each workflow can only manage one schedule resource, using serial testing instead.
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckModelArtsWorkflowId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testAccV2WorkflowSchedule_nonExistentWorkflow(),
				ExpectError: regexp.MustCompile(`error creating ModelArts workflow schedule`),
			},
			{
				Config: testAccV2WorkflowSchedule_basic_step1(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workflow_id", acceptance.HW_MODELARTS_WORKFLOW_ID),
					resource.TestCheckResourceAttr(rName, "content", `{"cron":"0 0 12 * * ?","method":"fixed"}`),
					resource.TestCheckResourceAttr(rName, "enable", "true"),
					resource.TestCheckResourceAttr(rName, "policies.#", "1"),
					resource.TestCheckResourceAttrSet(rName, "policies.0.on_failure"),
					resource.TestCheckResourceAttrSet(rName, "policies.0.on_running"),
					resource.TestCheckResourceAttrSet(rName, "user_id"),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccV2WorkflowSchedule_basic_step2(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workflow_id", acceptance.HW_MODELARTS_WORKFLOW_ID),
					resource.TestCheckResourceAttr(rName, "content", `{"cron":"0 0 13 * * Mon,Tue","method":"fixed"}`),
					resource.TestCheckResourceAttr(rName, "enable", "true"),
					resource.TestCheckResourceAttr(rName, "policies.#", "1"),
					resource.TestCheckResourceAttrSet(rName, "policies.0.on_failure"),
					resource.TestCheckResourceAttrSet(rName, "policies.0.on_running"),
					resource.TestCheckResourceAttrSet(rName, "user_id"),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccV2WorkflowScheduleImportStateIDFunc(rName),
			},
		},
	})
}

func testAccV2WorkflowScheduleImportStateIDFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}

		workflowID := rs.Primary.Attributes["workflow_id"]
		if workflowID == "" {
			return "", fmt.Errorf("attribute (workflow_id) of resource (%s) not found", name)
		}
		return fmt.Sprintf("%s/%s", workflowID, rs.Primary.ID), nil
	}
}

func testAccV2WorkflowSchedule_nonExistentWorkflow() string {
	randomUUID, _ := uuid.NewRandom()

	return fmt.Sprintf(`
resource "huaweicloud_modelartsv2_workflow_schedule" "test" {
  workflow_id = "%[1]s"
  content     = jsonencode({
    cron   = "0 0 12 * * ?"
    method = "fixed"
  })
}
`, randomUUID.String())
}

func testAccV2WorkflowSchedule_basic_step1() string {
	return fmt.Sprintf(`
resource "huaweicloud_modelartsv2_workflow_schedule" "test" {
  workflow_id = "%[1]s"
  content     = jsonencode({
    cron   = "0 0 12 * * ?"
    method = "fixed"
  })
}
`, acceptance.HW_MODELARTS_WORKFLOW_ID)
}

func testAccV2WorkflowSchedule_basic_step2() string {
	return fmt.Sprintf(`
resource "huaweicloud_modelartsv2_workflow_schedule" "test" {
  workflow_id = "%[1]s"
  content     = jsonencode({
    cron   = "0 0 13 * * Mon,Tue"
    method = "fixed"
  })
}
`, acceptance.HW_MODELARTS_WORKFLOW_ID)
}

func TestAccV2WorkflowSchedule_withPolicies(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_modelartsv2_workflow_schedule.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getV2WorkflowScheduleResourceFunc)
	)

	// Each workflow can only manage one schedule resource, using serial testing instead.
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckModelArtsWorkflowId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV2WorkflowSchedule_withPolicies_step1(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workflow_id", acceptance.HW_MODELARTS_WORKFLOW_ID),
					resource.TestCheckResourceAttr(rName, "content", `{"cron":"0 0 12 * * ?","method":"fixed"}`),
					resource.TestCheckResourceAttr(rName, "enable", "true"),
					resource.TestCheckResourceAttr(rName, "policies.#", "1"),
					resource.TestCheckResourceAttrSet(rName, "policies.0.on_failure"),
					resource.TestCheckResourceAttrSet(rName, "policies.0.on_running"),
					resource.TestCheckResourceAttrSet(rName, "user_id"),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccV2WorkflowSchedule_withPolicies_step2(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workflow_id", acceptance.HW_MODELARTS_WORKFLOW_ID),
					resource.TestCheckResourceAttr(rName, "content", `{"cron":"0 0 13 * * Mon,Tue","method":"fixed"}`),
					resource.TestCheckResourceAttr(rName, "enable", "true"),
					resource.TestCheckResourceAttr(rName, "policies.#", "1"),
					resource.TestCheckResourceAttrSet(rName, "policies.0.on_failure"),
					resource.TestCheckResourceAttrSet(rName, "policies.0.on_running"),
					resource.TestCheckResourceAttrSet(rName, "user_id"),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccV2WorkflowScheduleImportStateIDFunc(rName),
			},
		},
	})
}

func testAccV2WorkflowSchedule_withPolicies_step1() string {
	return fmt.Sprintf(`
resource "huaweicloud_modelartsv2_workflow_schedule" "test" {
  workflow_id = "%[1]s"
  content     = jsonencode({
    cron   = "0 0 12 * * ?"
    method = "fixed"
  })
  policies {
    on_failure = "retry"
    on_running = "cancel"
  }
}
`, acceptance.HW_MODELARTS_WORKFLOW_ID)
}

func testAccV2WorkflowSchedule_withPolicies_step2() string {
	return fmt.Sprintf(`
resource "huaweicloud_modelartsv2_workflow_schedule" "test" {
  workflow_id = "%[1]s"
  content     = jsonencode({
    cron   = "0 0 13 * * Mon,Tue"
    method = "fixed"
  })
  policies {
    on_failure = "retry"
    on_running = "cancel"
  }
}
`, acceptance.HW_MODELARTS_WORKFLOW_ID)
}
