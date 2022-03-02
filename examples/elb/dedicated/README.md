# Configure a dedicated Elastic Load Balance service

Configuration in this directory creates an elastic load balancing service including elastic load balancer, listener,
backend server group, and bind the backend servers and enable health check for each servers.
This sample contains the ACL and security group rules (in_v4_elb_member) that must be configured to run the elastic
load balancing service normally.

To check the health of backend servers, dedicated load balancers use the IP addresses from the VPC where they work to
send heartbeat requests to backend servers. Therefore, the corresponding network segment needs to be opened in the ACL
and security group configuration.

To run, configure your Huaweicloud provider as described in the
[document](https://registry.terraform.io/providers/huaweicloud/huaweicloud/latest/docs).

If you want to use key-pair or password to create the ECS instance, please refer the examples of the ECS service.

## Usage

```shell
terraform init
terraform plan
terraform apply
terraform destroy
```

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.33.0 |
