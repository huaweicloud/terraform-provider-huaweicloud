# Basic VPC peering

This example provisions a VPC peering connection is used to connect subnet A in VPC A to subnet B in VPC B,
the local VPC is VPC A, and the peer VPC is VPC B.

* a VPC with a subnet:
  > Create a VPC A with subnet A and create a VPC B with subnet B.

* a VPC peering:
  > The local VPC is VPC A, and the peer VPC is VPC B, VPC A and VPC B cidrs are the same. Add a route in the route
  table of VPC B, set destination to 0.0.0.0/0 and next hop to the created VPC peering connection between VPC A and
  VPC B.

## Usage

```shell
terraform init
terraform plan
terraform apply
terraform destroy
```
