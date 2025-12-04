---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_zone_retrieval_verify"
description: |-
  Manages a DNS zone retrieval verify resource within HuaweiCloud.
---

# huaweicloud_dns_zone_retrieval_verify

Manages a DNS zone retrieval verify resource within HuaweiCloud.

## Example Usage

```hcl
variable "retrieval_id" {}

resource "huaweicloud_dns_zone_retrieval_verify" "test" {
  retrieval_id = var.retrieval_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `retrieval_id` - (Required, String, NonUpdatable) Specifies the retrieval ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `message` - Indicates the message.
