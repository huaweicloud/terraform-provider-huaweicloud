# Using a public NAT gateway and VPC peering to enable communications between VPCs and the internet

This example provisions VPC A and VPC B are in the same region, a public NAT gateway is configured for subnet A in
VPC A and you can add SNAT and DNAT rules for internet connectivity, subnet B connects to subnet A through a VPC
peering connection and uses the public NAT gateway of subnet A to communicate with the internet.

* A VPC:
  > Create a VPC A with subnet A, and create a VPC B with subnet B.

* A NAT gateway:
  > The NAT gateway is created in VPC A, and subnet B can use the NAT gateway to communicate the internet.

* An SNAT rule:
  > An SNAT rule for subnet A, and an SNAT rule for subnet B.

* A DNAT rule:
  > A DNAT rule for subnet A, select VPC for scenario and enter an IP address of a server in subnet A
  for private IP address. Add a DNAT rule for subnet B, set scenario to Direct Connect/Cloud Connect and
  enter an IP address of a server in subnet B for private IP address.

* An security group:
  > An security group with two ingress rules.

* An ECS instance:
  > VPC A DNAT rule associate an ECS instance and VPC B DNAT rule associate an ECS instance.

* A VPC peering:
  > A VPC peering connection is used to connect subnet A in VPC A to subnet B in VPC B, the local VPC is VPC A,
  and the peer VPC is VPC B, the VPC A and VPC B cidrs are the same. Add a route in the route table of VPC B, set
  destination to 0.0.0.0/0 and next hop to the created VPC peering connection between VPC A and VPC B.

## Usage

```shell
terraform init
terraform plan
terraform apply
terraform destroy
```
