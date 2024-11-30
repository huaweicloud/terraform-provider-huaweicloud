package cdm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceCdmJobExecutionRecords_basic(t *testing.T) {
	rName := "data.huaweicloud_cdm_job_execution_records.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()
	bucketName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceCdmJobExecutionRecords_basic(name, bucketName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "records.0.is_incrementing"),
					resource.TestCheckResourceAttrSet(rName, "records.0.is_execute_auto"),
					resource.TestCheckResourceAttrSet(rName, "records.0.is_delete_job"),
					resource.TestCheckResourceAttrSet(rName, "records.0.creation_user"),
					resource.TestCheckResourceAttrSet(rName, "records.0.creation_date"),
					resource.TestCheckResourceAttrSet(rName, "records.0.external_id"),
					resource.TestCheckResourceAttrSet(rName, "records.0.progress"),
					resource.TestCheckResourceAttrSet(rName, "records.0.submission_id"),
					resource.TestCheckResourceAttrSet(rName, "records.0.execute_date"),
					resource.TestCheckResourceAttrSet(rName, "records.0.status"),
					resource.TestCheckResourceAttrSet(rName, "records.0.counters.#"),
				),
			},
		},
	})
}

func testAccDatasourceCdmJobExecutionRecords_basic(name, bucketName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cdm_job_execution_records" "test" {
  depends_on = [huaweicloud_cdm_job.test]

  cluster_id = huaweicloud_cdm_cluster.test.id
  job_name   = huaweicloud_cdm_job.test.name
}`, testAccCdmJob_basic(name, bucketName))
}
