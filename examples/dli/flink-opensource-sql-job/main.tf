resource "huaweicloud_dli_elastic_resource_pool" "test" {
  name                  = var.elastic_resource_pool_name
  description           = var.elastic_resource_pool_description
  min_cu                = var.elastic_resource_pool_min_cu
  max_cu                = var.elastic_resource_pool_max_cu
  cidr                  = var.elastic_resource_pool_cidr
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
  label                 = var.elastic_resource_pool_label
}

resource "huaweicloud_dli_queue" "test" {
  elastic_resource_pool_name = huaweicloud_dli_elastic_resource_pool.test.name
  resource_mode              = 1
  name                       = var.queue_name
  queue_type                 = "general"
  cu_count                   = var.queue_cu_count
  description                = var.queue_description
  enterprise_project_id      = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}

resource "huaweicloud_dli_flinksql_job" "test" {
  name                   = var.job_name
  type                   = "flink_opensource_sql_job"
  run_mode               = "exclusive_cluster"
  queue_name             = huaweicloud_dli_queue.test.name
  sql                    = var.job_sql
  flink_version          = var.job_flink_version
  execution_agency_urn   = var.job_execution_agency_urn
  cu_number              = var.job_cu_number
  parallel_number        = var.job_parallel_number
  manager_cu_number      = var.job_manager_cu_number
  tm_cus                 = var.job_tm_cus
  tm_slot_num            = var.job_tm_slot_num
  checkpoint_enabled     = var.job_checkpoint_enabled
  checkpoint_mode        = var.job_checkpoint_mode
  checkpoint_interval    = var.job_checkpoint_interval
  obs_bucket             = var.job_obs_bucket
  log_enabled            = var.job_log_enabled
  smn_topic              = var.job_smn_topic
  restart_when_exception = var.job_restart_when_exception
  resume_max_num         = var.job_resume_max_num
  idle_state_retention   = var.job_idle_state_retention
  runtime_config         = var.job_runtime_config
  tags                   = var.job_tags
  description            = var.job_description

  depends_on = [huaweicloud_dli_queue.test]
}
