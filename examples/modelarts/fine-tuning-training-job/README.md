# Create a fine-tuning training job

This example provides best practice code for using Terraform to create a ModelArts fine-tuning training job in
HuaweiCloud ModelArts service.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* A ModelArts dataset for fine-tuning
* An existing ModelArts asset
* An OBS path for saving the output model

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the ModelArts service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `training_job_name` - The name of the training job
* `training_job_asset_id` - The asset ID of the training job
* `training_job_asset_model` - The asset model configuration of the training job
  - `name` - The model name (required)
  - `version` - The model version (required)
  - `type` - The model type (required)
  - `code` - The model code (optional)
  - `desc` - The model description (optional)
  - `series` - The model series (optional)
* `training_job_ftjob_config` - The fine-tuning training job configuration
  - `envs` - The fine-tuning training environment variables (required).
    + `env_name` - The name of the environment variable (required)
    + `env_type` - The type of the environment variable (required)
    + `value` - The value of the environment variable (required)
    + `label` - The label of the environment variable (optional)
    + `des` - The description of the environment variable (optional)
    + `modifiable` - Whether the environment variable is modifiable (optional)
    + `displayable` - Whether the environment variable is displayable (optional)
  - `checkpoint_config` - The checkpoint configuration (optional)
    + `checkpoint_id` - The checkpoint ID (optional)
    + `save_checkpoints_max` - The maximum number of checkpoints to save (optional)
    + `skipped_steps` - The number of steps to skip (optional)
    + `restore_training` - Whether to restore training from a checkpoint (optional)

#### Optional Variables

* `workspace_name` - The name of the workspace to create (default: ""). Cannot be configured together with `workspace_id`
* `topic_name` - The name of the SMN topic to create for notifications (default: ""). Cannot be configured together
  with `training_job_notification_topic_urn`. If specified, a new SMN topic will be created
* `training_job_train_type` - The training type of the fine-tuning job (default: "SFT").
* `workspace_id` - The existing workspace ID (default: ""). Cannot be configured together with `workspace_name`

  -> If `workspace_id` and `workspace_name` are omitted, the default workspace is used.

* `training_job_annotations` - The annotations of the training job (default: {})
* `training_job_description` - The description of the training job (default: null)
* `training_job_inputs` - The inputs of the training job (default: [])
  - `dataset` - The dataset input configuration (required)
    + `id` - The dataset ID (required)
    + `name` - The dataset name (required)
    + `version_id` - The dataset version ID (optional)
    + `service_type` - The dataset service type (optional), e.g. **V3**
    + `dataset_proportion` - The dataset proportion (optional)
  - `local_dir` - The local directory mapped in the container (optional)
* `training_job_environments` - The environment variables of the training job (default: null)
* `resource_flavor_id` - The flavor ID of the public resource pool (default: "")
* `resource_node_count` - The number of resource replicas used by the training job (default: 1)
* `training_job_notification_topic_urn` - The existing SMN topic URN for training job notifications (default: "").
  Cannot be configured together with `topic_name`
* `training_job_notification_events` - The notification events of the training job (default: []).
  Required when `training_job_notification_topic_urn` or `topic_name` is configured
* `training_job_output_model` - The output model configuration of the training job (default: null)
  - `obs_path` - The OBS path where the output model is saved (required)
  - `local_path` - The local path where the output model is saved (optional)
* `training_job_tags` - The tags of the training job (default: null)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  training_job_name   = "tf_test_training_job"
  training_job_inputs = [
    {
      dataset = {
        id                 = "your-dataset-id"
        name               = "your-dataset-name"
        dataset_proportion = 100
      }
    }
  ]

  resource_flavor_id    = "your-public-resource-pool-flavor"
  training_job_asset_id = "your-asset-id"

  training_job_asset_model = {
    name    = "new-fine-tuning"
    version = "1.0.0"
    type    = "NEW_ASSET"
    desc    = "assert new fine-tuning model"
  }

  training_job_output_model = {
    obs_path = "your-obs-path-for-output-model"
  }

  training_job_ftjob_config = {
    envs = [
      {
        label       = "MIN_LR"
        des         = "Minimum learning rate"
        env_name    = "MIN_LR"
        env_type    = "string"
        value       = "1.25e-7"
        modifiable  = true
        displayable = true
      },
      {
        label       = "LR"
        des         = "Learning rate"
        env_name    = "LR"
        env_type    = "string"
        value       = "1.25e-6"
        modifiable  = true
        displayable = true
      },
    ]
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

## Notes

* Make sure to keep your credentials secure and never commit them to version control
* All resources will be created in the specified region
* `workspace_name` and `workspace_id` cannot be configured at the same time
* `topic_name` and `training_job_notification_topic_urn` cannot be configured at the same time
* `training_job_notification_events` is required when `topic_name` or `training_job_notification_topic_urn` is configured
* `training_job_ftjob_config.envs` must match the environment variables required by the selected training model

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.93.0 |
