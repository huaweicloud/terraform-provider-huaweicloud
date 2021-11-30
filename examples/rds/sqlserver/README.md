# Create the Microsoft SQL Server databases

The smallest management unit of RDS is the DB instance.
The configuration in this directory creates a single-type and an HA-type relational database of Microsoft SQL Server engine.

To run, configure your Huaweicloud provider as described in the
[document](https://registry.terraform.io/providers/huaweicloud/huaweicloud/latest/docs).

## Configuration

The data source [`huaweicloud_rds_engine_versions`](https://registry.terraform.io/providers/huaweicloud/huaweicloud/latest/docs/data-sources/rds_engine_versions)
provides a way to query the list of available engine versions of SQL Server database.

These examples assumes that you have created a random password. For more details about random password resource, please
refer to the corresponding provider
[document](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/password).

### Single Type

All editions supports single database creation:

* Standard edition
* Enterprise edition
* Web edition

### HA Type

Only standard editions support HA database creation.
If the HA-type is used, [ha_replication_mode](https://registry.terraform.io/providers/huaweicloud/huaweicloud/latest/docs/resources/rds_instance#ha_replication_mode)
is required.

## Usage

```shell
terraform init
terraform plan
terraform apply
terraform destroy
```

During the creation process, you may notice that the database instance on the console has been displayed as "running",
but the terraform is still running at this time. This is because the database needs to be backed up after the database
is created.

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.31.0 |
