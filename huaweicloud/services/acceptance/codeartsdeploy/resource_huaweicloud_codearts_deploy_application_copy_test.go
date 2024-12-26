package codeartsdeploy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccCodeArtsDeployApplicationCopy_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testCodeArtsDeployApplicationCopy_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("applications_len_before_copy", "true"),
					resource.TestCheckOutput("applications_len_after_copy", "true"),
				),
			},
		},
	})
}

func testCodeArtsDeployApplicationCopy_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_codearts_deploy_applications" "beforeCopy" {
  depends_on = [huaweicloud_codearts_deploy_application.test]

  project_id = huaweicloud_codearts_project.test.id
}

resource "huaweicloud_codearts_deploy_application_copy" "test" {
  app_id = huaweicloud_codearts_deploy_application.test.id
}

data "huaweicloud_codearts_deploy_applications" "afterCopy" {
  depends_on = [huaweicloud_codearts_deploy_application_copy.test]

  project_id = huaweicloud_codearts_project.test.id
}

output "applications_len_before_copy" {
  value = length(data.huaweicloud_codearts_deploy_applications.beforeCopy.applications) == 1
}

output "applications_len_after_copy" {
  value = length(data.huaweicloud_codearts_deploy_applications.afterCopy.applications) == 2
}
`, testDeployApplication_basic(name))
}
