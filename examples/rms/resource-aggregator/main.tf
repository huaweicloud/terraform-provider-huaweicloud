# Create a resource aggregator
resource "huaweicloud_rms_resource_aggregator" "test" {
  name        = var.aggregator_name
  type        = var.aggregator_type
  account_ids = var.account_ids
}
