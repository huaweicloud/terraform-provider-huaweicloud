# Create a cloud desktop machine

Configuration in this directory describes how to enable the Workspace service and create a cloud desktop machine.

Notes the cloud desktop machine and the user are dependent on the service configuration.  
Please read the implicit and explicit dependencies in the script carefully.

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
| huaweicloud | >= 1.67.0 |
