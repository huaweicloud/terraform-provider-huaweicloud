# Query all images owned by the current user
data "huaweicloud_images_images" "test" {
  visibility = "private"
  image_type = var.image_type != "" ? var.image_type : null
  os         = var.image_os != "" ? var.image_os : null
  name_regex = var.image_name_regex != "" ? var.image_name_regex : null
}

# ST.003 Disable
# Filter active images
locals {
  active_images = [
    for image in data.huaweicloud_images_images.test.images : image if image.status == "active"
  ]
}
# ST.003 Enable

# Create OBS bucket for storing exported images
resource "huaweicloud_obs_bucket" "test" {
  bucket        = var.obs_bucket_name
  storage_class = "STANDARD"
  region        = var.region_name

  tags = var.obs_bucket_tags
}

# Export each image to OBS bucket
resource "huaweicloud_ims_image_export" "test" {
  count = length(local.active_images)

  region      = var.region_name
  image_id    = local.active_images[count.index].id
  bucket_url  = "${huaweicloud_obs_bucket.test.bucket}:${local.active_images[count.index].name}-${local.active_images[count.index].id}.${var.file_format}"
  file_format = var.file_format

  depends_on = [huaweicloud_obs_bucket.test]
}
