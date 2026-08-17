# Variable definitions for authentication
variable "region_name" {
  description = "The region where the DAS resources are located"
  type        = string
}

variable "access_key" {
  description = "The access key of the IAM user"
  type        = string
  sensitive   = true
}

variable "secret_key" {
  description = "The secret key of the IAM user"
  type        = string
  sensitive   = true
}

# Variable definitions for huaweicloud_das_instance_group
variable "ops_datastore_type" {
  description = "The database type"
  type        = string
}

variable "ops_group_name" {
  description = "The instance group name"
  type        = string
}

variable "ops_group_description" {
  description = "The description of the instance group"
  type        = string
}

# Variable definitions for huaweicloud_das_instance_group_assign
variable "ops_group_instance_ids" {
  description = "The list of instance IDs to be assigned to the group"
  type        = list(string)
}

# Variable definitions for huaweicloud_das_email_template
variable "ops_email_template_name" {
  description = "The name of the email template"
  type        = string
}

variable "ops_email_health_rank" {
  description = "The list of health ranks"
  type        = list(string)
}

variable "ops_email_inspection_time" {
  description = "The diagnosis time"
  type        = string
}

variable "ops_email_send_time" {
  description = "The send time"
  type        = string
}

variable "ops_email_time_zone" {
  description = "The time zone"
  type        = string
}

variable "ops_email_address" {
  description = "The email address for notification"
  type        = string
  default     = null
}

variable "ops_email_topic" {
  description = "The topic ID for notification"
  type        = string
  default     = null
}

variable "ops_email_topic_urn" {
  description = "The topic URN for notification"
  type        = string
  default     = null
}

variable "ops_email_obs_bucket_name" {
  description = "The OBS bucket name for storing inspection reports"
  type        = string
  default     = null
}

# Variable definitions for huaweicloud_das_email_templates_batch_action
variable "ops_email_subscribe" {
  description = "Whether to subscribe to the email templates"
  type        = bool
}
