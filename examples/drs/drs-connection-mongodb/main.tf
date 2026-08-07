# DRS connection for MongoDB with sharding
resource "huaweicloud_drs_connection" "test" {
  name        = var.connection_name
  db_type     = "mongodb"
  description = var.description

  endpoint {
    endpoint_name = "mongodb"
    ip            = var.endpoint_ip
    db_user       = var.db_user
    db_password   = var.db_password
    db_name       = var.db_name

    source_sharding {
      endpoint_name = "mongodb"
      ip            = var.shard1_ip
      db_user       = var.db_user
      db_password   = var.db_password
      db_name       = var.db_name
    }

    source_sharding {
      endpoint_name = "mongodb"
      ip            = var.shard2_ip
      db_user       = var.db_user
      db_password   = var.db_password
      db_name       = var.db_name
    }
  }

  ssl {
    ssl_link = false
  }

  config {
    driver_name = var.driver_name
  }

  lifecycle {
    ignore_changes = [
      endpoint.0.db_password,
      endpoint.0.source_sharding.0.db_password,
      endpoint.0.source_sharding.0.endpoint_name,
      endpoint.0.source_sharding.1.db_password,
      endpoint.0.source_sharding.1.endpoint_name,
    ]
  }
}
