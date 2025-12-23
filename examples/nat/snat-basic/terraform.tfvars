# VPC Configuration
vpc_name = "terraform-test-vpc"
vpc_cidr = "192.168.0.0/16"

# Subnet Configuration
subnet_name       = "terraform-test-subnet"
subnet_cidr       = "192.168.1.0/24"
subnet_gateway_ip = "192.168.1.1"

# NAT Gateway Configuration
nat_gateway_name    = "terraform-test-nat-gateway"
gateway_spec        = "1"
gateway_description = "NAT gateway for terraform test"

# EIP Configuration
eip_bandwidth_name        = "terraform-test-eip-bandwidth"
eip_bandwidth_size        = 5
eip_bandwidth_share_type  = "PER"
eip_bandwidth_charge_mode = "traffic"

# SNAT Rule Configuration (VPC Scenario)
snat_source_type = 0
snat_cidr        = "192.168.0.0/16"
snat_description = "SNAT rule for terraform test"

# ECS Instance Configuration
ecs_instance_name           = "terraform-test-ecs"
ecs_flavor_id               = ""
ecs_flavor_performance_type = "normal"
ecs_flavor_cpu_core_count   = 2
ecs_flavor_memory_size      = 4
ecs_image_id                = ""
ecs_image_visibility        = "public"
ecs_image_os                = "Ubuntu"
ecs_security_group_name     = "terraform-test-ecs-sg"
ecs_admin_password          = "YourPassword123!"
ecs_system_disk_type        = "GPSSD"
ecs_system_disk_size        = 40
ecs_instance_tags           = {
  Environment = "test"
  Project     = "terraform-example"
}
