resource "huaweicloud_tms_tags" "test" {
	tag {
		key = var.tags_key
		value = var.tags_value
	}
}