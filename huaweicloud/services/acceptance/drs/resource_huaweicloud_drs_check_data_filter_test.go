package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceCheckDataFilter_basic(t *testing.T) {
	resourceName := "huaweicloud_drs_check_data_filter.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCheckDataFilter_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
		},
	})
}

func testAccResourceCheckDataFilter_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_check_data_filter" "test" {
  job_id = "%s"

  data_process_info {
    filter_conditions {
      filtering_type = "contentConditionalFilter"
      value          = "id>1"
    }

    db_object {
      object_scope = "table"

      object_info = <<EOT
      {
        "dyh4" : {
          "name" : "dyh4",
          "all" : false,
          "tables" : {
            "test1_table1" : {
              "name" : "test1_table1",
              "type" : "table",
              "all" : true
            },
            "test1_table10" : {
              "name" : "test1_table10",
              "type" : "table",
              "all" : true
            },
            "test1_table11" : {
              "name" : "test1_table11",
              "type" : "table",
              "all" : true
            }
          }
        }
      }
      EOT
    }
  }
}
`, acceptance.HW_DRS_JOB_ID)
}
