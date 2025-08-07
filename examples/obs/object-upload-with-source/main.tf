resource "huaweicloud_kms_key" "test" {
  count = var.bucket_encryption && var.bucket_encryption_key_id == "" ? 1 : 0

  key_alias = var.key_alias
  key_usage = var.key_usage
}

locals {
  upload_object_name   = var.object_extension_name != "" ? format("%s%s", var.object_name, var.object_extension_name) : var.object_name
  upload_zip_file_name = format("%s.zip", var.object_name)
}

resource "huaweicloud_obs_bucket" "test" {
  bucket        = var.bucket_name
  storage_class = var.bucket_storage_class
  acl           = var.bucket_acl
  encryption    = var.bucket_encryption
  sse_algorithm = var.bucket_encryption ? var.bucket_sse_algorithm : null
  kms_key_id    = var.bucket_encryption ? var.bucket_encryption_key_id != "" ? var.bucket_encryption_key_id : huaweicloud_kms_key.test[0].id : null
  force_destroy = var.bucket_force_destroy
  tags          = merge(var.bucket_tags, {
    upload_object_name   = local.upload_object_name
    upload_zip_file_name = local.upload_zip_file_name
  })

  provisioner "local-exec" {
    command = "echo '${var.object_upload_content}' >> ${local.upload_object_name}\nzip -r ${local.upload_zip_file_name} ${local.upload_object_name}"
  }
  provisioner "local-exec" {
    command = "rm ${self.tags.upload_object_name} ${self.tags.upload_zip_file_name}"
    when    = destroy
  }

  lifecycle {
    ignore_changes = [
      sse_algorithm
    ]
  }
}

resource "huaweicloud_obs_bucket_object" "test" {
  bucket = huaweicloud_obs_bucket.test.id
  key    = local.upload_zip_file_name
  source = abspath(format("./%s", local.upload_zip_file_name))
}
