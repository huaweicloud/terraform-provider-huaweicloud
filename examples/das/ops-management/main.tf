# Create an instance group for batch instance management
resource "huaweicloud_das_instance_group" "test" {
  datastore_type = var.ops_datastore_type
  group_name     = var.ops_group_name
  description    = var.ops_group_description
}

# Assign instances to the instance group
resource "huaweicloud_das_instance_group_assign" "test" {
  group_id     = huaweicloud_das_instance_group.test.id
  instance_ids = var.ops_group_instance_ids
}

# Create an email template for inspection report notification
resource "huaweicloud_das_email_template" "test" {
  datastore_type  = var.ops_datastore_type
  name            = var.ops_email_template_name
  groups          = [huaweicloud_das_instance_group.test.id]
  health_rank     = var.ops_email_health_rank
  inspection_time = var.ops_email_inspection_time
  send_time       = var.ops_email_send_time
  time_zone       = var.ops_email_time_zone
  email           = var.ops_email_address
  topic           = var.ops_email_topic
  topic_urn       = var.ops_email_topic_urn
  obs_bucket_name = var.ops_email_obs_bucket_name
}

# Batch subscribe to email templates
resource "huaweicloud_das_email_templates_batch_action" "test" {
  subscribe          = var.ops_email_subscribe
  email_template_ids = [huaweicloud_das_email_template.test.id]
}
