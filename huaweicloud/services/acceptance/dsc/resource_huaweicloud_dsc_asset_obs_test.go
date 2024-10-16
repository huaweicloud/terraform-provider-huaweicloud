package dsc

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

func getAssetObsResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getAssetObs: Query the asset OBS
	var (
		getAssetObsHttpUrl = "v1/{project_id}/sdg/asset/obs/buckets"
		getAssetObsProduct = "dsc"
	)
	getAssetObsClient, err := cfg.NewServiceClient(getAssetObsProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating AssetObs Client: %s", err)
	}

	getAssetObsPath := getAssetObsClient.Endpoint + getAssetObsHttpUrl
	getAssetObsPath = strings.ReplaceAll(getAssetObsPath, "{project_id}", getAssetObsClient.ProjectID)
	getAssetObsPath += "?added=true"

	getAssetObsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getAssetObsResp, err := getAssetObsClient.Request("GET", getAssetObsPath, &getAssetObsOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving AssetObs: %s", err)
	}

	getAssetObsRespBody, err := utils.FlattenResponse(getAssetObsResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving AssetObs: %s", err)
	}

	assetObs := utils.PathSearch(fmt.Sprintf("buckets[?id=='%s']|[0]", state.Primary.ID), getAssetObsRespBody, nil)
	if assetObs == nil {
		return nil, fmt.Errorf("error retrieving AssetObs: %s", err)
	}
	return assetObs, nil
}

func TestAccAssetObs_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	obsName := acceptance.RandomAccResourceNameWithDash()
	rName := "huaweicloud_dsc_asset_obs.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAssetObsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAssetObs_basic(name, obsName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "bucket_name", "huaweicloud_obs_bucket.test", "bucket"),
					resource.TestCheckResourceAttr(rName, "bucket_policy", "private"),
				),
			},
			{
				Config: testAssetObs_basic(name+"_update", obsName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"_update"),
					resource.TestCheckResourceAttrPair(rName, "bucket_name", "huaweicloud_obs_bucket.test", "bucket"),
					resource.TestCheckResourceAttr(rName, "bucket_policy", "private"),
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

func testAssetObs_basic(name, obsName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%s"
  acl           = "private"
  force_destroy = true

  lifecycle {
    ignore_changes = [
      logging,
    ]
  }
}

resource "huaweicloud_dsc_asset_obs" "test" {
  name          = "%s"
  bucket_name   = huaweicloud_obs_bucket.test.bucket
  bucket_policy = "private"

  depends_on = [huaweicloud_dsc_instance.test]
}
`, testDscInstance_basic(), obsName, name)
}
