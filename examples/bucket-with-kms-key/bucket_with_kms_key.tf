
provider "huaweicloud" {
  user_name   = "${var.user_name}"
  tenant_name = "${var.tenant_name}"
  password    = "${var.password}"
  auth_url    = "${var.auth_url}"
  region      = "${var.region}"
  domain_name = "${var.domain_name}"
  access_key  = "${var.access_key}"
  secret_key  = "${var.secret_key}"
  insecure    = "true"
}


resource "huaweicloud_kms_key_v1" "kms_key1" {
  key_alias       = "test_kms_key_1_for_obs"
  key_description = "des test"
  realm           = "${var.region}"
  is_enabled      = true
  pending_days = "7"
}


resource "huaweicloud_s3_bucket" "bucket1" {
  region = "${var.region}"
  bucket = "test-bucket-with-kms-key-1234567890000"
  acl    = "public-read-write"
}

resource "huaweicloud_s3_bucket_object" "object1" {
  source = "./s3_file_test.txt"
  bucket = "${huaweicloud_s3_bucket.bucket1.id}"
  key    = "test-bucket-testfile"
  acl    = "public-read"
  server_side_encryption = "aws:kms"
  sse_kms_key_id         = "${huaweicloud_kms_key_v1.kms_key1.id}"
  depends_on             = ["huaweicloud_kms_key_v1.kms_key1"]
}
