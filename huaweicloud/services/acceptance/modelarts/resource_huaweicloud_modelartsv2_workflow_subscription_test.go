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

func getV2WorkflowSubscriptionResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("modelarts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ModelArts client: %s", err)
	}

	return modelarts.GetV2WorkflowSubscriptionById(client, state.Primary.Attributes["workflow_id"], state.Primary.ID)
}

func TestAccV2WorkflowSubscription_basic(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_modelartsv2_workflow_subscription.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getV2WorkflowSubscriptionResourceFunc)

		name = acceptance.RandomAccResourceName()
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckModelArtsWorkflowId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testAccV2WorkflowSubscription_nonExistentWorkflow(name),
				ExpectError: regexp.MustCompile(`error creating ModelArts workflow subscription`),
			},
			{
				Config: testAccV2WorkflowSubscription_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workflow_id", acceptance.HW_MODELARTS_WORKFLOW_ID),
					resource.TestCheckResourceAttr(rName, "topic_urns.#", "2"),
					resource.TestCheckResourceAttrPair(rName, "topic_urns.0", "huaweicloud_smn_topic.test.0", "topic_urn"),
					resource.TestCheckResourceAttrPair(rName, "topic_urns.1", "huaweicloud_smn_topic.test.1", "topic_urn"),
					resource.TestCheckResourceAttr(rName, "events.#", "3"),
					resource.TestCheckResourceAttr(rName, "events.0", "service_step:wait_inputs,hold,completed"),
					resource.TestCheckResourceAttr(rName, "events.1", "labeling:wait_inputs,hold,failed,create_failed"),
					resource.TestCheckResourceAttr(rName, "events.2", "*:wait_inputs,hold,completed"),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccV2WorkflowSubscription_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workflow_id", acceptance.HW_MODELARTS_WORKFLOW_ID),
					resource.TestCheckResourceAttr(rName, "topic_urns.#", "2"),
					resource.TestCheckResourceAttrPair(rName, "topic_urns.0", "huaweicloud_smn_topic.test.1", "topic_urn"),
					resource.TestCheckResourceAttrPair(rName, "topic_urns.1", "huaweicloud_smn_topic.test.2", "topic_urn"),
					resource.TestCheckResourceAttr(rName, "events.#", "3"),
					resource.TestCheckResourceAttr(rName, "events.0", "release:wait_inputs,hold,failed,create_failed"),
					resource.TestCheckResourceAttr(rName, "events.1", "model_step:completed,failed,create_failed"),
					resource.TestCheckResourceAttr(rName, "events.2", "training_job:wait_inputs,hold,completed"),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccV2WorkflowSubscription_basic_step3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workflow_id", acceptance.HW_MODELARTS_WORKFLOW_ID),
					resource.TestCheckResourceAttr(rName, "topic_urns.#", "2"),
					resource.TestCheckResourceAttrPair(rName, "topic_urns.0", "huaweicloud_smn_topic.test.1", "topic_urn"),
					resource.TestCheckResourceAttrPair(rName, "topic_urns.1", "huaweicloud_smn_topic.test.0", "topic_urn"),
					resource.TestCheckResourceAttr(rName, "events.#", "3"),
					resource.TestCheckResourceAttr(rName, "events.0", "*:wait_inputs,hold,completed"),
					resource.TestCheckResourceAttr(rName, "events.1", "labeling:wait_inputs,hold,failed,create_failed"),
					resource.TestCheckResourceAttr(rName, "events.2", "service_step:wait_inputs,hold,completed"),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccV2WorkflowSubscriptionImportStateIDFunc(rName),
			},
		},
	})
}

func testAccV2WorkflowSubscriptionImportStateIDFunc(name string) resource.ImportStateIdFunc {
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

func testAccV2WorkflowSubscription_smnTopics(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "test" {
  count = 3

  name = format("%%s-%%d", "%[1]s", count.index)
}
`, name)
}

func testAccV2WorkflowSubscription_nonExistentWorkflow(name string) string {
	randomUUID, _ := uuid.NewRandom()

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelartsv2_workflow_subscription" "test" {
  workflow_id = "%[2]s"
  topic_urns  = slice(huaweicloud_smn_topic.test[*].topic_urn, 0, 2)

  events = [
    "service_step:wait_inputs,hold,completed",
    "labeling:wait_inputs,hold,failed,create_failed",
    "*:wait_inputs,hold,completed",
  ]
}
`, testAccV2WorkflowSubscription_smnTopics(name), randomUUID.String())
}

func testAccV2WorkflowSubscription_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelartsv2_workflow_subscription" "test" {
  workflow_id = "%[2]s"
  topic_urns  = slice(huaweicloud_smn_topic.test[*].topic_urn, 0, 2)

  events = [
    "service_step:wait_inputs,hold,completed",
    "labeling:wait_inputs,hold,failed,create_failed",
    "*:wait_inputs,hold,completed",
  ]
}
`, testAccV2WorkflowSubscription_smnTopics(name), acceptance.HW_MODELARTS_WORKFLOW_ID)
}

func testAccV2WorkflowSubscription_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelartsv2_workflow_subscription" "test" {
  workflow_id = "%[2]s"
  topic_urns  = slice(huaweicloud_smn_topic.test[*].topic_urn, 1, 3)

  events = [
    "release:wait_inputs,hold,failed,create_failed",
    "model_step:completed,failed,create_failed",
    "training_job:wait_inputs,hold,completed",
  ]
}
`, testAccV2WorkflowSubscription_smnTopics(name), acceptance.HW_MODELARTS_WORKFLOW_ID)
}

// The configuration is the same as the basic step 1, but the order of topic urns and events are different.
func testAccV2WorkflowSubscription_basic_step3(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelartsv2_workflow_subscription" "test" {
  workflow_id = "%[2]s"
  topic_urns  = reverse(slice(huaweicloud_smn_topic.test[*].topic_urn, 0, 2))

  events = [
    "*:wait_inputs,hold,completed",
    "labeling:wait_inputs,hold,failed,create_failed",
    "service_step:wait_inputs,hold,completed",
  ]
}
`, testAccV2WorkflowSubscription_smnTopics(name), acceptance.HW_MODELARTS_WORKFLOW_ID)
}
