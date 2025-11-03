---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aadv2_instances"
description: |-
  Use this data source to get the list of Advanced Anti-DDos v2 instances within HuaweiCloud.
---

# huaweicloud_aadv2_instances

Use this data source to get the list of Advanced Anti-DDos v2 instances within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_aadv2_instances" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_access_type` - (Optional, String) Specifies the access type.  
  The valid values are as follows:
  + **0**: Website instance.
  + **1**: IP access instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `items` - The list of instances.

  The [items](#items_struct) structure is documented below.

<a name="items_struct"></a>
The `items` block supports:

* `instance_id` - The instance ID.

* `instance_name` - The name of the instance.

* `enterprise_project_id` - The enterprise project ID.

* `instance_access_type` - The instance access type.

* `pp_support` - Whether PP is supported. `1` indicates supported, and `0` indicates not supported.

* `pp_enable` - Whether the customer has enabled PP. `1` indicates enabled, and `0` indicates disabled.

* `overseas_type` - The protection region. `0` indicates mainland China, and `1` indicates overseas.

* `vips` - The high-defense IP information.

  The [vips](#vips_struct) structure is documented below.

<a name="vips_struct"></a>
The `vips` block supports:

* `ip` - The IP address.

* `isp` - The line.
