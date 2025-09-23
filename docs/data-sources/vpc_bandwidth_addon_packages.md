---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_bandwidth_addon_packages"
description: |-
  Use this data source to get a list of bandwidth add-on packages.
---

# huaweicloud_vpc_bandwidth_addon_packages

Use this data source to get a list of bandwidth add-on packages.

## Example Usage

```hcl
data "huaweicloud_vpc_bandwidth_addon_packages" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `bandwidth_pkgs` - Indicates the bandwidth packages list.

  The [bandwidth_pkgs](#bandwidth_pkgs_struct) structure is documented below.

<a name="bandwidth_pkgs_struct"></a>
The `bandwidth_pkgs` block supports:

* `resource_id` - Indicates the refueling package ID.

* `resource_name` - Indicates the refueling package name.

* `bandwidth_id` - Indicates the  ID of the original bandwidth bound to the add-on package.

* `pkg_size` - Indicates the size of the add-on packet.
  And it is the increased bandwidth over the original bandwidth.

* `billing_info` - Indicates the information about an add-on package order.

* `status` - Indicates the resource status of an add-on package.
  Value can be **pending**, **active**, **completed**, **error**.

* `start_time` - Indicates the start time when an add-on package takes effect.

* `end_time` - Indicates the end time when an add-on package takes effect.

* `processed_time` - Indicates the resource creation time.
