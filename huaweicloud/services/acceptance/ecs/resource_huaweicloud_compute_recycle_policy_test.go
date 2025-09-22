package ecs

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

func getRecyclePolicyResourceFunc(cfg *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v1/{project_id}/recycle-bin"
		product = "ecs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating ECS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving ECS recycle: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	switchValue := utils.PathSearch("switch", getRespBody, "").(string)
	if switchValue == "off" {
		return nil, golangsdk.ErrDefault404{}
	}

	return getRespBody, nil
}

func TestAccComputeRecyclePolicy_Basic(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_compute_recycle_policy.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getRecyclePolicyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRecyclePolicy_basic(10, 20),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "retention_hour", "10"),
					resource.TestCheckResourceAttr(resourceName, "recycle_threshold_day", "20"),
				),
			},
			{
				Config: testAccComputeRecyclePolicy_basic(15, 30),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "retention_hour", "15"),
					resource.TestCheckResourceAttr(resourceName, "recycle_threshold_day", "30"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeRecyclePolicy_basic(retentionHour, recycleThresholdDay int) string {
	return fmt.Sprintf(`
resource "huaweicloud_compute_recycle_policy" "test" {
  retention_hour        = %[1]d
  recycle_threshold_day = %[2]d
}
`, retentionHour, recycleThresholdDay)
}
