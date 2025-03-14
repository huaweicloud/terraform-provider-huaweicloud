---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_network"
description: |-
  Manages a CCI Network resource within HuaweiCloud.
---

# huaweicloud_cciv2_network

Manages a CCI Network resource within HuaweiCloud.

## Example Usage

```hcl

```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String) Specifies the name of the CCI network.

* `namespace` - (Required, String) Specifies the namespace.

* `annotations` - (Optional, Map) Specifies the annotations of the CCI network.

* `ip_families` - (Optional, List) Specifies the IP families of the CCI network.

* `security_group_ids` - (Optional, List) Specifies the security group IDs of the CCI network.

* `subnets` - (Optional, List) Specifies the subnets of the CCI network.
  The [subnets](#block--subnets) structure is documented below.

<a name="block--subnets"></a>
The `subnets` block supports:

* `subnet_id` - (Optional, String) Specifies the subnet ID of the CCI network.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `api_version` - The API version of the CCI network.

* `kind` - The kind of the CCI network.

* `resource_version` - The resource version of the namespace.

* `self_link` - The self link of the namespace.

* `status` - The status of the namespace.

* `uid` - The uid of the namespace.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The xxx can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_cciv2_network.test <id>
```
