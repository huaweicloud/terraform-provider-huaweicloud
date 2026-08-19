resource "huaweicloud_vod_media_category" "test" {
  name = var.media_category_name
}

resource "huaweicloud_vod_media_asset" "test" {
  name        = var.media_asset_name
  media_type  = "MP4"
  url         = var.media_asset_url
  description = var.media_asset_description
  category_id = huaweicloud_vod_media_category.test.id
  labels      = var.media_asset_labels
}
