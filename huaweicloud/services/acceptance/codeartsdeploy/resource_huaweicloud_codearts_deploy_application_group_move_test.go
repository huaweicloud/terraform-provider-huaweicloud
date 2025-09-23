package codeartsdeploy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDeployApplicationGroupMove_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_deploy_application_group.test2"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDeployApplicationGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDeployApplicationGroupMove_base(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "ordinal", "2"),
				),
			},
			{
				Config: testDeployApplicationGroupMove_basic(name),
			},
			{
				Config: testDeployApplicationGroupMove_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "ordinal", "1"),
				),
			},
		},
	})
}

func testDeployApplicationGroupMove_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_deploy_application_group" "test1" {
  project_id = huaweicloud_codearts_project.test.id
  name       = "%[2]s-1"
}

resource "huaweicloud_codearts_deploy_application_group" "test2" {
  depends_on = [huaweicloud_codearts_deploy_application_group.test1]

  project_id = huaweicloud_codearts_project.test.id
  name       = "%[2]s-2"
}
`, testProject_basic(name), name)
}

func testDeployApplicationGroupMove_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_deploy_application_group_move" "test" {
  project_id = huaweicloud_codearts_project.test.id
  group_id   = huaweicloud_codearts_deploy_application_group.test2.id
  movement   = -1
}
`, testDeployApplicationGroupMove_base(name))
}
