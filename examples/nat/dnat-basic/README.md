## Example: Basic NAT gateway and dnat rule

This example provisions:
- a ECS instance.
  > It can mount one or more data disks
  > Using CentOS 7.3 64bit as the operating system image
- a basic NAT gateway.
- a dnat rule.
  > DNAT rule specifies port 8080 as the port for Tomcat to provide external services, which can be used to build Java Web services.
  > DNAT function is bound with EIP, and EIP is shared across VPC by binding IP mapping, which provides services for the Internet.
- a VPC and subnet instance.
  > VPC provides an isolated cloud service environment.
  