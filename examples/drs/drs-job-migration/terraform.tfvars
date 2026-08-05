vpc_name            = "your_vpc"
subnet_name         = "your_subnet"
security_group_name = "your_security_group"
source_rds_name     = "your_source_rds"
dest_rds_name       = "your_dest_rds"
rds_flavor          = "rds.mysql.x1.large.2.ha"
source_rds_fixed_ip = "192.168.0.58"
dest_rds_fixed_ip   = "192.168.0.59"
db_password         = "TestDrs@123"
job_name            = "your_drs_job"
tags                = {
  foo = "bar"
  key = "value"
}
