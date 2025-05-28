# Create a CBH basic instance

This example creates a CBH basic instance. The CBH dedicated instance requires a VPC,
Subnet and security group. In this example, we all create them in the simplest
way. You can replace them with resources already created in Huawei Cloud.

To run, configure your HuaweiCloud provider as described in the
[document](https://registry.terraform.io/providers/huaweicloud/huaweicloud/latest/docs).

## The CBH instance configuration

| Attributes    | Value           |
|---------------|-----------------|
| name          | Cbh_demo        |
| flavor_id     | cbh.basic.50    |
| password      | Cbh@Huawei123   |
| charging_mode | prePaid         |
| period_unit   | mouth           |
| period        | 1               |

## Usage

```shell
terraform init
terraform plan
terraform apply
terraform destroy
```

It takes about 15 minutes to create a CBH basic instance.

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 0.12.0 |
| huaweicloud | >= 1.40.0 |
