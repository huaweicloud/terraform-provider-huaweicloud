# Batch enable SQL switch for multiple instances
resource "huaweicloud_das_batch_set_sql_switch" "test" {
  engine_type     = var.sql_diagnostics_engine_type
  switch_on       = var.batch_sql_switch_on
  switch_type     = var.batch_sql_switch_type
  instance_ids    = var.batch_sql_instance_ids
  retention_hours = var.batch_sql_retention_hours
}

# Enable search path switch for a database connection
resource "huaweicloud_das_search_path_switch" "test" {
  connection_id = var.search_path_connection_id
  switch_on     = var.search_path_switch_on
}
