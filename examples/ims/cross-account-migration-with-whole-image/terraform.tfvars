# Authentication variables
region_name = "cn-north-4"

# Variable definitions for sharer account
access_key = "Your_sharer_access_key"
secret_key = "Your_sharer_secret_key"

# Variable definitions for accepter account
accepter_access_key = "Your_accepter_access_key"
accepter_secret_key = "Your_accepter_secret_key"

# Resource variables
vpc_name               = "tf_test_whole_image_vpc"
subnet_name            = "tf_test_whole_image_subnet"
security_group_name    = "tf_test_whole_image_sg"
instance_name          = "tf_test_whole_image_ecs"
administrator_password = "YourPassword@12!"
instance_data_disks    = [
  {
    size = 10
    type = "SAS"
  }
]

vault_name             = "tf_test_sharer_vault"
whole_image_name       = "tf_test_whole_image"
accepter_project_ids   = ["your_accepter_project_id"]
accepter_vault_name    = "tf_test_accepter_vault"
accepter_instance_name = "tf_test_accepter_instance"
