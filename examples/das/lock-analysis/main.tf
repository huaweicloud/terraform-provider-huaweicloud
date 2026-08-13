# Enable full dead lock detection for a database instance
resource "huaweicloud_das_full_dead_lock_switch" "test" {
  instance_id = var.lock_analysis_instance_id

  switch_on = var.full_dead_lock_switch_on

  retention_hours = var.full_dead_lock_retention_hours
}

# Enable history transaction switch for a database instance
resource "huaweicloud_das_history_transaction_switch" "test" {
  instance_id    = var.lock_analysis_instance_id
  status         = var.history_transaction_status
  datastore_type = var.lock_analysis_datastore_type
}

# Export history transactions to OBS bucket
resource "huaweicloud_das_history_transaction_export_task" "test" {
  instance_id = var.lock_analysis_instance_id
  bucket_name = var.history_transaction_bucket_name
  start_time  = var.history_transaction_start_time
  end_time    = var.history_transaction_end_time
  file_path   = var.history_transaction_file_path
  time_zone   = var.history_transaction_time_zone
  order_field = var.history_transaction_order_field
  order_by    = var.history_transaction_order_by

  last_sec_min = var.history_transaction_last_sec_min
  last_sec_max = var.history_transaction_last_sec_max

  depends_on = [huaweicloud_das_history_transaction_switch.test]
}
