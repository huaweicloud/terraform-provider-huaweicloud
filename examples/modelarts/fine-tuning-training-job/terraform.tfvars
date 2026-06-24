training_job_name = "tf_test_training_job"

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

  checkpoint_config = {
    save_checkpoints_max = 5
  }
}
