# Create a basic AS configuration

This example creates a basic AS configuration based on the example
[examples/as/configuration-basic](https://github.com/huaweicloud/terraform-provider-huaweicloud/tree/master/examples/as/configuration-basic).
You can replace the VPC, Subnet and Security Group with the resources already created in HuaweiCloud.

To run, configure your HuaweiCloud provider as described in the
[document](https://registry.terraform.io/providers/huaweicloud/huaweicloud/latest/docs).

## The basic AS configuration

| Attributes       | Value |
|------------------|-------|
| Disk size        | 40    |
| Disk volume type | SSD   |
| Disk type        | SYS   |

## Usage

```shell
terraform init
terraform plan
terraform apply
terraform destroy
```

It takes several minutes to create a basic AS configuration.

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 0.12.0 |
| huaweicloud | >= 1.49.0 |
