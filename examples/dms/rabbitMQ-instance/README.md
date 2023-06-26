# Create a RabbitMQ cluster

This example creates a RabbitMQ cluster based on the example
[examples/dms/rabbitMQ-instance](https://github.com/huaweicloud/terraform-provider-huaweicloud/tree/master/examples/dms/rabbitMQ-instance).
You can replace the VPC, Subnet and Security Group with the resources already created in HuaweiCloud.

To run, configure your HuaweiCloud provider as described in the
[document](https://registry.terraform.io/providers/huaweicloud/huaweicloud/latest/docs).

## The RabbitMQ instance configuration

| Attributes    | Value           |
|---------------|-----------------|
| Version       | 3.8.35          |
| Instance Type | cluster         |
| Flavor        | c6.2u4g.cluster |
| Broker Number | 3               |
| Storage Space | 600 GB          |
| IO Type       | ultra-high I/O  |

### FAQ

- **How to obtain the number of nodes of the RabbitMQ instances?**

  The **number of nodes** is included in the product and do not need to be
  specified when creating a resource.

  The expected `flavor_id` can be obtained in the following way.

  ```hcl
  data "huaweicloud_dms_rabbitmq_flavors" "test" {
    type              = "cluster"
    flavor_id         = "c6.2u4g.cluster"
    storage_spec_code = "dms.physical.storage.ultra.v2"
  }
  ```

## Usage

```shell
terraform init
terraform plan
terraform apply
terraform destroy
```

It takes about 20 to 50 minutes to create a RabbitMQ cluster depending on the flavor and broker number.

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 0.12.0 |
| huaweicloud | >= 1.49.0 |
