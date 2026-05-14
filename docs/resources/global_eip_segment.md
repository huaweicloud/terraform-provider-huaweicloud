---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_global_eip_segment"
description: |-
  Manages a global EIP segment resource within HuaweiCloud.
---

# huaweicloud_global_eip_segment

Manages a global EIP segment resource within HuaweiCloud.

## Example Usage

```hcl
variable "geip_pool_name" {}
variable "access_site" {}

resource "huaweicloud_global_eip_segment" "test" {
  geip_pool_name = var.geip_pool_name
  access_site    = var.access_site
  mask           = 29

  tags {
    key   = "foo"
    value = "bar"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the resource is located.  
  If omitted, the provider-level region will be used. Change this parameter will create a new resource.

* `geip_pool_name` - (Required, String, NonUpdatable) Specifies the global EIP pool name.  
  The value can be queried use the `huaweicloud_global_eip_pools` data source.

* `access_site` - (Required, String, NonUpdatable) Specifies the access site name.  
  The value can be queried use the `huaweicloud_global_eip_pools` data source.

* `mask` - (Required, Int, NonUpdatable) Specifies the mask length of the segment.  
  The value can be queried use the `huaweicloud_global_eip_segment_support_masks` data source

* `name` - (Optional, String) Specifies the name of the global EIP segment.

* `description` - (Optional, String) Specifies the description of the global EIP segment.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID to which the global
  EIP segment belongs. If omitted, the default enterprise project is used.

* `tags` - (Optional, List) Specifies the tags of the global EIP segment.

  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - (Required, String) Specifies the tag key.

* `value` - (Optional, String) Specifies the tag value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also as the segment ID.

* `domain_id` - The domain ID.

* `isp` - The ISP line.

* `ip_version` - The IP version.

* `cidr` - The IPv4 CIDR block.

* `cidr_v6` - The IPv6 CIDR block.

* `freezen` - Whether the segment is frozen.

* `status` - The status of the segment.

* `created_at` - The creation time of the segment.

* `updated_at` - The update time of the segment.

* `is_pre_paid` - Whether the segment is prepaid.

* `is_charged` - Whether the segment is charged.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.

## Import

The global EIP segment can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_global_eip_segment.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `mask`. It is generally recommended
running `terraform plan` after importing the resource. You can then decide if changes should be applied to the resource,
or the resource definition should be updated to align with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_global_eip_segment" "test" {
  ...

  lifecycle {
    ignore_changes = [
      mask,
    ]
  }
}
```
