# Create EVS v3 volume and snapshot

This example provides best practice code for using Terraform to create an EVS v3 volume and v3 snapshot in HuaweiCloud.

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

* `volume_image_id` - The ID of the image used to create the volume (default: "")
* `volume_image_visibility` - The visibility of the volume image (default: "public")
* `volume_image_os` - The OS of the volume image (default: "Ubuntu")
* `volume_type` - The type of the volume, the value can be **SATA**, **SAS**, **GPSSD**, **SSD**, **ESSD**, **GPSSD2**
  or **ESSD2** (default: "GPSSD")
* `volume_availability_zone` - The availability zone for the volume (default: "")
* `volume_description` - The description of the volume (default: "")
* `volume_metadata` - The metadata of the volume (default: {})
* `volume_multiattach` - The volume is shared volume or not (default: false)
* `volume_size` - The size of the volume (default: 40)
* `volume_tags` - The tags of the volume (default: {})
* `snapshot_metadata` - The metadata information of the snapshot (default: {})
* `snapshot_description` - The description of the snapshot (default: "")

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

| Name | Version   |
| ---- |-----------|
| terraform | >= 1.1.0  |
| huaweicloud | >= 1.77.5 |
