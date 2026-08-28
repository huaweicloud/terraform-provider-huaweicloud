# GaussDB download backups authorize

This example provides best practice code for using Terraform to authorize downloading GaussDB backups within
HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* An existing GaussDB backup ID (UUID format, 36 characters)

## Required Variables

### Authentication Variables

* `region_name` - The region where the GaussDB download backups authorize resource is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `backup_id` - The backup ID in UUID format (36 characters) that uniquely identifies a GaussDB backup

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  backup_id = "your_gaussdb_backup_id"
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

## Note

* This resource is a one-time action resource that authorizes users to download GaussDB backups
* The resource uses the GaussDB (opengauss) service endpoint
* The create operation sends a POST request to `/v3/{project_id}/backups/{backup_id}/download/authorization`
* The resource ID equals the `backup_id`
* The `backup_id` parameter is non-updatable and will trigger resource replacement if changed (FlexibleForceNew)
* The read operation is a no-op (returns nil), no API call is made during read
* The update operation is a no-op (returns nil)
* Deleting this resource only removes it from Terraform state with a warning, the authorization is not revoked in the
  cloud
* The computed attributes `bucket` and `file_paths` are populated from the create API response
* The `bucket` attribute indicates the OBS bucket name where the backup file is stored
* The `file_paths` attribute indicates the paths from which backups can be downloaded using OBS Browser+
* This resource does not support import

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 0.14.0 |
| huaweicloud | >= 1.95.0 |
