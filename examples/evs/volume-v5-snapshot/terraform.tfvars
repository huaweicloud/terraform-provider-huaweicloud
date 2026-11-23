vpc_name           = "tf_test_evs_vpc"
subnet_name        = "tf_test_evs_subnet"
secgroup_name      = "tf_test_evs_sg"
ecs_instance_name  = "tf_test_evs_ecs"
volume_name        = "tf_test_evs_volume"
volume_description = "terraform volume test"

volume_tags = {
  foo = "bar"
  key = "value"
}

volume_metadata = {
  metadata_key = "volume"
  key          = "value"
}

snapshot_name        = "tf_test_evs_v5_snapshot"
snapshot_description = "terraform v5 snapshot test"

snapshot_tags = {
  foo = "bar"
  key = "value"
}

snapshot_metadata = {
  metadata_key = "snapshot"
  key          = "value"
}
