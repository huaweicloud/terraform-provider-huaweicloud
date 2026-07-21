output "top_sql_infos" {
  description = "The list of Top SQL information"
  value       = data.huaweicloud_gaussdb_top_sql_statements.test.top_sql_infos
}
