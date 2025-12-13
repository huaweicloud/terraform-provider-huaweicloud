resource "huaweicloud_waf_cloud_instance" "test" {
  resource_spec_code = var.cloud_instance_resource_spec_code

  dynamic "bandwidth_expack_product" {
    for_each = var.cloud_instance_bandwidth_expack_product

    content {
      resource_size = bandwidth_expack_product.value["resource_size"]
    }
  }

  dynamic "domain_expack_product" {
    for_each = var.cloud_instance_domain_expack_product

    content {
      resource_size = domain_expack_product.value["resource_size"]
    }
  }

  dynamic "rule_expack_product" {
    for_each = var.cloud_instance_rule_expack_product

    content {
      resource_size = rule_expack_product.value["resource_size"]
    }
  }

  charging_mode         = var.cloud_instance_charging_mode
  period_unit           = var.cloud_instance_period_unit
  period                = var.cloud_instance_period
  auto_renew            = var.cloud_instance_auto_renew
  enterprise_project_id = var.enterprise_project_id
}

resource "huaweicloud_waf_domain" "test" {
  domain                = var.cloud_domain
  certificate_id        = var.cloud_certificate_id
  certificate_name      = var.cloud_certificate_name
  proxy                 = var.cloud_proxy
  enterprise_project_id = var.enterprise_project_id
  description           = var.cloud_description
  website_name          = var.cloud_website_name
  protect_status        = var.cloud_protect_status
  forward_header_map    = var.cloud_forward_header_map

  dynamic "custom_page" {
    for_each = var.cloud_custom_page

    content {
      http_return_code = custom_page.value["http_return_code"]
      block_page_type  = custom_page.value["block_page_type"]
      page_content     = custom_page.value["page_content"]
    }
  }

  dynamic "timeout_settings" {
    for_each = var.cloud_timeout_settings

    content {
      connection_timeout = timeout_settings.value["connection_timeout"]
      read_timeout       = timeout_settings.value["read_timeout"]
      write_timeout      = timeout_settings.value["write_timeout"]
    }
  }

  dynamic "traffic_mark" {
    for_each = var.cloud_traffic_mark

    content {
      ip_tags     = traffic_mark.value["ip_tags"]
      session_tag = traffic_mark.value["session_tag"]
      user_tag    = traffic_mark.value["user_tag"]
    }
  }

  dynamic "server" {
    for_each = var.cloud_server

    content {
      client_protocol = server.value["client_protocol"]
      server_protocol = server.value["server_protocol"]
      address         = server.value["address"]
      port            = server.value["port"]
      type            = server.value["type"]
      weight          = server.value["weight"]
    }
  }

  depends_on = [
    huaweicloud_waf_cloud_instance.test
  ]
}
