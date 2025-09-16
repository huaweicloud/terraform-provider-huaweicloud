security_group_name                = "tf_test_security_group"
security_group_rule_configurations = [
  # Allow all IPv4 ingress traffic of the ICMP protocol
  {
    direction        = "ingress"
    ethertype        = "IPv4"
    protocol         = "icmp"
    ports            = null
    remote_ip_prefix = "0.0.0.0/0"
  },
  # Allow some ports for IPv4 ingress traffic of the TCP protocol
  {
    direction        = "ingress"
    ethertype        = "IPv4"
    protocol         = "tcp"
    ports            = "22-23,443,3389,30100-30130"
    remote_ip_prefix = "0.0.0.0/0"
  },
  # Allow all IPv4 egress traffic
  {
    direction        = "egress"
    ethertype        = "IPv4"
    remote_ip_prefix = "0.0.0.0/0"
  },
  # Allow all IPv6 egress traffic
  {
    direction        = "egress"
    ethertype        = "IPv6"
    remote_ip_prefix = "::/0"
  }
]
