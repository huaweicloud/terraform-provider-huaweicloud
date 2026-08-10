# Create a DIS stream with auto scaling and data format configuration
resource "huaweicloud_dis_stream" "test" {
  stream_name      = var.stream_name
  partition_count  = var.stream_partition_count
  stream_type      = var.stream_type
  retention_period = var.stream_retention_period

  auto_scale_min_partition_count = var.stream_auto_scale_min_partition_count
  auto_scale_max_partition_count = var.stream_auto_scale_max_partition_count

  compression_format = var.stream_compression_format
  data_type          = var.stream_data_type
  csv_delimiter      = var.stream_csv_delimiter
  data_schema        = var.stream_data_schema

  tags = var.stream_tags
}
