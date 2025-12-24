resource "huaweicloud_cdn_rule_engine_rule" "test" {
  domain_name = var.domain_name
  name        = var.rule_name
  status      = var.rule_status
  priority    = var.rule_priority
  conditions  = var.conditions != "" ? var.conditions : null

  dynamic "actions" {
    for_each = var.cache_rule != null ? [var.cache_rule] : []

    content {
      cache_rule {
        ttl           = actions.value.ttl
        ttl_unit      = actions.value.ttl_unit
        follow_origin = lookup(actions.value, "follow_origin", null)
        force_cache   = lookup(actions.value, "force_cache", null)
      }
    }
  }

  dynamic "actions" {
    for_each = var.access_control != null ? [var.access_control] : []

    content {
      access_control {
        type = actions.value.type
      }
    }
  }

  dynamic "actions" {
    for_each = length(var.http_response_headers) > 0 ? var.http_response_headers : []

    content {
      http_response_header {
        name   = actions.value.name
        value  = actions.value.value
        action = actions.value.action
      }
    }
  }

  dynamic "actions" {
    for_each = var.browser_cache_rule != null ? [var.browser_cache_rule] : []

    content {
      browser_cache_rule {
        cache_type = actions.value.cache_type
      }
    }
  }

  dynamic "actions" {
    for_each = var.request_url_rewrite != null ? [var.request_url_rewrite] : []

    content {
      request_url_rewrite {
        execution_mode = actions.value.execution_mode
        redirect_url   = actions.value.redirect_url
      }
    }
  }

  dynamic "actions" {
    for_each = length(var.flexible_origins) > 0 ? var.flexible_origins : []

    content {
      flexible_origin {
        sources_type      = actions.value.sources_type
        ip_or_domain      = actions.value.ip_or_domain
        priority          = actions.value.priority
        weight            = actions.value.weight
        http_port         = lookup(actions.value, "http_port", null)
        https_port        = lookup(actions.value, "https_port", null)
        origin_protocol   = lookup(actions.value, "origin_protocol", null)
        host_name         = lookup(actions.value, "host_name", null)
        obs_bucket_type   = lookup(actions.value, "obs_bucket_type", null)
        bucket_access_key = lookup(actions.value, "bucket_access_key", null)
        bucket_secret_key = lookup(actions.value, "bucket_secret_key", null)
        bucket_region     = lookup(actions.value, "bucket_region", null)
        bucket_name       = lookup(actions.value, "bucket_name", null)
      }
    }
  }

  dynamic "actions" {
    for_each = length(var.origin_request_headers) > 0 ? var.origin_request_headers : []

    content {
      origin_request_header {
        action = actions.value.action
        name   = actions.value.name
        value  = lookup(actions.value, "value", null)
      }
    }
  }

  dynamic "actions" {
    for_each = var.origin_request_url_rewrite != null ? [var.origin_request_url_rewrite] : []

    content {
      origin_request_url_rewrite {
        rewrite_type = actions.value.rewrite_type
        target_url   = actions.value.target_url
      }
    }
  }

  dynamic "actions" {
    for_each = var.origin_range != null ? [var.origin_range] : []

    content {
      origin_range {
        status = actions.value.status
      }
    }
  }

  dynamic "actions" {
    for_each = var.request_limit_rule != null ? [var.request_limit_rule] : []

    content {
      request_limit_rule {
        limit_rate_after = actions.value.limit_rate_after
        limit_rate_value = actions.value.limit_rate_value
      }
    }
  }

  dynamic "actions" {
    for_each = length(var.error_code_cache) > 0 ? var.error_code_cache : []

    content {
      error_code_cache {
        code = actions.value.code
        ttl  = actions.value.ttl
      }
    }
  }

  lifecycle {
    ignore_changes = [
      conditions
    ]
  }
}
