package vod

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	vod "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vod/v1/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getResourceAsset(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.HcVodV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating VOD client: %s", err)
	}

	return client.ShowAssetDetail(&vod.ShowAssetDetailRequest{AssetId: state.Primary.ID})
}

func TestAccMediaAsset_basic(t *testing.T) {
	var asset vod.ShowAssetDetailResponse
	rName := acceptance.RandomAccResourceNameWithDash()
	updateName := rName + "-update"
	description := "test video"
	descriptionUpdate := "test video update"
	resourceName := "huaweicloud_vod_media_asset.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&asset,
		getResourceAsset,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckVODMediaAsset(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccMediaAsset_basic(testAccMediaAsset_base(rName), rName, description),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "media_type", "MP4"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "labels", "test_label_1,test_lable_2,test_label_3"),
					resource.TestCheckResourceAttr(resourceName, "category_id", "-1"),
					resource.TestCheckResourceAttr(resourceName, "media_name", acceptance.HW_VOD_MEDIA_ASSET_FILE),
				),
			},
			{
				Config: testAccMediaAsset_basic(testAccMediaAsset_base(rName), updateName, descriptionUpdate),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "media_type", "MP4"),
					resource.TestCheckResourceAttr(resourceName, "description", descriptionUpdate),
					resource.TestCheckResourceAttr(resourceName, "labels", "test_label_1,test_lable_2,test_label_3"),
					resource.TestCheckResourceAttr(resourceName, "category_id", "-1"),
					resource.TestCheckResourceAttr(resourceName, "media_name", acceptance.HW_VOD_MEDIA_ASSET_FILE),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"input_bucket", "input_path", "thumbnail",
				},
			},
		},
	})
}

func testAccMediaAsset_base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "test" {
  bucket = "%s"
  acl    = "private"
}

resource "huaweicloud_obs_bucket_object" "test" {
  bucket = huaweicloud_obs_bucket.test.bucket
  key    = "input/%[2]s"
  source = "%[2]s"
}`, rName, acceptance.HW_VOD_MEDIA_ASSET_FILE)
}

func testAccMediaAsset_basic(baseConfig, rName, description string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vod_media_asset" "test" {
  name         = "%s"
  media_type   = "MP4"
  input_bucket = huaweicloud_obs_bucket.test.bucket
  input_path   = huaweicloud_obs_bucket_object.test.id
  description  = "%s"
  labels       = "test_label_1,test_lable_2,test_label_3"

  thumbnail {
    type = "time"
    time = 1
  }
}
`, baseConfig, rName, description)
}
