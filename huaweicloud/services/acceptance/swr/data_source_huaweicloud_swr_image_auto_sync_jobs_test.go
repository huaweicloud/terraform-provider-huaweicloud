package swr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrImageAutoSyncJobs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_image_auto_sync_jobs.test"
	dc := acceptance.InitDataSourceCheck(dataSource)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSwrTargetRegion(t)
			acceptance.TestAccPreCheckSwrTargetOrigination(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrImageAutoSyncJobs_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.organization"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.override"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.remote_organization"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.remote_region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.repo_name"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.sync_operator_id"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.sync_operator_name"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.tag"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.updated_at"),
				),
			},
		},
	})
}

func testDataSourceSwrImageAutoSyncJobs_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_swr_image_auto_sync_jobs" "test" {
  depends_on = [huaweicloud_swr_image_auto_sync.test]
  
  organization = huaweicloud_swr_organization.test.name
  repository   = huaweicloud_swr_repository.test.name
}
`, testSwrImageAutoSync_basic(name))
}
