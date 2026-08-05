vpc_name                    = "your_vpc1"
subnet_name                 = "your_subnet"
security_group_name         = "your_security_group"
instance_name               = "your_geminidb_influxdb_instance"
instance_backup_time_window = "03:00-04:00"
instance_backup_keep_days   = 14
tags                        = {
  foo = "bar"
  key = "value"
}
