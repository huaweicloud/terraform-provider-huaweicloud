# Execute a DataArts Factory script

This example provides best practice code for using Terraform to create a DataArts Factory script and execute it once
in HuaweiCloud DataArts Studio service.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* An existing DataArts Studio workspace with a DLI data connection configured
* A DLI SQL queue (the `default` queue is used by default)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

Configure authentication variables in `authentication.auto.tfvars` (recommended) or `terraform.tfvars`:

* `region_name` - The region where the DataArts Studio instance is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `dli_database_name` - The name of the DLI database created for the script
* `dli_table_name` - The name of the DLI table created for the script
* `script_name` - The name of the DataArts Factory script

#### Conditionally Required Variables

* `workspace_id` - The ID of the workspace. Required when you already know the target workspace ID
* `instance_id` - The ID of the DataArts Studio instance. Required when `workspace_id` is omitted, used together
  with the `huaweicloud_dataarts_studio_workspaces` data source to resolve the workspace

  -> Configure either `workspace_id` or `instance_id`. When `workspace_id` is provided, the workspace data source is
  skipped and `instance_id` is not used.

#### Optional Variables

* `workspace_name` - The name of the workspace used to filter results (default: ""). Only effective when `workspace_id`
  is omitted
* `connection_name` - The name of the DLI data connection (default: ""). If omitted, the first DLI connection in the
  workspace is resolved via the `huaweicloud_dataarts_studio_data_connections` data source
* `queue_name` - The DLI queue name associated with the script (default: "default")
* `dli_database_description` - The description of the DLI database (default: "")
* `dli_table_description` - The description of the DLI table (default: "")
* `dli_table_columns` - The column definitions of the DLI table (default: name/string and age/int)
* `script_type` - The type of the DataArts Factory script (default: "DLISQL")
* `script_directory` - The directory path where the script is stored in DataArts Factory (default: "/terraform")
* `script_description` - The description of the DataArts Factory script (default: "")
* `script_content` - The SQL content of the DataArts Factory script (default: "", auto-generated SELECT statement if
  empty)
* `script_configuration` - The user-defined configuration parameters of the DataArts Factory script (default: {})
* `script_execute_params` - The execution parameters passed to the script content (default: Spark SQL adaptive
  configuration)

## Architecture Overview

This example follows a script execution workflow with conditional data source lookups:

1. **Resolve workspace** (optional):
   + When `workspace_id` is omitted, query `huaweicloud_dataarts_studio_workspaces` by `instance_id` (and optionally
     `workspace_name`) to obtain the target workspace ID
2. **Resolve DLI connection** (optional):
   + When `connection_name` is omitted, query `huaweicloud_dataarts_studio_data_connections` with type `DLI` to obtain
     the first available connection name
3. **Prepare DLI data sources** by creating a database and table
4. **Create a DataArts Factory script** that references the DLI database and table
5. **Execute the script once** using `huaweicloud_dataarts_factory_script_execute`

## Usage

* Copy this example to your working directory.

* Create an `authentication.auto.tfvars` file for credentials:

  ```hcl
  region_name = "your-region-name"
  access_key  = "your-access-key"
  secret_key  = "your-secret-key"
  ```

* Create a `terraform.tfvars` file and fill in the required variables.

  Example A — provide `workspace_id` directly:

  ```hcl
  workspace_id      = "your-dataarts-studio-workspace-id"
  dli_database_name = "tf_test_database"
  dli_table_name    = "tf_test_table"
  script_name       = "tf_test_factory_script"
  ```

  Example B — resolve workspace from instance:

  ```hcl
  instance_id       = "your-dataarts-studio-instance-id"
  dli_database_name = "tf_test_database"
  dli_table_name    = "tf_test_table"
  script_name       = "tf_test_factory_script"
  ```

* Initialize Terraform:

  ```bash
  terraform init
  ```

* Review the Terraform plan:

  ```bash
  terraform plan
  ```

* Apply the configuration:

  ```bash
  terraform apply
  ```

* To clean up the resources:

  ```bash
  terraform destroy
  ```

## Notes

* Make sure to keep your credentials secure and never commit them to version control.
* All resources will be created in the specified region.
* `huaweicloud_dataarts_factory_script_execute` is a one-time action resource. Deleting this resource will not remove
  the execution record from DataArts Studio, but will only remove the resource information from the tfstate file.
* Re-running `terraform apply` after a successful execution will trigger a new script execution because `workspace_id`
  and `script_name` are non-updatable parameters on the execute resource.
* Ensure the workspace has a valid DLI data connection before applying this example.
* When `workspace_id` is omitted, `instance_id` must be provided or Terraform will fail the workspace data source
  precondition check.
* After apply, verify the script execution status in the DataArts Studio console or by running
  `terraform state show huaweicloud_dataarts_factory_script_execute.test`.

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.3.0 |
| huaweicloud | >= 1.91.0 |
