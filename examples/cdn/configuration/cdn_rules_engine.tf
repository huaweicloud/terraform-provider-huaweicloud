# 此示例展示了如何创建和管理CDN规则引擎规则

# 华为云CDN加速域名创建示例
resource "huaweicloud_cdn_domain" "example" {
  name = "terraform-test.nanguapi.com"
  type = "web"
  sources {
    origin = "1.2.3.4"
    origin_type = "ipaddr"
    active = "1"
    http_port = 80
    https_port = 443
  }
  service_area = "mainland_china"
}

# 创建CDN规则引擎规则 - 复杂条件示例
resource "huaweicloud_cdn_rules_engine" "complex_rule" {
  domain_name = huaweicloud_cdn_domain.example.name
  name        = "complex_rule"
  status      = "on"
  priority    = 8

  conditions {
    match {
      logic = "and"
      criteria {
        match_target_type = "method"
        match_type        = "contains"
        match_pattern     = ["POST", "PUT"]
        negate           = false
        case_sensitive   = false
      }
      
      criteria {
        match_target_type = "header"
        match_target_name = "User-Agent"
        match_type        = "contains"
        match_pattern     = ["Mobile"]
        negate           = false
        case_sensitive   = false
      }
    }
  }

  actions {
    access_control {
      type = "trust"
    }
  }

  actions {
    http_response_header {
      name   = "X-Custom-Header"
      value  = "mobile-user"
      action = "set"
    }
  }
}

# 查询所有已创建规则
data "huaweicloud_cdn_rules_engine" "all_rules" {
  domain_name = huaweicloud_cdn_domain.example.name
}

output "rules" {
  value = data.huaweicloud_cdn_rules_engine.all_rules
}

