package cpts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cpts"
)

func getTaskResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("cpts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CPTS client: %s", err)
	}

	return cpts.GetTaskDetail(client, state.Primary.ID)
}

func TestAccTask_basic(t *testing.T) {
	var obj interface{}

	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_cpts_task.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getTaskResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testTask_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "2"),
					resource.TestCheckResourceAttr(resourceName, "benchmark_concurrency", "200"),
					resource.TestCheckResourceAttrPair(resourceName, "project_id",
						"huaweicloud_cpts_project.test", "id"),
				),
			},
			{
				Config: testTask_basic_update1(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-update", rName)),
					resource.TestCheckResourceAttr(resourceName, "status", "2"),
					resource.TestCheckResourceAttr(resourceName, "benchmark_concurrency", "102"),
					resource.TestCheckResourceAttrPair(resourceName, "project_id",
						"huaweicloud_cpts_project.test", "id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster_type", "operation"},
			},
		},
	})
}

func testTask_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cpts_task" "test" {
  name                  = "%[2]s"
  project_id            = huaweicloud_cpts_project.test.id
  benchmark_concurrency = 200
}
`, testProject_basic(rName), rName)
}

func testTask_basic_update1(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cpts_task" "test" {
  name                  = "%[2]s-update"
  project_id            = huaweicloud_cpts_project.test.id
  benchmark_concurrency = 102
}
`, testProject_basic(rName), rName)
}
