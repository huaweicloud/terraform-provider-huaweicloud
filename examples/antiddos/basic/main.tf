resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = var.vpc_eip_publicip_type
  }

  bandwidth {
    share_type  = var.vpc_eip_bandwidth_share_type
    name        = var.vpc_eip_bandwidth_name
    size        = var.vpc_eip_bandwidth_size
    charge_mode = var.vpc_eip_bandwidth_charge_mode
  }
}

resource "huaweicloud_smn_topic" "test" {
  name         = var.smn_topic_name
  display_name = var.smn_topic_display_name
}

resource "huaweicloud_smn_subscription" "test" {
  topic_urn = huaweicloud_smn_topic.test.id
  endpoint  = var.smn_subscription_endpoint
  protocol  = var.smn_subscription_protocol
  remark    = var.smn_subscription_remark
}

resource "huaweicloud_antiddos_basic" "test" {
  traffic_threshold = var.antiddos_traffic_threshold
  eip_id            = huaweicloud_vpc_eip.test.id
  topic_urn         = huaweicloud_smn_topic.test.id
}
