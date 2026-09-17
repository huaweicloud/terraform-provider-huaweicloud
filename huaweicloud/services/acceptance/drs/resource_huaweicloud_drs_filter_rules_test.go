package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccFilterRules_basic(t *testing.T) {
	var (
		rName = "huaweicloud_drs_filter_rules.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testFilterRules_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "job_id", acceptance.HW_DRS_JOB_ID),
					resource.TestCheckResourceAttr(rName, "data_filter_type", "common"),
					resource.TestCheckResourceAttr(rName, "source", "job"),
					resource.TestCheckResourceAttr(rName, "general_filtering_list.#", "1"),
					resource.TestCheckResourceAttr(rName, "general_filtering_list.0.source_filter_condition", "id > 1"),
				),
			},
			{
				Config: testFilterRules_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "job_id", acceptance.HW_DRS_JOB_ID),
					resource.TestCheckResourceAttr(rName, "general_filtering_list.0.source_filter_condition", "id > 10"),
				),
			},
		},
	})
}

func testFilterRules_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_filter_rules" "test" {
  job_id           = "%[1]s"
  data_filter_type = "common"
  source           = "job"

  general_filtering_list {
    filter_object_list {
      id          = "test_drs-*-*-mysql1"
      parent_id   = "test_drs"
      object_name = "mysql1"
    }

    source_filter_condition = "id > 1"
  }
}
`, acceptance.HW_DRS_JOB_ID)
}

func testFilterRules_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_filter_rules" "test" {
  job_id           = "%[1]s"
  data_filter_type = "common"
  source           = "job"

  general_filtering_list {
    filter_object_list {
      id          = "test_drs-*-*-mysql1"
      parent_id   = "test_drs"
      object_name = "mysql1"
    }

    source_filter_condition = "id > 10"
  }
}
`, acceptance.HW_DRS_JOB_ID)
}
