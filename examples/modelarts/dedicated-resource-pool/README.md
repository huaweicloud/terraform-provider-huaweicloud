# Create a dedicated resource pool

This example provides best practice code for using Terraform to create a ModelArts dedicated resource pool in
HuaweiCloud ModelArts service. You can either provide an existing `network_id` or create a ModelArts network. SFS Turbo
connections are optional when creating a ModelArts network.

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

* `resource_pool_name` - The name of the dedicated resource pool
* `resource_pool_resources` - The list of resource specifications in the resource pool
  - `flavor_id` - The resource flavor ID (optional). If empty, the first available flavor in the current AZ is used
  - `count` - The number of resources of the corresponding flavors (required)
  - `max_count` - The max number of resources of the corresponding flavors (optional)
  - `extend_params` - The extend parameters of the resource pool, in JSON format (optional)
  - `root_volume` - The root volume of the resource pool nodes (optional)
    * `volume_type` - The type of the root volume
    * `size` - The size of the root volume
  - `data_volumes` - The data volumes of the resource pool nodes (optional)
    * `volume_type` - The type of the data volume
    * `size` - The size of the data volume
    * `extend_params` - The extend parameters of the data volume, in JSON format (optional)
    * `count` - The count of the current data volume configuration (optional)
  - `volume_group_configs` - The extend configurations of the volume groups (optional)
    * `volume_group` - The name of the volume group
    * `docker_thin_pool` - The percentage of container volumes to data volumes (optional)
    * `types` - The storage types of the volume group
    * `lvm_config` - The LVM management configuration (optional)
      + `lv_type` - The LVM write mode
      + `path` - The volume mount path (optional)
  - `os` - The image information for the specified OS (optional)
    * `name` - The OS name of the image (optional)
    * `image_id` - The image ID (optional)
    * `image_type` - The image type (optional)
  - `driver` - The driver information (optional)
    * `version` - The driver version
  - `creating_step` - The creation step configuration of the resource pool nodes (optional)
    * `step` - The creation step of the resource pool nodes
    * `type` - The type of the resource pool nodes

-> Either specify `network_id`, or specify `network_name` to create a ModelArts network.

#### Optional Variables

* `turbo_name` - The name of the SFS Turbo to create (default: ""). Can only be specified when `network_name` is
  specified. If specified, VPC, subnet, security group, and SFS Turbo will be created and connected to the ModelArts
  network
* `vpc_name` - The VPC name (default: ""). Can only be specified when `turbo_name` is specified
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `enterprise_project_id` - The enterprise project ID (default: null)
* `subnet_name` - The subnet name (default: ""). Can only be specified when `turbo_name` is specified
* `subnet_cidr` - The CIDR block of the subnet (default: auto-calculated from VPC CIDR)
* `subnet_gateway_ip` - The gateway IP of the subnet (default: auto-calculated from subnet CIDR)
* `security_group_name` - The name of the security group (default: ""). Can only be specified when `turbo_name` is
  specified
* `turbo_size` - The size of the SFS Turbo (default: 1228)
* `turbo_share_proto` - The share protocol of the SFS Turbo (default: "NFS")
* `turbo_share_type` - The share type of the SFS Turbo (default: "HPC")
* `turbo_hpc_bandwidth` - The HPC bandwidth of the SFS Turbo (default: "40M")
* `network_sfs_turbos` - The existing SFS Turbo connections for the ModelArts network (default: []). Can only be
  specified when `network_name` is specified and `turbo_name` is not specified
  - `id` - The SFS Turbo ID (required)
  - `name` - The SFS Turbo name (required)
* `network_name` - The name of the ModelArts network to create (default: ""). Cannot be configured together with
  `network_id`
* `network_cidr` - The CIDR block of the ModelArts network (default: "10.168.0.0/16")
* `workspace_name` - The name of the workspace to create (default: ""). Cannot be configured together with `workspace_id`
* `resource_pool_description` - The description of the dedicated resource pool (default: null)
* `resource_pool_scope` - The scope of the dedicated resource pool (default: ["Train", "Infer", "Notebook"])
* `network_id` - The existing ModelArts network ID (default: ""). Cannot be configured together with `network_name`
* `workspace_id` - The existing workspace ID (default: ""). Cannot be configured together with `workspace_name`

  -> If `workspace_id` and `workspace_name` are omitted, the default workspace is used.

* `resource_pool_metadata_annotations` - The annotations of the resource pool, in JSON format (default: null)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  network_id              = "your_modelarts_network_id"
  resource_pool_name      = "tf-test-resource-pool"
  resource_pool_resources = [
    {
      count = 1
    }
  ]
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

### If you want to use an existing ModelArts network, you need configure the specified network ID

```hcl
network_id              = "your_modelarts_network_id"
resource_pool_name      = "tf-test-resource-pool"
resource_pool_resources = [
  {
    count = 1
  }
]
```

### If you want to create a network without SFS Turbo, you need configure the specified network name

```hcl
network_name            = "tf-test-network"
resource_pool_name      = "tf-test-resource-pool"
resource_pool_resources = [
  {
    count = 1
  }
]
```

### If you want to associate existing SFS Turbos when creating a new network, you need to specify network_sfs_turbos

```hcl
network_name       = "tf-test-network"
resource_pool_name = "tf-test-resource-pool"
network_sfs_turbos = [
  {
    id   = "your_sfs_turbo_id"
    name = "your_sfs_turbo_name"
  }
]

resource_pool_resources = [
  {
    count = 1
  }
]
```

### If you want to create a new network with a new SFS Turbo, you need to specify network name and turbo name

```hcl
vpc_name                = "tf_test_vpc"
subnet_name             = "tf_test_subnet"
security_group_name     = "tf_test_security_group"
turbo_name              = "tf_test_sfs_turbo"
network_name            = "tf-test-network"
resource_pool_name      = "tf-test-resource-pool"
resource_pool_resources = [
  {
    count = 1
  }
]
```

## Notes

* Make sure to keep your credentials secure and never commit them to version control
* All resources will be created in the specified region
* Exactly one of `network_id` or `network_name` must be specified
* If `flavor_id` in `resource_pool_resources` is not specified, the example automatically selects the first available
  flavor in the current AZ
* `workspace_name` and `workspace_id` cannot be configured at the same time

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.92.0 |
