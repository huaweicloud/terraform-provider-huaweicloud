package drs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceDrsSmnBatchSet_basic(t *testing.T) {
	resourceName := "huaweicloud_drs_smn_batch_set.test"

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
				Config: testAccResourceDrsSmnBatchSet_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "results.#"),
					resource.TestMatchResourceAttr(resourceName, "results.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.status"),
				),
			},
		},
	})
}

func testAccResourceDrsSmnBatchSet_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_smn_batch_set" "test" {
  jobs {
    job_id      = "%s"
    status      = "START_JOB_FAILED"
    engine_type = "mysql"
  }

  alarm_notify_info {
    topic_urn     = "%s"
    delay_time    = 1200
    rto_delay     = 20
    rpo_delay     = 20
    alarm_to_user = false
  }
}
`, acceptance.HW_DRS_JOB_ID, acceptance.HW_SMN_TOPIC_URN)
}
