# Create a cluster DMS Kafka instances

This example creates a cluster DMS kafka instances based on the example
[examples/dms/kafka-instance](https://github.com/huaweicloud/terraform-provider-huaweicloud/tree/master/examples/dms/kafka-instance)
. You can replace the VPC, Subnet and Security Group with the resources already created in HuaweiCloud.

To run, configure your HuaweiCloud provider as described in the
[document](https://registry.terraform.io/providers/huaweicloud/huaweicloud/latest/docs).

## The Kafka instance configuration

| Attributes | Value |
| ---- | ---- |
| Instance Type | cluster |
| Version | 2.3.0 |
| Bandwidth | 100MB |
| TPS | 50,000 |
| Number of Partitions | 300 |
| IO Type | high |

### FAQ

- **How to obtain the bandwidth, TPS, and number of partitions of the Kafka instances?**

  The **bandwidth**, **TPS**, and the **number of partitions** are included in the product and do not need to be
  specified when creating a resource. If you want to modify them, please modify the argument `product_id`.

  The expected `product_id` can be obtained in the following way.

  ```hcl
  data "huaweicloud_dms_product" "product_1" {
    engine            = "kafka"
    instance_type     = "cluster"
    version           = "2.3.0"
    bandwidth         = "100MB"
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

The creation of the DMS Kafka instance takes about 20 minutes.

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.31.1 |
