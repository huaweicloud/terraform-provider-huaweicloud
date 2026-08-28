package as

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getAsWarmPoolResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/scaling-groups/{scaling_group_id}/warm-pool"
		product = "autoscaling"
		region  = acceptance.HW_REGION_NAME
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating AS Client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{scaling_group_id}", fmt.Sprintf("%v", state.Primary.Attributes["scaling_group_id"]))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving AS warm pool: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, fmt.Errorf("error flatten response: %s", err)
	}
	status := utils.PathSearch("warm_pool.status", getRespBody, "").(string)
	if status == "CLOSED" {
		return nil, golangsdk.ErrDefault404{}
	}
	return getRespBody, nil
}

func TestAccAsWarmPool_basic(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_as_warm_pool.test"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAsWarmPoolResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAsWarmPool_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "min_capacity", "1"),
					resource.TestCheckResourceAttr(rName, "max_capacity", "1"),
					resource.TestCheckResourceAttr(rName, "instance_init_wait_time", "30"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
				),
			},
			{
				Config: testAsWarmPool_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "min_capacity", "2"),
					resource.TestCheckResourceAttr(rName, "max_capacity", "2"),
					resource.TestCheckResourceAttr(rName, "instance_init_wait_time", "60"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAsWarmPool_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_as_warm_pool" "test" {
  scaling_group_id 		  = huaweicloud_as_group.acc_as_group.id
  min_capacity     		  = 1
  max_capacity     		  = 1
  instance_init_wait_time = 30
}
`, testASGroup_basic(name))
}

func testAsWarmPool_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_as_warm_pool" "test" {
  scaling_group_id 		  = huaweicloud_as_group.acc_as_group.id
  min_capacity     		  = 2
  max_capacity     		  = 2
  instance_init_wait_time = 60
}
`, testASGroup_basic(name))
}
