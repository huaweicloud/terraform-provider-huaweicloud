data "huaweicloud_gaussdb_top_sql_statements" "test" {
  instance_id    = var.instance_id
  node_ids       = var.node_ids
  start_time     = var.start_time
  end_time       = var.end_time
  support_system = var.support_system

  dynamic "multi_queries" {
    for_each = var.multi_queries

    content {
      name      = multi_queries.value.name
      condition = multi_queries.value.condition
      values    = multi_queries.value.values
      is_fuzzy  = multi_queries.value.is_fuzzy
    }
  }
}
