# Create a volume snapshot

This example provides best practice code for using Terraform to create a volume sanpshot in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Required Variables

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the resources are located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `volume_name` - The volume name
* `snapshot_name` - The snapshot name

#### Optional Variables

* `volume_availability_zone` - The availability zone for the volume (default: "")
* `volume_type` - The type of the volume (default: "SAS")
* `volume_size` - The size of the volume (default: 20)
* `voluem_description` - The description of the volume (default: "")
* `vouleme_multiattach` - The volume is shared volume or not (default: false)
* `snapshot_description` - The description of the snapshot (default: "")
* `snapshot_metadata` - The metadata information of the snapshot (default: {})
* `snapshot_force` - The flag for forcibly creating a snapshot (default: false)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  volume_name   = "your_volume_name"
  snapshot_name = "your_snapshot_name"
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
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.62.0 |
