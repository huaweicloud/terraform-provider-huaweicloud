# Create notebook with dedicated resource pool

This example provides best practice code for using Terraform to create a ModelArts notebook on a dedicated
resource pool in HuaweiCloud ModelArts service. The notebook can optionally mount **OBS Parallel File System (PFS)** storage.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the ModelArts service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The VPC name
* `subnet_name` - The subnet name
* `security_group_name` - The name of the security group
* `turbo_name` - The name of the SFS Turbo
* `network_name` - The name of the ModelArts network
* `resource_pool_name` - The name of the dedicated resource pool
* `notebook_name` - The name of the notebook

#### Optional Variables

* `enterprise_project_id` - The enterprise project ID (default: null)
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: auto-calculated from VPC CIDR)
* `subnet_gateway_ip` - The gateway IP of the subnet (default: auto-calculated from subnet CIDR)
* `turbo_size` - The size of the SFS Turbo (default: 1228)
* `turbo_share_proto` - The share protocol of the SFS Turbo (default: "NFS")
* `turbo_share_type` - The share type of the SFS Turbo (default: "HPC")
* `turbo_hpc_bandwidth` - The HPC bandwidth of the SFS Turbo (default: "40M")
* `network_cidr` - The CIDR block of the ModelArts network (default: "10.168.0.0/16")
* `workspace_name` - The name of the workspace to create (default: ""). Cannot be configured together with `workspace_id`
* `resource_pool_flavor_id` - The flavor ID of the dedicated resource pool (default: "")
* `resource_pool_scope` - The scope of the dedicated resource pool (default: ["Notebook", "Train", "Infer"])
* `workspace_id` - The existing workspace ID (default: ""). Cannot be configured together with `workspace_name`

  -> If `workspace_id` and `workspace_name` are omitted, the default workspace is used.

* `notebook_flavor_id` - The flavor ID of the notebook (default: ""). If empty, the first available dedicated flavor
  is used
* `notebook_flavor_category` - The processor type of the notebook flavor (default: "CPU"). Valid values are **CPU**,
  **GPU**, and **ASCEND**
* `notebook_image_id` - The image ID of the notebook (default: ""). If empty, the first available image is used
* `notebook_image_type` - The type of the notebook image (default: "BUILD_IN"). Valid values are **BUILD_IN** and
  **DEDICATED**
* `notebook_key_pair_name` - The existing key pair name for remote SSH access (default: "")
* `allowed_access_ips` - The public IP addresses that are allowed for remote SSH access (default: [])
* `keypair_name` - The name of the KPS key pair to create automatically (default: "")
* `notebook_description` - The description of the notebook (default: null)
* `notebook_tags` - The tags of the notebook (default: {})
* `notebook_mount_storage_path` - The OBS path of PFS or its folders to mount (default: "")
* `notebook_mount_storage_local_directory` - The local mount directory for the OBS storage (default: "")

## Usage

* Copy the example files (`main.tf`, `variables.tf`, and `providers.tf`) to your working directory.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "tf_test_vpc"
  subnet_name         = "tf_test_subnet"
  security_group_name = "tf_test_security_group"
  turbo_name          = "tf_test_sfs_turbo"
  network_name        = "tf_test_network"
  resource_pool_name  = "tf_test_resource_pool"
  notebook_name       = "tf_test_notebook"
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

## Configuration Details

### If you want to enable SSH access with IP restriction

Use an existing key pair:

```hcl
notebook_key_pair_name = "your_key_pair_name"
allowed_access_ips     = ["your_public_ip_address"]
```

Or create a new key pair automatically:

```hcl
keypair_name       = "your_key_pair_name"
allowed_access_ips = ["your_public_ip_address"]
```

### If you want to dynamically mount OBS PFS storage

```hcl
notebook_mount_storage_path            = "obs://your-bucket/folder/"
notebook_mount_storage_local_directory = "/data/test/"
```

## Notes

* Make sure to keep your credentials secure and never commit them to version control
* All resources will be created in the specified region
* If `resource_pool_flavor_id`, `notebook_flavor_id`, or `notebook_image_id` is not specified, the example
  automatically selects the first available option
* When `allowed_access_ips` is configured, either `notebook_key_pair_name` or `keypair_name` must be specified
* `key_pair` and `allowed_access_ips` must be specified together
* `workspace_name` and `workspace_id` cannot be configured at the same time
* The dedicated notebook uses EFS storage backed by SFS Turbo connected through the ModelArts network

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.92.0 |
