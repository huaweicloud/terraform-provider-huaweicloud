---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_resolver_access_log"
description: |-
  Manages a DNS resolver access log resource within HuaweiCloud.
---

# huaweicloud_dns_resolver_access_log

Manages a DNS resolver access log resource within HuaweiCloud.

## Example Usage

```hcl
variable "lts_group_id" {}
variable "lts_stream_id" {}
variable "vpc_ids" {
  type = list(string)
}

resource "huaweicloud_dns_resolver_access_log" "test" {
  lts_group_id = var.lts_group_id
  lts_topic_id = var.lts_stream_id
  vpc_ids      = var.vpc_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the resolver access log is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `lts_group_id` - (Required, String, NonUpdatable) Specifies the ID of the log group.

* `lts_topic_id` - (Required, String, NonUpdatable) Specifies the ID of the log stream.

* `vpc_ids` - (Required, List) Specifies the list of VPC IDs associated with the resolver access log.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

Resolver access log can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_dns_resolver_access_log.test <id>
```
