package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocIncidentTasks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_incident_tasks.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: "0.12.1",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocIncidentTasks_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.key"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.state"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.operations.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.operations.0.task_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.operations.0.key"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocIncidentTasks_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "time_sleep" "wait_10_seconds" {
  depends_on = [huaweicloud_coc_incident.test]

  create_duration = "10s"
}

data "huaweicloud_coc_incident_tasks" "test" {
  depends_on  = [time_sleep.wait_10_seconds]
  incident_id = huaweicloud_coc_incident.test.id
}
`, testIncident_basic(name))
}
