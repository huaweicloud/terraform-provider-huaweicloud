package as

import (
	"fmt"
	"strings"
	"testing"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

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
	getASWarmPoolClient, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating AS Client: %s", err)
	}

	getPath := getASWarmPoolClient.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}",
		getASWarmPoolClient.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{scaling_group_id}",
		fmt.Sprintf("%v", state.Primary.Attributes["scaling_group_id"]))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	getASWarmPoolResp, err := getASWarmPoolClient.Request("GET", getPath,
		&getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving AS warm pool: %s", err)
	}

	getASWarmPoolRespBody, err := utils.FlattenResponse(getASWarmPoolResp)
	if err != nil {
		return nil, fmt.Errorf("error flatten response: %s", err)
	}
	status := utils.PathSearch("warm_pool.status", getASWarmPoolRespBody, "").(string)
	if status == "CLOSED" {
		return nil, golangsdk.ErrDefault404{}
	}
	return getASWarmPoolRespBody, nil
}

func TestAccAsWarmPool_basic(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_as_warm_pool.test"

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
				Config: testAsWarmPool_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "min_capacity", "1"),
					resource.TestCheckResourceAttr(rName, "max_capacity", "1"),
					resource.TestCheckResourceAttr(rName, "instance_init_wait_time", "30"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAsWarmPoolImportState(rName),
			},
			{
				Config: testAsWarmPool_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "min_capacity", "2"),
					resource.TestCheckResourceAttr(rName, "max_capacity", "2"),
					resource.TestCheckResourceAttr(rName, "instance_init_wait_time", "60"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
				),
			},
		},
	})
}

func testAsWarmPool_basic() string {
	return fmt.Sprintf(`

resource "huaweicloud_as_warm_pool" "test" {
  scaling_group_id 		  = "%[1]s"
  min_capacity     		  = 1
  max_capacity     		  = 1
  instance_init_wait_time = 30
}
`, acceptance.HW_AS_SCALING_GROUP_ID)
}

func testAsWarmPool_basic_update() string {
	return fmt.Sprintf(`

resource "huaweicloud_as_warm_pool" "test" {
  scaling_group_id 		  = "%[1]s"
  min_capacity     		  = 2
  max_capacity     		  = 2
  instance_init_wait_time = 60
}
`, acceptance.HW_AS_SCALING_GROUP_ID)
}

func testAsWarmPoolImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		var scalingGroupID string
		if scalingGroupID = rs.Primary.Attributes["scaling_group_id"]; scalingGroupID == "" {
			return "", fmt.Errorf("attribute (scaling_group_id) of Resource (%s) not found: %s", name, rs)
		}
		return fmt.Sprintf("%s", scalingGroupID), nil
	}
}
