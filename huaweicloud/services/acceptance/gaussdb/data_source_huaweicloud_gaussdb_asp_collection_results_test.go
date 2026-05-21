package gaussdb

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussDbAspCollectionResults_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_asp_collection_results.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbAspCollectionResults_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "asp.#"),
					resource.TestCheckResourceAttrSet(dataSource, "asp.0.job_id"),
					resource.TestCheckResourceAttrSet(dataSource, "asp.0.file_size"),
					resource.TestCheckResourceAttrSet(dataSource, "asp.0.file_path"),
					resource.TestCheckResourceAttrSet(dataSource, "asp.0.file_name"),
					resource.TestCheckResourceAttrSet(dataSource, "asp.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "asp.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "asp.0.download_url"),
					resource.TestCheckResourceAttrSet(dataSource, "asp.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "asp.0.obs_bucket.#"),
					resource.TestCheckResourceAttrSet(dataSource, "asp.0.obs_bucket.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "asp.0.obs_bucket.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "asp.0.obs_bucket.0.url"),
					resource.TestCheckResourceAttrSet(dataSource, "asp.0.obs_bucket.0.port"),
					resource.TestCheckResourceAttrSet(dataSource, "asp.0.obs_bucket.0.domain_id"),
					resource.TestCheckOutput("time_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceGaussdbAspCollectionResults_basic() string {
	startTime := time.Now().UTC().Add(-8 * time.Hour).Format("2006-01-02T15:04:05+0000")
	endTime := time.Now().UTC().Add(8 * time.Hour).Format("2006-01-02T15:04:05+0000")
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_asp_collection_results" "test" {
  instance_id = "%[1]s"
}

data "huaweicloud_gaussdb_asp_collection_results" "time_filter" {
  instance_id = "%[1]s"
  start_time  = "%[2]s"
  end_time    = "%[3]s"
}
output "time_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_asp_collection_results.time_filter.asp) > 0
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID, startTime, endTime)
}
