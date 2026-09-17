package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccCheckFilterRules_basic(t *testing.T) {
	var (
		rName = "huaweicloud_drs_check_filter_rules.test"
	)

	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testCheckFilterRules_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "job_id", acceptance.HW_DRS_JOB_ID),
					resource.TestCheckResourceAttr(rName, "data_filter_type", "common"),
					resource.TestCheckResourceAttr(rName, "source", "job"),
					resource.TestCheckResourceAttrSet(rName, "query_id"),
				),
			},
		},
	})
}

func testCheckFilterRules_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_check_filter_rules" "test" {
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
