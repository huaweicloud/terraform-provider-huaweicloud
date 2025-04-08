# Create an ER instance to share with other accounts

Configuration in this directory creates an ER instance and RAM share resource. The example includes an ER instance,
a RAM share resource, a RAM share accepter resource, a VPC configuration, an ER attachment and
an ER attachment accepter resource.
Configuration in this directory describes how to share an ER instance with other accounts, and the sharer accepts or
rejects attachment requests from other accounts.

To run, configure your Huaweicloud provider as described in the
[document](https://registry.terraform.io/providers/huaweicloud/huaweicloud/latest/docs).

## Usage

```
terraform init
terraform plan
terraform apply
terraform destroy
```

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.73.4 |
