volume_type        = "GPSSD"
volume_description = "Created by terraform"
volume_name        = "tf_test_volume"

volume_metadata = {
  test = "terraform volume"
}

volume_tags = {
  foo = "bar"
  key = "value"
}

snapshot_name        = "tf_test_snapshot"
snapshot_description = "Created by terraform"

snapshot_metadata = {
  test = "terraform snapshot"
}
