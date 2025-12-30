resource "huaweicloud_tms_tags" "test" {
  dynamic "tags" {
    for_each = var.preset_tags

    content {
      key   = tags.value.key
      value = tags.value.value
    }
  }
}
