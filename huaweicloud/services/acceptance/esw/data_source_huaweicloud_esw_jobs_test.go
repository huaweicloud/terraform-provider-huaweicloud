package esw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccEswJobs_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_esw_jobs.test"
	dc := acceptance.InitDataSourceCheck(rName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceEswFJobs_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "jobs.#"),
					resource.TestCheckResourceAttrSet(rName, "jobs.0.id"),
					resource.TestCheckResourceAttrSet(rName, "jobs.0.name"),
					resource.TestCheckResourceAttrSet(rName, "jobs.0.status"),
					resource.TestCheckResourceAttrSet(rName, "jobs.0.begin_time"),
					resource.TestCheckResourceAttrSet(rName, "jobs.0.end_time"),
					resource.TestCheckResourceAttrSet(rName, "jobs.0.resource_id"),
					resource.TestCheckResourceAttrSet(rName, "jobs.0.resource_name"),
					resource.TestCheckResourceAttrSet(rName, "jobs.0.resource_type"),
					resource.TestCheckResourceAttrSet(rName, "jobs.0.project_id"),
				),
			},
		},
	})
}

func testAccDatasourceEswFJobs_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_esw_jobs" "test" {
  resource_id = huaweicloud_esw_instance.test.id
}
`, testAccEswInstance_basic(name))
}
