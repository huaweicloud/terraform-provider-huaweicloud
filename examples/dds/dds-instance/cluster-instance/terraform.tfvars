vpc_name            = "tf_test_instance"
subnet_name         = "tf_test_instance"
security_group_name = "tf_test_instance"
instance_name       = "tf_test_instance"
instance_flavors    = [
  {
    type      = "mongos"
    num       = 2
    spec_code = "dds.mongodb.s6.large.2.mongos"
  },
  {
    type      = "shard"
    num       = 2
    spec_code = "dds.mongodb.s6.large.2.shard"
    storage   = "ULTRAHIGH"
    size      = 20
  },
  {
    type      = "config"
    num       = 1
    spec_code = "dds.mongodb.s6.large.2.config"
    storage   = "ULTRAHIGH"
    size      = 20
  }
]
