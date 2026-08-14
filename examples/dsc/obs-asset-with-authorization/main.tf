resource "huaweicloud_obs_bucket" "test" {
  bucket        = var.bucket_name
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_dsc_asset_authorization" "test" {
  type                 = "OBS"
  authorization_status = true
}

resource "huaweicloud_dsc_asset_obs" "test" {
  name          = var.asset_name
  bucket_name   = huaweicloud_obs_bucket.test.bucket
  bucket_policy = "private"

  depends_on = [huaweicloud_dsc_asset_authorization.test]
}
