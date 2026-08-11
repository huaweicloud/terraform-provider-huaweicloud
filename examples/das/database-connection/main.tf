# Create a database instance connection
resource "huaweicloud_das_database_instance_connection" "test" {
  instance_id      = var.connection_instance_id
  engine_type      = var.connection_engine_type
  network_type     = var.connection_network_type
  username         = var.connection_username
  password         = var.connection_password
  is_save_password = var.connection_is_save_password

  port            = var.connection_port
  database_name   = var.connection_database_name
  sql_record_flag = var.connection_sql_record_flag
  description     = var.connection_description
  node_ids        = var.connection_node_ids
}

# Create a database user
resource "huaweicloud_das_database_user" "test" {
  instance_id = var.connection_instance_id
  name        = var.db_user_name
  password    = var.db_user_password
}

# Share the connection with another IAM user
resource "huaweicloud_das_shared_connection" "test" {
  connection_id = huaweicloud_das_database_instance_connection.test.id
  user_id       = var.shared_user_id
  user_name     = var.shared_user_name
  expired_at    = var.shared_expired_at
}
