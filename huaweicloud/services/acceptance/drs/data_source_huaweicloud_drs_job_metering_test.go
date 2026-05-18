package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDrsJobMetering_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_drs_job_metering.test"
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
				Config: testAccDataSourceDrsJobMetering_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "product_info_list.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "product_info_list.0.cloud_service_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "product_info_list.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "product_info_list.0.resource_spec_code"),
				),
			},
		},
	})
}

func testAccDataSourceDrsJobMetering_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_job_metering" "test" {
  job_id = "%s"
}
`, acceptance.HW_DRS_JOB_ID)
}
