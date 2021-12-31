# Create a cluster DMS RabbitMQ instances

This example creates a cluster DMS RabbitMQ instances based on the example
[examples/dms/rabbitMQ-instance](https://github.com/huaweicloud/terraform-provider-huaweicloud/tree/master/examples/dms/rabbitMQ-instance)
. You can replace the VPC, Subnet and Security Group with the resources already created in HuaweiCloud.

To run, configure your HuaweiCloud provider as described in the
[document](https://registry.terraform.io/providers/huaweicloud/huaweicloud/latest/docs).

## The RabbitMQ instance configuration

| Attributes | Value |
| ---- | ---- |
| Version | 3.7.17 |
| Instance Type | cluster |
| Nodes | 3 |
| Storage Capacity | 300GB |
| IO Type | high |

### FAQ

- **How to obtain the number of nodes of the RabbitMQ instances?**

  The **number of nodes** is included in the product and do not need to be
  specified when creating a resource. If you want to modify it, please modify the argument `product_id`.

  The expected `product_id` can be obtained in the following way.

  ```hcl
  data "huaweicloud_dms_product" "product_1" {
    engine            = "rabbitmq"
    instance_type     = "cluster"
    version           = "3.7.17"
    storage_spec_code = "dms.physical.storage.high"
  }
  ```

## Usage

```shell
terraform init
terraform plan
terraform apply
terraform destroy
```

The creation of the DMS RabbitMQ instance takes about 20 minutes.

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.31.1 |
