# VPC and subnet
vpc_name          = "vpc-dnat-basic"
vpc_cidr          = "172.16.0.0/16"
subnet_name       = "subnet-dnat-basic"
subnet_cidr       = "172.16.0.0/24"
subnet_gateway_ip = "172.16.0.1"

# Backend ECS
security_group_name = "sg-dnat-backend"
ingress_cidr        = "192.168.0.0/16"
instance_name       = "ecs-dnat-backend"

# Optional: Configure data source query parameters (used when IDs are not specified)
ecs_flavor_performance_type = "normal"
ecs_flavor_cpu_core_count   = 2
ecs_flavor_memory_size      = 4
ecs_image_visibility        = "public"
ecs_image_os                = "Ubuntu"
ecs_system_disk_type        = "SSD"
ecs_system_disk_size        = 40

# DNAT ports and protocol
backend_protocol  = "tcp"
backend_port      = 2288
frontend_protocol = "tcp"
frontend_port     = 2288

# EIP and NAT gateway
eip_bandwidth_name       = "eip-dnat-basic"
eip_bandwidth_size       = 5
eip_bandwidth_share_type = "PER"
eip_bandwidth_charge_mode = "traffic"
nat_gateway_name         = "nat-gateway-dnat-basic"
nat_gateway_description  = ""
nat_gateway_spec         = "1"
