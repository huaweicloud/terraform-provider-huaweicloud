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

resource "huaweicloud_dli_flinkjar_job" "test" {
  name                   = var.job_name
  description            = var.job_description
  queue_name             = huaweicloud_dli_queue.test.name
  main_class             = var.job_main_class
  entrypoint             = var.job_entrypoint
  entrypoint_args        = var.job_entrypoint_args
  dependency_jars        = var.job_dependency_jars
  dependency_files       = var.job_dependency_files
  execution_agency_urn   = var.job_execution_agency_urn
  feature                = var.job_feature
  flink_version          = var.job_flink_version
  image                  = var.job_image
  cu_num                 = var.job_cu_num
  parallel_num           = var.job_parallel_num
  obs_bucket             = var.job_obs_bucket
  log_enabled            = var.job_log_enabled
  smn_topic              = var.job_smn_topic
  restart_when_exception = var.job_restart_when_exception
  manager_cu_num         = var.job_manager_cu_num
  tm_cu_num              = var.job_tm_cu_num
  tm_slot_num            = var.job_tm_slot_num
  resume_checkpoint      = var.job_resume_checkpoint
  resume_max_num         = var.job_resume_max_num
  runtime_config         = var.job_runtime_config
  checkpoint_path        = var.job_checkpoint_path
  tags                   = var.job_tags
  checkpoint_enabled     = var.job_checkpoint_enabled
  checkpoint_mode        = var.job_checkpoint_mode
  checkpoint_interval    = var.job_checkpoint_interval

  depends_on = [huaweicloud_dli_queue.test]
}
