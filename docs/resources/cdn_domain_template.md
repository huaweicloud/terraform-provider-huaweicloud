---
subcategory: "Content Delivery Network (CDN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_domain_template"
description: |-
  Manages a CDN domain template resource within HuaweiCloud.
---

# huaweicloud_cdn_domain_template

Manages a CDN domain template resource within HuaweiCloud.

## Example Usage

```hcl
variable "template_name" {}

resource "huaweicloud_cdn_domain_template" "test" {
  name        = var.template_name
  description = "Created by terraform"
  configs     = jsonencode({
    "cache_rules": [
      {
        "force_cache": "on",
        "follow_origin": "off",
        "match_type": "all",
        "priority": 1,
        "stale_while_revalidate": "off",
        "ttl": 20,
        "ttl_unit": "d",
        "url_parameter_type": "full_url",
        "url_parameter_value": ""
      }
    ],
    "http_response_header": [
      {
        "action": "set",
        "name": "Content-Disposition",
        "value": "1235"
      }
    ],
    "origin_follow302_status": "off",
    "compress": {
      "type": "gzip,br",
      "status": "on",
      "file_type": ".js,.html,.css,.xml,.json,.shtml,.htmx"
    },
    "origin_range_status": "on",
    "referer": {
      "type": "black",
      "value": "1.2.1.1",
      "include_empty": false
    },
    "ip_filter": {
      "type": "white",
      "value": "1.1.2.2"
    },
    "user_agent_filter": {
      "type": "white",
      "ua_list": [
        "1.1.3.3"
      ],
      "include_empty": false
    },
    "flow_limit_strategy": [
      {
        "strategy_type": "instant",
        "item_type": "bandwidth",
        "limit_value": 1000001,
        "alarm_percent_threshold": null,
        "ban_time": 60
      }
    ]
  })
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the name of the domain template.

* `configs` - (Required, String) Specifies the configuration of the domain template, in JSON format.  
  The configuration includes cache rules, HTTP response headers, compression settings, origin settings, access control,
  and other CDN configurations.

* `description` - (Optional, String) Specifies the description of the domain template.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `type` - The type of the domain template.
  + **1** - System preset template.
  + **2** - Tenant custom template.

* `account_id` - The account ID.

* `create_time` - The creation time of the domain template, in RFC3339 format.

* `modify_time` - The modification time of the domain template, in RFC3339 format.

## Import

The domain template can be imported using the `id` or `name`, e.g.

```bash
$ terraform import huaweicloud_cdn_domain_template.test <id>
```

or

```bash
$ terraform import huaweicloud_cdn_domain_template.test <name>
```
