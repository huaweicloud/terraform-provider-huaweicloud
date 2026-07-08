# Create a Basic ASM Mesh for Single Cluster

This example provides best practice code for using Terraform to create an ASM (Application Service Mesh) mesh for
a single CCE cluster in HuaweiCloud. This is a basic scenario suitable for development and testing environments.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* An existing CCE cluster in the target region
* At least one node in the CCE cluster for installing ASM mesh components

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

Configure authentication variables in `authentication.auto.tfvars` (recommended) or `terraform.tfvars`:

* `region_name` - The region where the ASM mesh and CCE cluster are located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `mesh_name` - The name of the ASM mesh. Requirements:
  + `4` to `64` characters
  + Can include letters, digits, and hyphens (-)
  + Must start with a letter
  + Cannot end with a hyphen (-)
  + Example: `dev-basic-mesh`
* `mesh_version` - The version of the ASM mesh
* `cluster_id` - The ID of the CCE cluster to be associated with the ASM mesh
* `node_id` - The ID of the node where ASM mesh components will be installed

#### Optional Variables

* `mesh_type` - The type of the ASM mesh. Currently, only "InCluster" is supported (default: "InCluster")
* `tags` - The key/value pairs to associate with the ASM mesh (default: {})

## Architecture Overview

This example follows a simple ASM mesh creation workflow:

1. **Create ASM mesh**:

  * Create an ASM mesh with basic configuration using `huaweicloud_asm_mesh`
  * Configure mesh name and version
  * Associate the mesh with a single CCE cluster
  * Specify the node for installing mesh components using field selector
  * Apply tags for resource categorization

## Usage

* Copy this example to your working directory.

* Create an `authentication.auto.tfvars` file for credentials:

  ```hcl
  region_name = "your-region-name"
  access_key  = "your-access-key"
  secret_key  = "your-secret-key"
  ```

* Create a `terraform.tfvars` file and fill in the variables:

  ```hcl
  # ASM mesh configuration
  mesh_name    = "test-basic-mesh"
  mesh_version = "1.18.7-r7"

  # CCE cluster and node
  cluster_id = "your-cce-cluster-id"
  node_id    = "your-node-id"

  # Tags for the mesh
  tags = {
    foo = "bar"
    key = "value"
  }
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
* The mesh name must be unique in the region. If a mesh with the same name already exists, the deployment will fail.
* The ASM mesh type, version, name, and cluster configuration are non-updatable. Changing them will recreate the mesh
  resource.
* The ASM mesh components will be installed on the specified node using a field selector with the node UID.
* The mesh creation and deletion operations have a timeout of `30` minutes.

## Example Outputs

After successful deployment, you can query the mesh status:

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 1.1.0  |
| huaweicloud | >= 1.75.0 |
