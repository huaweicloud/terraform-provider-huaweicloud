package codeartsdeploy

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCodeartsDeployApplicationDeploymentRecords_basic(t *testing.T) {
	dataSource := "data.huaweicloud_codearts_deploy_application_deployment_records.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCodeartsDeployApplicationDeploymentRecords_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "records.#", "1"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.duration"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.state"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.operator"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.release_id"),
				),
			},
		},
	})
}

func testDataSourceCodeartsDeployApplicationDeploymentRecords_basic(name string) string {
	date := strings.Split(time.Now().Format("2006-01-02T15:04:05Z"), "T")[0]

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_codearts_deploy_application_deployment_records" "test" {
  depends_on = [huaweicloud_codearts_deploy_application_deploy.test]

  project_id = huaweicloud_codearts_project.test.id
  task_id    = huaweicloud_codearts_deploy_application.test.task_id
  start_date = "%[2]s"
  end_date   = "%[2]s"
}
`, testAccDeployApplicationDeploy_basic(name), date)
}
