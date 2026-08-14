data "huaweicloud_das_database_users" "test" {
  instance_id = var.log_analysis_instance_id
}

locals {
  user_id = try(data.huaweicloud_das_database_users.test.users[0].id, null)
}

# Export slow logs to OBS bucket
resource "huaweicloud_das_slow_log_export_task" "test" {
  instance_id = var.log_analysis_instance_id
  bucket_name = var.slow_log_bucket_name
  start_time  = var.slow_log_start_time
  end_time    = var.slow_log_end_time
  file_path   = var.slow_log_file_path
  export_type = var.slow_log_export_type
  sort_field  = var.slow_log_sort_field
  sort_asc    = var.slow_log_sort_asc
  time_zone   = var.slow_log_time_zone
}

# Create a binlog parse task
resource "huaweicloud_das_binlog_parse_task" "test" {
  user_id     = local.user_id
  binlog_type = var.binlog_binlog_type
  file_name   = var.binlog_file_name
  backup_id   = var.binlog_backup_id
}

# Export binlog parse task results to OBS bucket
resource "huaweicloud_das_binlog_parse_task_export" "test" {
  user_id     = local.user_id
  task_id     = huaweicloud_das_binlog_parse_task.test.id
  bucket_name = var.binlog_export_bucket_name

  filter_condition {
    db_names            = var.binlog_filter_db_names
    tb_names            = var.binlog_filter_tb_names
    start_time          = var.binlog_filter_start_time
    end_time            = var.binlog_filter_end_time
    types               = var.binlog_filter_types
    parse_double_insert = var.binlog_filter_parse_double_insert
  }
}
