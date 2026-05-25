package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDrsDbObject_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_drs_db_object.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDrsDbObject_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "max_table_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "object_scope"),
					resource.TestCheckResourceAttrSet(dataSourceName, "object_info.#"),
				),
			},
		},
	})
}

func testDataSourceDrsDbObject_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_db_object" "test" {
  job_id = "%s"
  type   = "modified"
}
`, acceptance.HW_DRS_JOB_ID)
}
