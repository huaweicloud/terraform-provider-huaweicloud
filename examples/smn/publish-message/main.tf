resource "huaweicloud_smn_topic" "test" {
  name                  = var.topic_name
  display_name          = var.topic_display_name
  enterprise_project_id = var.enterprise_project_id
}

resource "huaweicloud_smn_subscription" "test" {
  topic_urn = huaweicloud_smn_topic.test.id
  protocol  = var.subscription_protocol
  endpoint  = var.subscription_endpoint
  remark    = var.subscription_description
}

resource "huaweicloud_smn_message_template" "test" {
  count    = var.template_name != "" ? 1 : 0

  name     = var.template_name
  protocol = var.template_protocol
  content  = var.template_content
}

resource "huaweicloud_smn_message_publish" "test" {
  topic_urn             = huaweicloud_smn_topic.test.topic_urn
  subject               = var.pulblish_subject
  message               = var.pulblish_message != "" ? var.pulblish_message : null
  message_structure     = var.pulblish_message_structure != "" ? var.pulblish_message_structure : null
  message_template_name = var.template_name != "" ? try(huaweicloud_smn_message_template.test[0].name, null) : null
  time_to_live          = var.pulblish_time_to_live
  tags                  = var.pulblish_tags

  dynamic "message_attributes" {
    for_each = var.pulblish_message_attributes

    content {
      name   = message_attributes.value.name
      type   = message_attributes.value.type
      value  = message_attributes.value.value
      values = message_attributes.value.values
    }
  }
}
