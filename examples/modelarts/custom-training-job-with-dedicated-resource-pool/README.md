# Create custom training job with dedicated resource pool

This example provides best practice code for using Terraform to create a ModelArts training job on a dedicated
resource pool in HuaweiCloud ModelArts service. The training job supports both **custom image** and **preset image**
configurations.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* Training code uploaded to OBS

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
* `training_job_name` - The name of the training job
* `training_job_code_dir` - The OBS code directory of the training job
* `training_job_command` - The container startup command for the custom image training job
* `training_job_engine` - The engine configuration of the training job
  - `image_url` - The custom image URL obtained from SWR (optional, required for custom image)
  - `id` - The engine specification ID (optional, required for preset image)
  - `version` - The engine specification version (optional, required for preset image)
  - `name` - The engine specification name (optional, required for preset image)

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
* `resource_pool_flavor_id` - The flavor ID of the dedicated resource pool (default: ""). If empty, the first available
  flavor in the current AZ is used
* `resource_pool_scope` - The scope of the dedicated resource pool (default: ["Notebook", "Train", "Infer"])
* `topic_name` - The name of the SMN topic to create for notifications (default: ""). Cannot be configured together
  with `training_job_notification_topic_urn`, if specified, a new SMN topic will be created
* `workspace_id` - The existing workspace ID (default: ""). Cannot be configured together with `workspace_name`

  -> If `workspace_id` and `workspace_name` are omitted, the default workspace is used.

* `training_job_annotations` - The annotations of the training job (default: {})
* `training_job_description` - The description of the training job (default: null)
* `training_job_inputs` - The inputs of the training job (default: [])
  - `local_dir` - The local directory mapped in the container (required)
  - `dataset` - The dataset input configuration (required)
    * `id` - The dataset ID (required)
    * `name` - The dataset name (optional)
    * `version_id` - The dataset version ID (optional)
    * `service_type` - The dataset service type (optional), e.g. **V3**
* `training_job_environments` - The environment variables of the training job (default: {})
* `resource_node_count` - The number of resource replicas used by the training job (default: 1)
* `training_job_volumes` - The volume mount configuration of the training job (default: [])
  - `nfs` - The NFS volume mount configuration (optional)
    * `nfs_server_path` - The NFS server path (optional)
    * `local_path` - The local mount path in the training container (optional)
    * `read_only` - Whether the NFS volume is read-only in the container (optional, default: false)
  - `pfs` - The PFS volume mount configuration (optional)
    * `pfs_path` - The OBSFS path (optional)
    * `local_path` - The local mount path in the training container (optional)
  - `obs` - The OBS volume mount configuration (optional)
    * `obs_path` - The OBS path to be mounted (optional)
    * `local_path` - The local mount path in the training container (optional)
* `training_job_log_export_path_obs_url` - The OBS log export path of the training job (default: "")
* `training_job_log_export_config_version` - The log export config version of the training job (default: "")
* `training_job_auto_stop_duration` - The auto stop duration of the training job in hours (default: 0)
* `training_job_notification_topic_urn` - The existing SMN topic URN for training job notifications (default: "").
  Cannot be configured together with `topic_name`
* `training_job_notification_events` - The notification events of the training job (default: [])
* `custom_metrics` - The custom metrics collection configuration of the training job (default: null)
  - `exec` - The command-based metrics collection configuration (optional)
    * `command` - The command used to collect metrics (required when `exec` is specified)
  - `http_get` - The HTTP-based metrics collection configuration (optional)
    * `path` - The HTTP path used to collect metrics (required when `http_get` is specified)
    * `port` - The HTTP port used to collect metrics (required when `http_get` is specified)
* `training_job_asset_model` - The asset model configuration of the training job (default: null)
  - `name` - The model name (required)
  - `version` - The model version (required)
  - `type` - The model type (required)
  - `code` - The model code (optional)
  - `desc` - The model description (optional)
  - `series` - The model series (optional)
* `training_job_output_model` - The output model configuration of the training job (default: null)
  - `obs_path` - The OBS path where the output model is saved (required)
  - `local_path` - The local path where the output model is saved (optional)
* `training_job_tags` - The tags of the training job (default: {})

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name                = "tf_test_vpc"
  subnet_name             = "tf_test_subnet"
  security_group_name     = "tf_test_security_group"
  turbo_name              = "tf_test_sfs_turbo"
  network_name            = "tf_test_network"
  resource_pool_name      = "tf_test_resource_pool"
  training_job_name       = "tf_test_training_job"
  training_job_code_dir   = "your_training_code_dir"
  training_job_command    = "your_training_command"
  resource_pool_flavor_id = "your_resource_pool_flavor_id"

  training_job_engine = {
    image_url = "your_swr_image_url"
  }
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

### If you want to configure a custom image job, you need configure the specified image URL

```hcl
training_job_engine = {
  image_url = "your_swr_image_url"
}
```

### If you want to configure a preset image job, you need configure the specified engine

```hcl
training_job_engine = {
  id      = "your_engine_id"
  version = "your_engine_version"
  name    = "your_engine_name"
}
```

## Notes

* Make sure to keep your credentials secure and never commit them to version control
* All resources will be created in the specified region
* `image_url` must be a valid SWR image path in the same region as ModelArts
* If `resource_pool_flavor_id` is not specified, the example automatically selects the first available flavor in the
  current AZ
* `workspace_name` and `workspace_id` cannot be configured at the same time
* `topic_name` and `training_job_notification_topic_urn` cannot be configured at the same time

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.93.0 |
