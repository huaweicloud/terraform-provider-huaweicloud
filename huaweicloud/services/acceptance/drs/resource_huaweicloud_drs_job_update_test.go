package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDrsJobUpdate_policy(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDrsJobUpdate_policy(),
			},
		},
	})
}

func testAccDrsJobUpdate_policy() string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_job_update" "test" {
  job_id = "%[1]s"
  type   = "policy"

  policy_config {
    sync_type_policy = "tableData"
  }
}
`, acceptance.HW_DRS_JOB_ID)
}

func TestAccDrsJobUpdate_name(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDrsJobUpdate_name(),
			},
		},
	})
}

func testAccDrsJobUpdate_name() string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_job_update" "test" {
  job_id = "%[1]s"
  type   = "name"

  base_info {
    name        = "drs-test-rename"
    description = "updated by terraform"
  }
}
`, acceptance.HW_DRS_JOB_ID)
}

func TestAccDrsJobUpdate_speed_limit(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDrsJobUpdate_speed_limit(),
			},
		},
	})
}

func testAccDrsJobUpdate_speed_limit() string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_job_update" "test" {
  job_id = "%[1]s"
  type   = "speed_limit"

  speed_limit {
    speed      = "10"
    start_time = "08:00"
    end_time   = "22:00"
    is_effect  = true
  }
}
`, acceptance.HW_DRS_JOB_ID)
}
