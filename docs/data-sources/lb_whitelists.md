---
subcategory: "Elastic Load Balance (ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lb_whitelists"
description: |-
  Use this data source to get the list of ELB whitelists.
---

# huaweicloud_lb_whitelists

Use this data source to get the list of ELB whitelists.

## Example Usage

```hcl
data "huaweicloud_lb_whitelists" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `whitelist_id` - (Optional, String) Specifies the whitelist ID.

* `enable_whitelist` - (Optional, String) Specifies whether to enable access control. Value options:
  + **true**: Access control is enabled.
  + **false**: Access control is disabled.

* `listener_id` - (Optional, String) Specifies the ID of the listener to which the whitelist is added.

* `whitelist` - (Optional, String) Specifies the IP addresses in the whitelist.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `whitelists` - Indicates the list of whitelist.

  The [whitelists](#whitelists_struct) structure is documented below.

<a name="whitelists_struct"></a>
The `whitelists` block supports:

* `id` - Indicates the whitelist ID.

* `listener_id` - Indicates the ID of the listener to which the whitelist is added.

* `enable_whitelist` - Indicates whether to enable access control.

* `whitelist` - Indicates the IP addresses in the whitelist.
