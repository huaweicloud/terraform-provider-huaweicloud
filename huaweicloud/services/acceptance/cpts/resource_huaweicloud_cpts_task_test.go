package cpts

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cpts/v1/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getTaskResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.HcCptsV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CPTS v1 client: %s", err)
	}

	id, err := strconv.ParseInt(state.Primary.ID, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("the Task ID must be integer: %s", err)
	}

	request := &model.ShowTaskRequest{
		TaskId: int32(id),
	}

	return client.ShowTask(request)
}

func TestAccTask_basic(t *testing.T) {
	var obj model.CreateTaskResponse

	rName := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
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
				Config: testTask_basic(rName, 200),
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
				Config: testTask_basic(updateName, 102),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
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
				ImportStateVerifyIgnore: []string{"cluster_type"},
			},
		},
	})
}

func testTask_basic(rName string, concurrency int) string {
	projectTaskConfig := testProject_basic(rName, "created by acc test")
	return fmt.Sprintf(`
%s

resource "huaweicloud_cpts_task" "test" {
  name                  = "%s"
  project_id            = huaweicloud_cpts_project.test.id
  benchmark_concurrency = %d
}
`, projectTaskConfig, rName, concurrency)
}
