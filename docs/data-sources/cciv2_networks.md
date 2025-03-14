---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_networks"
description: |-
  Use this data source to get the list of CCI networks within HuaweiCloud.
---

# huaweicloud_cciv2_networks

Use this data source to get the list of CCI networks within HuaweiCloud.

## Example Usage

```hcl

```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `networks` - <!-- please add the description of the attribute -->
  The [networks](#attrblock--networks) structure is documented below.

<a name="attrblock--networks"></a>
The `networks` block supports:

* `annotations` - The annotations of the namespace.

* `creation_timestamp` - The creation timestamp of the namespace.

* `ip_families` - Specifies the IP families of the CCI network.

* `labels` - The labels of the namespace.

* `name` - The name of the namespace.

* `namespace` - The name of the namespace.

* `resource_version` - The resource version of the namespace.

* `security_group_ids` - Specifies the security group IDs of the CCI network.

* `self_link` - The self link of the namespace.

* `status` - The status of the namespace.

* `subnets` - Specifies the subnets of the CCI network.
  The [subnets](#attrblock--networks--subnets) structure is documented below.

* `uid` - The uid of the namespace.

<a name="attrblock--networks--subnets"></a>
The `subnets` block supports:

* `subnet_id` - Specifies the subnet ID of the CCI network.
