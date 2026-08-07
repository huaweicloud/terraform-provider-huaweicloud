# Primary Region - Network
primary_vpc_name            = "example_vpc"
primary_vpc_cidr            = "172.16.0.0/16"
primary_subnet_names        = ["example_subnet_1", "example_subnet_2"]
primary_availability_zones  = ["primary_region_az_1", "primary_region_az_2"]
primary_security_group_name = "example_security_group"

# DR Region - Network
dr_vpc_name            = "example_vpc_dr"
dr_vpc_cidr            = "172.17.0.0/16"
dr_subnet_name         = "example_subnet_dr"
dr_availability_zone   = "dr_region_az"
dr_security_group_name = "example_security_group"
dr_region_name         = "dr_region"

# Cloud Connection
# cc_bandwidth = 10

# Primary Region - GaussDB Instance
instance_passwords                  = ["password_1", "password_2"]
primary_instance_name               = "example_primary_instance"
instance_flavor                     = "gaussdb.opengauss.ee.c3.xlarge.x864.ha"
primary_instance_availability_zones = "primary_region_az_1,primary_region_az_2,primary_region_az_3"
instance_db_port                    = 5432
enterprise_project_id               = "your_enterprise_project_id"
primary_instance_volume_type        = "ULTRAHIGH"
primary_instance_volume_size        = 40

# DR Region - GaussDB Instance
dr_instance_name               = "example_dr_instance"
dr_instance_availability_zones = "dr_region_az_1,dr_region_az_2,dr_region_az_3"
dr_instance_volume_type        = "ULTRAHIGH"
dr_instance_volume_size        = 40

# Disaster Recovery
dr_user_name     = "root"
dr_user_password = "your_dr_user_password"
