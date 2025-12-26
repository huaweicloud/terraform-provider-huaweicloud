security_group_rules = [
  {
    direction        = "ingress"
    description      = "Allow SSH access"
    ethertype        = "IPv4"
    protocol         = "tcp"
    ports            = "22"
    remote_ip_prefix = "1.2.3.4/32"  # Your IP address
  },
  {
    direction        = "ingress"
    description      = "Allow HTTP and HTTPS access"
    ethertype        = "IPv4"
    protocol         = "tcp"
    ports            = "80,443"
    remote_ip_prefix = "0.0.0.0/0"
  },
  {
    direction        = "ingress"
    description      = "Allow ICMP access"
    ethertype        = "IPv4"
    protocol         = "icmp"
    ports            = null
    remote_ip_prefix = "0.0.0.0/0"
  }
]

ecs_admin_password = "YourSecurePassword123!"

# VPC configurations
vpcs = [
  {
    vpc_name          = "my-vpc-one"
    vpc_cidr          = "192.168.0.0/16"
    subnet_name       = "my-subnet-one"
    subnet_cidr       = "192.168.1.0/24"
    subnet_gateway_ip = "192.168.1.1"
    instance_name     = "my-instance-one"
  },
  {
    vpc_name          = "my-vpc-other"
    vpc_cidr          = "10.0.0.0/16"
    subnet_name       = "my-subnet-other"
    subnet_cidr       = "10.0.2.0/24"
    subnet_gateway_ip = "10.0.2.1"
    instance_name     = "my-instance-other"
  }
]

# NAT Gateway configuration
nat_gateway_name = "my-nat-gateway"
bandwidth_size   = 10
