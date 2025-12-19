resource "huaweicloud_cdn_domain" "test" {
  name         = var.domain_name
  type         = var.domain_type
  service_area = var.service_area

  sources {
    origin      = var.origin_server
    origin_type = var.origin_type
    active      = 1
    http_port   = var.http_port
    https_port  = var.https_port
  }

  configs {
    origin_protocol               = var.origin_protocol
    ipv6_enable                   = var.ipv6_enable
    range_based_retrieval_enabled = var.range_based_retrieval_enabled
    description                   = var.domain_description

    dynamic "https_settings" {
      for_each = var.https_enabled ? [1] : []

      content {
        certificate_name     = var.https_enabled ? var.certificate_name : null
        certificate_source   = var.https_enabled ? var.certificate_source : null
        certificate_body     = var.https_enabled && var.certificate_body_path != "" ? file(var.certificate_body_path) : null
        private_key          = var.https_enabled && var.private_key_path != "" ? file(var.private_key_path) : null
        https_enabled        = var.https_enabled
        http2_enabled        = var.http2_enabled
        ocsp_stapling_status = var.ocsp_stapling_status
      }
    }
  }

  dynamic "cache_settings" {
    for_each = length(var.cache_rules) > 0 ? [var.cache_rules] : []

    content {
      dynamic "rules" {
        for_each = cache_settings.value

        content {
          rule_type           = rules.value.rule_type
          ttl                 = rules.value.ttl
          ttl_type            = rules.value.ttl_type
          priority            = rules.value.priority
          content             = rules.value.content
          url_parameter_type  = lookup(rules.value, "url_parameter_type", null)
          url_parameter_value = lookup(rules.value, "url_parameter_value", null)
        }
      }
    }
  }

  tags = var.domain_tags
}
