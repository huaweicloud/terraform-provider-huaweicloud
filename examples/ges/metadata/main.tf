# Create an OBS bucket for storing GES metadata schema files
resource "huaweicloud_obs_bucket" "test" {
  bucket        = var.bucket_name
  acl           = "private"
  force_destroy = true
}

# Create a GES metadata
resource "huaweicloud_ges_metadata" "test" {
  name          = var.metadata_name
  description   = var.metadata_description
  metadata_path = "${huaweicloud_obs_bucket.test.bucket}/${var.metadata_schema_file}"

  ges_metadata {
    labels {
      name       = "user"
      properties = var.metadata_properties
    }
  }

  depends_on = [
    huaweicloud_obs_bucket.test
  ]
}
