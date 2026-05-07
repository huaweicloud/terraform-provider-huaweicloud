---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_eip_publicip_count"
description: |-
  Use this data source to get the count of elastic public IPs.
---

# huaweicloud_eip_publicip_count

Use this data source to get the count of elastic public IPs.

## Example Usage

```hcl
data "huaweicloud_eip_publicip_count" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `elasticip_size` - Indicates the number of elastic public IPs.
