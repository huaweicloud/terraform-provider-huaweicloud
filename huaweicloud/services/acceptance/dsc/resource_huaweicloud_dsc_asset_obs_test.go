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
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v1/{project_id}/sdg/asset/obs/buckets"
		product = "dsc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += "?added=true"
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DSC asset OBS buckets: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	expression := fmt.Sprintf("buckets[?id=='%s']|[0]", state.Primary.ID)
	assetObs := utils.PathSearch(expression, respBody, nil)
	if assetObs == nil {
		return nil, golangsdk.ErrDefault404{}
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
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Configure a DSC instance with OBS authorization enabled.
			acceptance.TestAccPrecheckDscInstance(t)
		},
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
resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%s"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_dsc_asset_obs" "test" {
  name          = "%s"
  bucket_name   = huaweicloud_obs_bucket.test.bucket
  bucket_policy = "private"
}
`, obsName, name)
}
