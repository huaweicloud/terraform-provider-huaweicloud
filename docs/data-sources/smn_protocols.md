---
subcategory: "Simple Message Notification (SMN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_smn_protocols"
description: |-
  Use this data source to get the list of protocols supported by SMN.
---

# huaweicloud_smn_protocols

Use this data source to get the list of protocols supported by SMN.

## Example Usage

```hcl
data "huaweicloud_smn_protocols" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `protocols` - Indicates the list of protocol.
