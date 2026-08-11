elastic_resource_pool_name = "tf_test_resource_pool"
elastic_resource_pool_cidr = "172.16.0.0/18"

elastic_resource_pool_label = {
  spec = "basic"
}

queue_name               = "tf_test_queue"
job_name                 = "tf_test_flink_jar_job"
job_entrypoint           = "obs://your_bucket_path/your.jar"
job_flink_version        = "1.15"
job_execution_agency_urn = "your_agency_urn"
