# Rule Engine Rule Configuration
domain_name   = "example.com"
rule_name     = "test-rule-engine"
rule_status   = "on"
rule_priority = 1
# Conditions in JSON format (as a string)
# You need to provide the JSON string directly, or use heredoc syntax for multi-line strings.
conditions    = <<-JSON
{
  "match": {
    "logic": "and",
    "criteria": [
      {
        "match_target_type": "path",
        "match_type": "contains",
        "match_pattern": ["/api/"],
        "negate": false,
        "case_sensitive": true
      }
    ]
  }
}
JSON

# Origin request URL rewrite
origin_request_url_rewrite = {
  rewrite_type = "simple"
  target_url   = "/api/v2"
}
# Cache rule configuration
cache_rule                 = {
  ttl           = 10
  ttl_unit      = "m"
  follow_origin = "min_ttl"
  force_cache   = "off"
}
# Access control configuration
access_control             = {
  type = "trust"
}
# Browser cache rule
browser_cache_rule         = {
  cache_type = "follow_origin"
}
# Reuest URL rewrite
request_url_rewrite        = {
  execution_mode = "break"
  redirect_url   = "/new-path"
}
# Origin range
origin_range               = {
  status = "on"
}

# Request limit rule
request_limit_rule     = {
  limit_rate_after = 2
  limit_rate_value = 1048576
}
# Origin request headers
origin_request_headers = [
  {
    action = "set"
    name   = "X-Real-IP"
    value  = "$realip_from_header"
  }
]
# Flexible origins
flexible_origins       = [
  {
    sources_type    = "domain"
    ip_or_domain    = "target.domain.com"
    priority        = 1
    weight          = 10
    http_port       = 80
    https_port      = 443
    origin_protocol = "follow"
    host_name       = "target.domain.com"
  }
]
# HTTP response headers
http_response_headers  = [
  {
    name   = "Access-Control-Allow-Origin"
    value  = "*"
    action = "set"
  }
]
# Error code cache
error_code_cache       = [
  {
    code = 400
    ttl  = 60
  }
]
