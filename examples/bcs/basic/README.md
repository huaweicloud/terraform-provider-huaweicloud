# Create a BCS instance with basic configuration

This example provides best practice code for using Terraform to create a Blockchain Service (BCS) instance
with basic configuration in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* An existing CCE cluster (the BCS service needs to exclusively occupy the CCE cluster)
* An enterprise project ID

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

Configure authentication variables in `authentication.auto.tfvars` (recommended) or `terraform.tfvars`:

* `region_name` - The region where the BCS instance is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `instance_name` - The unique name of the BCS instance. The name consists of `4` to `24` characters,
  including letters, digits, chinese characters and hyphens (-), and the name cannot start with a hyphen
* `edition` - The service edition of the BCS instance. Valid values:
  + **1**
  + **2**
  + **4**
* `fabric_version` - The version of fabric for the BCS instance. Valid values:
  + **1.4**
  + **2.0**
* `consensus` - The consensus algorithm used by the BCS instance
  + The valid values of fabric 1.4 are **solo**, **kafka** and **SFLIC**
  + The valid values of fabric 2.0 are **SFLIC** and **etcdraft**
* `cce_cluster_id` - The ID of the CCE cluster to attach to the BCS instance. The BCS service needs to
  exclusively occupy the CCE cluster
* `enterprise_project_id` - The ID of the enterprise project that the BCS instance belongs to
* `instance_password` - The resource access and blockchain management password. The password consists
  of `8` to `12` characters and must consist at least three of following: uppercase letters, lowercase letters,
  digits, chinese characters, special characters(!@$%^-_=+[{}]:,./?)

#### Optional Variables

* `orderer_node_num` - The number of peers in the orderer organization (default: 1)
* `volume_type` - The storage volume type to attach to each organization of the BCS instance. Valid values:
  + **nfs**: SFS (default)
  + **efs**: SFS Turbo
* `org_disk_size` - The storage capacity of peer organization. The minimum storage capacity of `efs`
  `volume_type` is 500GB. The specifications are as follows when `volume_type` is `nfs`:
  + The minimum storage capacity of basic edition is `40` GB
  + The minimum storage capacity of enterprise and professional edition is `100` GB (default: 100)
* `block_info` - The configuration of block generation
  + `generation_interval` - The block generation time, the unit is second.(default: 2)
  + `transaction_quantity` - The number of transactions included in the block. (default: 500)
  + `block_size` - The volume of the block, the unit is MB. (default: 2)
* `peer_orgs` - The array of one or more peer organizations to attach to the BCS instance
  + `org_name`: The name of the peer organization (default: "organization")
  + `count`: The number of peers in organization (default: 2)
* `channels` - The array of one or more channels to attach to the BCS instance
  + `name`: The name of the channel (default: "channel")
  + `org_names`: The name of the peer organization (default: "organization")

## Architecture Overview

This example creates a BCS instance with basic configuration:

1. **Create BCS instance**:

  * Create a BCS instance with basic configuration using `huaweicloud_bcs_instance`
  * Configure instance name, edition, fabric version, and consensus algorithm
  * Attach the BCS instance to an existing CCE cluster
  * Configure block information and channels

## Usage

* Copy this example to your working directory.

* Create an `authentication.auto.tfvars` file for credentials:

  ```hcl
  region_name = "your-region-name"
  access_key  = "your-access-key"
  secret_key  = "your-secret-key"
  ```

* Create a `terraform.tfvars` file and fill in the required variables.

  Example:

  ```hcl
  instance_name         = "tf_test_bcs_instance"
  edition               = 1
  fabric_version        = "2.0"
  consensus             = "etcdraft"
  orderer_node_num      = 1
  cce_cluster_id        = "your-cce-cluster-id"
  enterprise_project_id = "your-enterprise-project-id"
  instance_password     = "your-instance-password"
  volume_type           = "nfs"
  org_disk_size         = 100

  block_info = [
    {
      generation_interval  = 2
      transaction_quantity = 500
      block_size           = 2
    }
  ]

  peer_orgs = [
    {
      org_name = "organization01"
      count    = 2
    }
  ]

  channels = [
    {
      name      = "channel01"
      org_names = ["organization01"]
    }
  ]
  ```

  **Note**: When using `volume_type = "efs"` with `edition = 4`, you must configure `sfs_turbo`:

  ```hcl
  sfs_turbo = [
    {
      share_type        = "STANDARD"
      type              = "efs-ha"
      availability_zone = ""
      flavor            = "sfs.turbo.20MBps"
    }
  ]
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
* The BCS service needs to exclusively occupy the CCE cluster. Please make sure that the CCE cluster
  is not occupied before deploying the BCS service.
* Changing parameters with `ForceNew` will create a new instance and delete the old one.
* This example uses only required fields for a basic configuration. For advanced configurations,
  refer to the `huaweicloud_bcs_instance` documentation.

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 1.9.0  |
| huaweicloud | >= 1.80.5 |
