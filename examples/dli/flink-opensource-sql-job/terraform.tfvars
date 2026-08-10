elastic_resource_pool_name = "tf_test_resource_pool"
elastic_resource_pool_cidr = "172.16.0.0/18"

elastic_resource_pool_label = {
  spec = "basic"
}

queue_name        = "tf_test_queue"
job_name          = "tf_test_flink_sql_job"
job_flink_version = "1.15"
job_sql           = <<-EOF
create table dataGenSource(
  user_id string,
  amount int
) with (
  'connector' = 'datagen',
  'rows-per-second' = '1',
  'fields.user_id.kind' = 'random',
  'fields.user_id.length' = '3'
);

create table printSink(
  user_id string,
  amount int
) with (
  'connector' = 'print'
);

insert into printSink select * from dataGenSource;
EOF
