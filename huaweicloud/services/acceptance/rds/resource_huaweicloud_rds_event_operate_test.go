package rds

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccRdsEventOperate_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_rds_event_operate.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsEventOperate_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "results.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.error_code"),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.error_msg"),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.success"),
					resource.TestCheckResourceAttrSet(resourceName, "results.1.id"),
					resource.TestCheckResourceAttrSet(resourceName, "results.1.error_code"),
					resource.TestCheckResourceAttrSet(resourceName, "results.1.error_msg"),
					resource.TestCheckResourceAttrSet(resourceName, "results.1.success"),
				),
			},
		},
	})
}

func testAccRdsEventOperate_basic(name string) string {
	randEventId1, _ := uuid.NewRandom()
	randEventId2, _ := uuid.NewRandom()
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_event_operate" "test" {
  event_instances {
    event_id    = "%[2]s"
    instance_id = huaweicloud_rds_instance.test.id
  }
  event_instances {
    event_id    = "%[3]s"
    instance_id = huaweicloud_rds_instance.test.id
  }

  operation_type = "reservation"

  event_schedule_window {
    planned_day = "2026-05-05"
    start_time  = "06:00"
    end_time    = "12:00"
  }
}
`, testAccRdsInstance_mysql_step1(name), randEventId1.String(), randEventId2.String())
}
