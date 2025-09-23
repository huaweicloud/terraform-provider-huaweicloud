# Create an EVS volume

This example provides best practice code for using Terraform to create an EVS volume in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variables Introdution

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the volume is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `volume_name` - The name of the volume

#### Optional Variables

* `volume_availability_zone` - The availability zone for the volume (default: "")
* `volume_type` - The type of the volume, the value can be **SAS**, **SSD**, **GPSSD**, **ESSD**, **GPSSD2** or
  **ESSD2** (default: "SSD")
* `voulme_size` - The size of the volume (default: 40)
* `volume_description` - The description of the volume (default: "")
* `volume_multiattach` - The volume is shared volume or not (default: false)
* `volume_iops` - The IOPS for the volume, only valid and required when `volume_type` is set to
  **GPSSD2** or **ESSD2** (default: null)
* `volume_throughput` - The throughput for the volume, only valid and required when `volume_type` is set to
  **GPSSD2** (default: null)
* `volume_device_type` - The device type of disk to create (default: "VBD")
* `enterprise_project_id` - The enterprise project ID of the volume (default: null)
* `volume_tags` - The tags of the volume (default: {})
* `charging_mode` - The charging mode of the volume (default: "postPaid")
* `period_unit` - The period unit of the volume, only required when `charging_mode` is **prePaid** (default: null)
* `period` - The period of the volume, only required when `charging_mode` is **prePaid** (default: null)
* `auto_renew` - The auto renew of the volume (default: "false")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  volume_name = "your_volume_name"
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
| huaweicloud | >= 1.53.0 |
