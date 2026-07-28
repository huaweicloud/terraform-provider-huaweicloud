# Manage GaussDB Parameter Templates

This example provides best practice code for using Terraform to manage parameter templates
of a GaussDB OpenGauss instance in HuaweiCloud, covering template creation, application.

## Prerequisites

* A HuaweiCloud account with GaussDB permissions
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the GaussDB instance is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the VPC subnet
* `security_group_name` - The name of the security group
* `instance_name` - The name of the GaussDB instance
* `template_name` - The name of the parameter template

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: `192.168.0.0/16`)
* `subnet_cidr` - The CIDR block of the VPC subnet (default: `192.168.0.0/24`)
* `subnet_gateway_ip` - The gateway IP of the VPC subnet (default: `""`, auto-calculated from subnet CIDR)
* `security_group_rule_ports` - The security group ingress rule ports, separated by commas
  (default: `2379-2380,5000-5001,5432-5532,6000,6500,12016,20050`)
* `instance_availability_zones` - The availability zones of the GaussDB instance, separated by commas (default: `""`)
* `instance_flavor` - The flavor spec code of the GaussDB instance (default: `gaussdb.opengauss.ee.c3.xlarge.x864.ha`)
* `instance_password` - The password of the GaussDB instance (default: `""`)
* `instance_db_port` - The database port of the GaussDB instance (default: `5432`)
* `enterprise_project_id` - The enterprise project ID of the GaussDB instance (default: `null`)
* `instance_ha_mode` - The HA mode of the GaussDB instance (default: `"centralization_standard"`)
* `instance_ha_replication_mode` - The HA replication mode of the GaussDB instance (default: `"sync"`)
* `instance_ha_consistency` - The HA consistency of the GaussDB instance (default: `"eventual"`)
* `instance_volume_type` - The volume type of the GaussDB instance (default: `ULTRAHIGH`)
* `instance_volume_size` - The volume size (GB) of the GaussDB instance (default: `40`)
* `engine_version` - The engine version of the parameter template (default: `"8.218"`)
* `instance_mode` - The instance mode of the parameter template (default: `"ha"`)
* `template_description` - The description of the parameter template (default: `""`)
* `template_parameters` - The parameters of the parameter template (default: `[]`)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "terraform-gaussdb-vpc"
  subnet_name         = "terraform-gaussdb-subnet"
  security_group_name = "terraform-gaussdb-secgroup"
  instance_name       = "terraform-gaussdb-instance"
  template_name       = "terraform-parameter-template"
  ```

* Initialize Terraform:

  ```bash
  $ terraform init
  ```

* Review the Terraform plan:

  ```bash
  $ terraform plan
  ```

* Apply the configuration:

  ```bash
  $ terraform apply
  ```

* To clean up the resources:

  ```bash
  $ terraform destroy
  ```

## Note

* Make sure to keep your credentials secure and never commit them to version control
* `huaweicloud_gaussdb_parameter_template` is a ForceNew resource, any field change will recreate the template.
  Use `huaweicloud_gaussdb_parameter_template_save` to modify parameter values instead
* When creating a parameter template, you must specify either `engine_version` + `instance_mode` (for a new template)
  or `source_configuration_id` (to copy from an existing template). These two modes are mutually exclusive.
  In this example, `engine_version` and `instance_mode` are configured via variables
* `huaweicloud_gaussdb_parameter_template_apply` is a one-time action resource, re-running `terraform apply`
  will not re-apply the template. Delete this resource from state first if you need to re-apply
* Some parameters require an instance restart after modification. Check the `need_restart` attribute in the
  template parameters to determine this
* Deleting `huaweicloud_gaussdb_parameter_template_apply`, `compare`, or `reset` resources only removes them
  from the Terraform state, it does not undo the operation
* `engine_version` and `instance_mode` must match the GaussDB instance, otherwise the parameter template
  cannot be applied to the instance. The `instance_mode` must match the deployment model of the GaussDB instance

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 1.1.0 |
| huaweicloud | >= 1.90.0 |
