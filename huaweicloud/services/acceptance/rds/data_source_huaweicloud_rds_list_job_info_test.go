package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsListJobInfo_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_list_job_info.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsListJobInfo_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "job.#"),
					resource.TestCheckResourceAttrSet(dataSource, "job.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "job.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "job.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "job.0.created"),
					resource.TestCheckResourceAttrSet(dataSource, "job.0.ended"),
					resource.TestCheckResourceAttrSet(dataSource, "job.0.entities"),
				),
			},
		},
	})
}

func testDataSourceRdsListJobInfo_basic(_ string) string {
	return fmt.Sprintf(`
  data "huaweicloud_rds_list_job_info" "test" {
    job_id = "%s"
  }
`, acceptance.HW_RDS_INSTANT_JOB_ID)
}
