---
subcategory: "Content Delivery Network (CDN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_domain_batch_copy"
description: |-
  Use this resource to batch copy the configurations of domain within HuaweiCloud.
---

# huaweicloud_cdn_domain_batch_copy

Use this resource to batch copy the configurations of domain within HuaweiCloud.

-> This resource is only a one-time action resource for performing domain batch copy operation. Deleting this resource
   will not clear the corresponding request record, but will only remove the resource information from the tfstate
   file.

## Example Usage

```hcl
variable "source_domain" {}
variable "target_domains" {}
variable "config_list" {
  type = list(string)
}

resource "huaweicloud_cdn_domain_batch_copy" "test" {
  source_domain  = var.source_domain
  target_domains = var.target_domain
  config_list    = var.config_list
}
```

## Argument Reference

* `source_domain` - (Required, String, NonUpdatable) Specifies the source domain whose configuration will be copied.

* `target_domains` - (Required, String, NonUpdatable) Specifies the target domain names to which configurations will be
  copied.

* `config_list` - (Required, List, NonUpdatable) Specifies the configuration items to copy.  
  The valid values are as follows:
  + **origin_request_header**
  + **http_response_header**
  + **url_auth**
  + **user_agent_black_and_white_list**
  + **ipv6Accelerate**
  + **range_status**
  + **cache_rules**
  + **follow_302_status**
  + **sources**
  + **compress**
  + **referer**
  + **ip_black_and_white_list**
  + **browserCacheRules**
  + **cacheValidErrorCode**

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `id` - The ID of resource.
