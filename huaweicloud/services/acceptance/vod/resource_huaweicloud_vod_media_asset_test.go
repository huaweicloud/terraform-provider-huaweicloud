package vod

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/vod"
)

func getResourceAsset(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "vod"
	)

	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating VOD client: %s", err)
	}

	return vod.ReadMediaAssetDetail(client, state.Primary.ID)
}

func TestAccMediaAsset_obs(t *testing.T) {
	var asset interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
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
				Config: testAccMediaAsset_obs(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "media_type", "MP4"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "labels", "test_label_1,test_lable_2,test_label_3"),
					resource.TestCheckResourceAttrPair(resourceName, "input_bucket",
						"huaweicloud_obs_bucket.test", "bucket"),
					resource.TestCheckResourceAttrPair(resourceName, "input_path",
						"huaweicloud_obs_bucket_object.test", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "category_id"),
					resource.TestCheckResourceAttrSet(resourceName, "category_name"),
					resource.TestCheckResourceAttrSet(resourceName, "media_name"),
					resource.TestCheckResourceAttrSet(resourceName, "thumbnail.#"),
				),
			},
			{
				Config: testAccMediaAsset_obs_update1(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s_update", rName)),
					resource.TestCheckResourceAttr(resourceName, "media_type", "MP4"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description update"),
					resource.TestCheckResourceAttr(resourceName, "labels", "test_label_1,test_lable_2"),
					resource.TestCheckResourceAttr(resourceName, "publish", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "category_id",
						"huaweicloud_vod_media_category.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "input_bucket",
						"huaweicloud_obs_bucket.test", "bucket"),
					resource.TestCheckResourceAttrPair(resourceName, "input_path",
						"huaweicloud_obs_bucket_object.test", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "category_name"),
					resource.TestCheckResourceAttrSet(resourceName, "media_name"),
					resource.TestCheckResourceAttrSet(resourceName, "thumbnail.#"),
				),
			},
			{
				Config: testAccMediaAsset_obs_update2(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s_update", rName)),
					resource.TestCheckResourceAttr(resourceName, "media_type", "MP4"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "labels", "test_label_1,test_lable_2"),
					resource.TestCheckResourceAttr(resourceName, "publish", "false"),
					resource.TestCheckResourceAttrPair(resourceName, "category_id",
						"huaweicloud_vod_media_category.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "input_bucket",
						"huaweicloud_obs_bucket.test", "bucket"),
					resource.TestCheckResourceAttrPair(resourceName, "input_path",
						"huaweicloud_obs_bucket_object.test", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "category_name"),
					resource.TestCheckResourceAttrSet(resourceName, "media_name"),
					resource.TestCheckResourceAttrSet(resourceName, "thumbnail.#"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"input_bucket", "input_path", "thumbnail", "publish",
				},
			},
		},
	})
}

func testAccMediaAsset_obs_base(rName string) string {
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

func testAccMediaAsset_obs(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vod_media_asset" "test" {
  name         = "%[2]s"
  media_type   = "MP4"
  input_bucket = huaweicloud_obs_bucket.test.bucket
  input_path   = huaweicloud_obs_bucket_object.test.id
  description  = "test description"
  labels       = "test_label_1,test_lable_2,test_label_3"

  thumbnail {
    type = "time"
    time = 1
  }
}
`, testAccMediaAsset_obs_base(rName), rName)
}

func testAccMediaAsset_obs_update1(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vod_media_category" "test" {
  name = "%[2]s"
}

resource "huaweicloud_vod_media_asset" "test" {
  name         = "%[2]s_update"
  media_type   = "MP4"
  input_bucket = huaweicloud_obs_bucket.test.bucket
  input_path   = huaweicloud_obs_bucket_object.test.id
  description  = "test description update"
  category_id  = huaweicloud_vod_media_category.test.id
  labels       = "test_label_1,test_lable_2"
  publish      = true

  thumbnail {
    type = "time"
    time = 1
  }
}
`, testAccMediaAsset_obs_base(rName), rName)
}

func testAccMediaAsset_obs_update2(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vod_media_category" "test" {
  name = "%[2]s"
}

resource "huaweicloud_vod_media_asset" "test" {
  name         = "%[2]s_update"
  media_type   = "MP4"
  input_bucket = huaweicloud_obs_bucket.test.bucket
  input_path   = huaweicloud_obs_bucket_object.test.id
  category_id  = huaweicloud_vod_media_category.test.id
  labels       = "test_label_1,test_lable_2"
  publish      = false

  thumbnail {
    type = "time"
    time = 1
  }
}
`, testAccMediaAsset_obs_base(rName), rName)
}

func TestAccMediaAsset_url(t *testing.T) {
	var asset interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
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
				Config: testAccMediaAsset_url(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "media_type", "MP4"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "url", "http://demo/test_video.mp4"),
					resource.TestCheckResourceAttr(resourceName, "labels", "test_label_1,test_lable_2,test_label_3"),
					resource.TestCheckResourceAttrSet(resourceName, "category_id"),
					resource.TestCheckResourceAttrSet(resourceName, "category_name"),
					resource.TestCheckResourceAttrSet(resourceName, "media_name"),
					resource.TestCheckResourceAttrSet(resourceName, "thumbnail.#"),
				),
			},
			{
				Config: testAccMediaAsset_obs_url1(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s_update", rName)),
					resource.TestCheckResourceAttr(resourceName, "media_type", "MP4"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description update"),
					resource.TestCheckResourceAttr(resourceName, "url", "http://demo/test_video.mp4"),
					resource.TestCheckResourceAttr(resourceName, "labels", "test_label_1,test_lable_2"),
					resource.TestCheckResourceAttr(resourceName, "publish", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "category_id",
						"huaweicloud_vod_media_category.test", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "category_name"),
					resource.TestCheckResourceAttrSet(resourceName, "media_name"),
					resource.TestCheckResourceAttrSet(resourceName, "thumbnail.#"),
				),
			},
			{
				Config: testAccMediaAsset_obs_url2(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s_update", rName)),
					resource.TestCheckResourceAttr(resourceName, "media_type", "MP4"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "url", "http://demo/test_video.mp4"),
					resource.TestCheckResourceAttr(resourceName, "labels", "test_label_1,test_lable_2"),
					resource.TestCheckResourceAttr(resourceName, "publish", "false"),
					resource.TestCheckResourceAttrPair(resourceName, "category_id",
						"huaweicloud_vod_media_category.test", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "category_name"),
					resource.TestCheckResourceAttrSet(resourceName, "media_name"),
					resource.TestCheckResourceAttrSet(resourceName, "thumbnail.#"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"url", "thumbnail", "publish",
				},
			},
		},
	})
}

func testAccMediaAsset_url(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vod_media_asset" "test" {
  name        = "%[1]s"
  media_type  = "MP4"
  url         = "http://demo/test_video.mp4"
  description = "test description"
  labels      = "test_label_1,test_lable_2,test_label_3"

  thumbnail {
    type = "time"
    time = 1
  }
}
`, rName)
}

func testAccMediaAsset_obs_url1(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vod_media_category" "test" {
  name = "%[1]s"
}

resource "huaweicloud_vod_media_asset" "test" {
  name        = "%[1]s_update"
  media_type  = "MP4"
  url         = "http://demo/test_video.mp4"
  description = "test description update"
  category_id = huaweicloud_vod_media_category.test.id
  labels      = "test_label_1,test_lable_2"
  publish     = true

  thumbnail {
    type = "time"
    time = 1
  }
}
`, rName)
}

func testAccMediaAsset_obs_url2(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vod_media_category" "test" {
  name = "%[1]s"
}

resource "huaweicloud_vod_media_asset" "test" {
  name        = "%[1]s_update"
  media_type  = "MP4"
  url         = "http://demo/test_video.mp4"
  category_id = huaweicloud_vod_media_category.test.id
  labels      = "test_label_1,test_lable_2"
  publish     = false

  thumbnail {
    type = "time"
    time = 1
  }
}
`, rName)
}
