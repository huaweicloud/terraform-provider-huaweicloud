package swr

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

func getSwrImageAutoSyncResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getSwrImageAutoSync: Query SWR image auto sync
	var (
		getSwrImageAutoSyncHttpUrl = "v2/manage/namespaces/{namespace}/repos/{repository}/sync_repo"
		getSwrImageAutoSyncProduct = "swr"
	)
	getSwrImageAutoSyncClient, err := cfg.NewServiceClient(getSwrImageAutoSyncProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating SWR client: %s", err)
	}

	organization := state.Primary.Attributes["organization"]
	repository := strings.ReplaceAll(state.Primary.Attributes["repository"], "/", "$")
	targetRegion := state.Primary.Attributes["target_region"]
	targetOrganization := state.Primary.Attributes["target_organization"]

	getSwrImageAutoSyncPath := getSwrImageAutoSyncClient.Endpoint + getSwrImageAutoSyncHttpUrl
	getSwrImageAutoSyncPath = strings.ReplaceAll(getSwrImageAutoSyncPath, "{namespace}", organization)
	getSwrImageAutoSyncPath = strings.ReplaceAll(getSwrImageAutoSyncPath, "{repository}", repository)

	getSwrImageSyncOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getSwrImageAutoSyncResp, err := getSwrImageAutoSyncClient.Request("GET",
		getSwrImageAutoSyncPath, &getSwrImageSyncOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving SWR image sync: %s", err)
	}

	getSwrImageAutoSyncRespBody, err := utils.FlattenResponse(getSwrImageAutoSyncResp)
	if err != nil {
		return nil, err
	}

	for _, res := range getSwrImageAutoSyncRespBody.([]interface{}) {
		resTargetRegion := utils.PathSearch("remoteRegionId", res, "").(string)
		resTargetOrganization := utils.PathSearch("remoteNamespace", res, "").(string)
		if resTargetRegion == targetRegion && resTargetOrganization == targetOrganization {
			return res, nil
		}
	}

	return nil, fmt.Errorf("error get SWR image auto sync")
}

func TestAccSwrImageAutoSync_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_swr_image_auto_sync.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getSwrImageAutoSyncResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSwrTargetRegion(t)
			acceptance.TestAccPreCheckSwrTargetOrigination(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testSwrImageAutoSync_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "organization",
						"huaweicloud_swr_organization.test", "name"),
					resource.TestCheckResourceAttrPair(rName, "repository",
						"huaweicloud_swr_repository.test", "name"),
					resource.TestCheckResourceAttr(rName, "target_region",
						acceptance.HW_SWR_TARGET_REGION),
					resource.TestCheckResourceAttr(rName, "target_organization",
						acceptance.HW_SWR_TARGET_ORGANIZATION),
					resource.TestCheckResourceAttr(rName, "override", "true"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
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

func testSwrImageAutoSync_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_swr_image_auto_sync" "test" {
  organization        = huaweicloud_swr_organization.test.name
  repository          = huaweicloud_swr_repository.test.name
  target_region       = "%[2]s"
  target_organization = "%[3]s"
  override            = true
}
`, testAccSWRRepository_basic(name), acceptance.HW_SWR_TARGET_REGION, acceptance.HW_SWR_TARGET_ORGANIZATION)
}
