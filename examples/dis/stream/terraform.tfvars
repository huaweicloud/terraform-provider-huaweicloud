stream_name                           = "tf_test_dis_stream"
stream_partition_count                = 2
stream_auto_scale_min_partition_count = 2
stream_auto_scale_max_partition_count = 4
stream_type                           = "COMMON"
stream_compression_format             = "zip"
stream_data_type                      = "CSV"
stream_csv_delimiter                  = ";"
stream_data_schema                    = "{\"type\":\"record\",\"name\":\"RecordName\",\"fields\":[{\"type\":\"string\",\"name\":\"name\"}]}"
stream_tags                           = {
  foo = "bar"
  key = "value"
}
