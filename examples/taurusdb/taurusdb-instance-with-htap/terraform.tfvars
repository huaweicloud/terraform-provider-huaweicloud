vpc_name               = "your_vpc"
subnet_name            = "your_subnet"
security_group_name    = "your_security_group"
taurusdb_instance_name = "your_taurusdb_instance"
htap_instance_name     = "your_htap_starrocks_instance"
be_parameter_values    = {
  alter_tablet_worker_count            = "1"
  base_compaction_num_threads_per_disk = "1"
}
fe_parameter_values    = {
  alter_table_timeout_second     = "21600"
  bdbje_heartbeat_timeout_second = "10"
}
