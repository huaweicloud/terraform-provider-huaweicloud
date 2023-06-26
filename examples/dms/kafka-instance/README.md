# Create a Kafka cluster

This example creates a Kafka cluster based on the example
[examples/dms/kafka-instance](https://github.com/huaweicloud/terraform-provider-huaweicloud/tree/master/examples/dms/kafka-instance).
You can replace the VPC, Subnet and Security Group with the resources already created in HuaweiCloud.

To run, configure your HuaweiCloud provider as described in the
[document](https://registry.terraform.io/providers/huaweicloud/huaweicloud/latest/docs).

## The Kafka instance configuration

| Attributes    | Value           |
|---------------|-----------------|
| Instance Type | cluster         |
| Version       | 2.7             |
| Flavor        | c6.2u4g.cluster |
| Broker Number | 3               |
| storage_space | 600 GB          |
| IO Type       | ultra-high I/O  |

### FAQ

- **How to obtain the TPS, and number of partitions of the Kafka instances?**

  The **TPS**, and the **number of partitions** are included in the flavor and do not need to be
  specified when creating a resource.

  The expected `flavor_id` can be obtained in the following way.

  ```hcl
  data "huaweicloud_dms_kafka_flavors" "test" {
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

It takes about 20 to 50 minutes to create a Kafka cluster depending on the flavor and broker number.

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 0.12.0 |
| huaweicloud | >= 1.40.0 |
