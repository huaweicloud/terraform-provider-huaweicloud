resource "huaweicloud_kms_key" "test" {
  count = var.bucket_encryption && var.bucket_encryption_key_id == "" ? 1 : 0

  key_alias = var.key_alias
  key_usage = var.key_usage
}

locals {
  website_page_names = [for k, v in var.website_configurations : v.file_name if k == "index" || k == "error"]
  index_page         = lookup(var.website_configurations, "index")
  error_page         = lookup(var.website_configurations, "error")
}

resource "huaweicloud_obs_bucket" "test" {
  bucket        = var.bucket_name
  storage_class = var.bucket_storage_class
  acl           = var.bucket_acl
  encryption    = var.bucket_encryption
  sse_algorithm = var.bucket_encryption ? var.bucket_sse_algorithm : null
  kms_key_id    = var.bucket_encryption ? var.bucket_encryption_key_id != "" ? var.bucket_encryption_key_id : huaweicloud_kms_key.test[0].id : null
  force_destroy = var.bucket_force_destroy
  tags          = merge(var.bucket_tags, {for k, v in var.website_configurations : k => v.file_name if k == "index" || k == "error"})

  provisioner "local-exec" {
    command = "echo '${lookup(local.index_page, "content")}' >> ${lookup(local.index_page, "file_name")}\necho '${lookup(local.error_page, "content")}' >> ${lookup(local.error_page, "file_name")}"
  }
  provisioner "local-exec" {
    command = "rm ${self.tags.index} ${self.tags.error}"
    when    = destroy
  }

  website {
    index_document = lookup(local.index_page, "file_name")
    error_document = lookup(local.error_page, "file_name")
  }

  lifecycle {
    ignore_changes = [
      sse_algorithm
    ]
  }
}

resource "huaweicloud_obs_bucket_policy" "test" {
  bucket = huaweicloud_obs_bucket.test.id
  policy = <<POLICY
{
  "Statement": [
    {
      "Sid": "AddPerm",
      "Effect": "Allow",
      "Principal": {"ID": "*"},
      "Action": ["GetObject"],
      "Resource": "${huaweicloud_obs_bucket.test.id}/*"
    }
  ]
}
  POLICY
}

resource "huaweicloud_obs_bucket_object" "test" {
  count = length(local.website_page_names)

  bucket = huaweicloud_obs_bucket.test.id
  key    = try(local.website_page_names[count.index], null)
  source = abspath(format("./%s", try(local.website_page_names[count.index], null)))
}
