vpc_name            = "tf_test_vpc"
subnet_name         = "tf_test_subnet"
security_group_name = "tf_test_security_group"

security_group_rule_configurations = [
  # Allow all IPv4 ingress traffic of the ICMP protocol
  {
    direction        = "ingress"
    ethertype        = "IPv4"
    protocol         = "icmp"
    remote_ip_prefix = "0.0.0.0/0"
  },
  # Allow some ports for IPv4 ingress traffic of the TCP protocol
  {
    direction        = "ingress"
    ethertype        = "IPv4"
    protocol         = "tcp"
    ports            = "22,3389"
    remote_ip_prefix = "10.1.0.7/32"
  },
  # Allow all IPv4 egress traffic
  {
    direction        = "egress"
    ethertype        = "IPv4"
    remote_ip_prefix = "0.0.0.0/0"
  },
]

instance_name                   = "tf_test_ecs"
internet_gateway_name           = "tf_test_igw"
global_eip_name                 = "tf_test_geip"
internet_bandwidth_name         = "tf_test_internet_bandwidth"
gc_bandwidth_name               = "tf_test_gc_bandwidth"
instance_administrator_password = "YourPassword@123"
