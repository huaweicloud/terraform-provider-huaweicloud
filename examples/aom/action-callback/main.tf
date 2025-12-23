resource "huaweicloud_smn_topic" "test" {
  name                  = var.smn_topic_name
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}

resource "huaweicloud_smn_subscription" "test" {
  count = length(var.alarm_callback_urls) > 0 ? 1 : 0

  topic_urn = huaweicloud_smn_topic.test.id
  protocol  = length(regexall("^https?://", var.alarm_callback_urls[count.index])) > 0 ? "https" : "http"
  endpoint  = var.alarm_callback_urls[count.index]
}

resource "huaweicloud_aom_message_template" "test" {
  name                  = var.alarm_notification_template_name
  locale                = var.alarm_notification_template_locale
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
  description           = var.alarm_notification_template_description

  templates {
    sub_type = var.alarm_notification_template_notification_type
    topic    = var.alarm_notification_template_notification_topic
    content  = var.alarm_notification_template_content
  }
}

resource "huaweicloud_aom_alarm_action_rule" "test" {
  depends_on = [huaweicloud_aom_message_template.test]

  name                  = var.alarm_action_rule_name
  user_name             = var.alarm_action_rule_user_name
  type                  = var.alarm_action_rule_type
  notification_template = huaweicloud_aom_message_template.test.name

  smn_topics {
    topic_urn = huaweicloud_smn_topic.test.topic_urn
  }
}
