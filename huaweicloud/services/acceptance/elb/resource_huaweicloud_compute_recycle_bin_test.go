package elb

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

func getRecycleBinResourceFunc(cfg *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/elb/recycle-bin"
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating ELB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving ELB recycle bin: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	enable := utils.PathSearch("recycle_bin.enable", getRespBody, false).(bool)
	if !enable {
		return nil, golangsdk.ErrDefault404{}
	}

	return getRespBody, nil
}

func TestAccElbRecycleBin_Basic(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_elb_recycle_bin.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getRecycleBinResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbRecycleBin_basic(10, 20),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "retention_hour", "10"),
					resource.TestCheckResourceAttr(resourceName, "recycle_threshold_hour", "20"),
				),
			},
			{
				Config: testAccElbRecycleBin_basic(15, 30),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "retention_hour", "15"),
					resource.TestCheckResourceAttr(resourceName, "recycle_threshold_hour", "30"),
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

func testAccElbRecycleBin_basic(retentionHour, recycleThresholdDay int) string {
	return fmt.Sprintf(`
resource "huaweicloud_elb_recycle_bin" "test" {
  retention_hour         = %[1]d
  recycle_threshold_hour = %[2]d
}
`, retentionHour, recycleThresholdDay)
}
