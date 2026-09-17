package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDrsAssociateSmnTopic_basic(t *testing.T) {
	resourceName := "huaweicloud_drs_associate_smn_topic.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobId(t)
			acceptance.TestAccPreCheckSmnTopicUrn(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDrsAssociateSmnTopic_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "job_id", acceptance.HW_DRS_JOB_ID),
					resource.TestCheckResourceAttr(resourceName, "topic_urn", acceptance.HW_SMN_TOPIC_URN),
					resource.TestCheckResourceAttr(resourceName, "is_alarm_to_user", "true"),
				),
			},
		},
	})
}

func testAccDrsAssociateSmnTopic_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_associate_smn_topic" "test" {
  job_id                     = "%[1]s"
  topic_urn                  = "%[2]s"
  is_alarm_to_user           = true
  increment_delay_threshold  = 500
}
`, acceptance.HW_DRS_JOB_ID, acceptance.HW_SMN_TOPIC_URN)
}
