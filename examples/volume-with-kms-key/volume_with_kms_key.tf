
provider "huaweicloud" {
  user_name   = "${var.user_name}"
  tenant_name = "${var.tenant_name}"
  password    = "${var.password}"
  auth_url    = "${var.auth_url}"
  region      = "${var.region}"
  domain_name = "${var.domain_name}"
  insecure    = "true"
}


resource "huaweicloud_kms_key_v1" "kms_key1" {
  key_alias       = "test_kms_key_1"
  key_description = "des test"
  realm           = "cn-north-1"
  is_enabled      = true
  pending_days    = "7"
}

resource "huaweicloud_blockstorage_volume_v2" "volume_1" {
  size = 2
  name = "volume_with_kms_key_test"
  metadata {
    __system__encrypted = "1"
    __system__cmkid     = "${huaweicloud_kms_key_v1.kms_key1.id}"
    region              = "${var.region}"
  }
}
