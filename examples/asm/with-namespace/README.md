# Create an ASM Mesh with Namespace Configuration

This example provides best practice code for using Terraform to create an ASM (Application Service Mesh) mesh
with namespace configuration for sidecar injection. This example demonstrates two approaches: using an existing
namespace or creating a new namespace through Terraform.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* An existing CCE cluster in the target region
* At least one node in the CCE cluster for installing ASM mesh components
* An existing namespace in the CCE cluster (if using existing namespace approach)

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
  + Example: `dev-with-namespace-mesh`
* `mesh_version` - The version of the ASM mesh
* `cluster_id` - The ID of the CCE cluster to be associated with the ASM mesh
* `node_id` - The ID of the node where ASM mesh components will be installed

#### Conditionally Required Variables

* `namespaces` - The list of existing namespace names for sidecar injection. You must choose one of the following
* `namespace_name` - The name of the namespace to create. Required when `namespaces` is empty (default behavior)

  -> Configure either `namespaces` or `namespace_name`. When `namespaces` is provided, the CCE namespace resource is
  skipped and `namespace_name` is not used

#### Optional Variables

* `mesh_type` - The type of the ASM mesh. Currently, only **InCluster** is supported (default: **InCluster**)
* `tags` - The key/value pairs to associate with the ASM mesh (default: {})

## Architecture Overview

This example demonstrates an ASM mesh creation workflow with namespace configuration:

1. **Namespace configuration (choose one approach)**:

    + **Option 1: Use existing namespaces**
      * Provide a list of existing namespace names via `namespaces` variable
      * The ASM mesh will use these namespaces for sidecar injection
      * No new namespace will be created

    + **Option 2: Create new namespace (default)**
      * Leave `namespaces` variable empty (empty list)
      * Provide `namespace_name` for the new namespace
      * A new CCE namespace will be created using `huaweicloud_cce_namespace`
      * The ASM mesh will use this newly created namespace for sidecar injection

## Usage

* Copy this example to your working directory.

* Create an `authentication.auto.tfvars` file for credentials:

  ```hcl
  region_name = "your-region-name"
  access_key  = "your-access-key"
  secret_key  = "your-secret-key"
  ```

### Example 1: Use Existing Namespace

* Create a `terraform.tfvars` file and configure to use existing namespace:

    ```hcl
    # ASM mesh configuration
    mesh_name    = "test-with-namespace-mesh"
    mesh_version = "1.18.7-r7"
  
    # CCE cluster and node
    cluster_id = "your-cce-cluster-id"
    node_id    = "your-node-id"
  
    # Use existing namespaces
    namespaces = ["your-existing-namespace-1", "your-existing-namespace-2"]
  
    # Tags for the mesh
    tags = {
      foo = "bar"
      key = "value"
    }
  ```

### Example 2: Create New Namespace

* Create a `terraform.tfvars` file and configure to create new namespace:

    ```hcl
    # ASM mesh configuration
    mesh_name    = "test-with-namespace-mesh"
    mesh_version = "1.18.7-r7"
  
    # CCE cluster and node
    cluster_id = "your-cce-cluster-id"
    node_id    = "your-node-id"
  
    # Create a new namespace
    namespace_name = "asm-test-namespace"
  
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
* You must choose either `namespaces` (existing namespace) or creating a new namespace. These two approaches are
  mutually exclusive.
* When using an existing namespace, make sure the namespace exists in the CCE cluster before creating the ASM mesh.
* The ASM mesh components will be installed on the specified node using a field selector with the node UID.
* The sidecar injection will be configured for the specified namespace using a field selector with the namespace name.
* The mesh creation and deletion operations have a timeout of 30 minutes.
* The namespace creation and deletion operations have a timeout of 5 minutes.

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 1.1.0  |
| huaweicloud | >= 1.75.0 |
