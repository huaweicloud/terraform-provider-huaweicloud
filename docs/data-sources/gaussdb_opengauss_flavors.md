---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_flavors"
description: |-
  Use this data source to get the list of GaussDB OpenGauss flavors.
---

# huaweicloud_gaussdb_opengauss_flavors

Use this data source to get the list of GaussDB OpenGauss flavors.

## Example Usage

```hcl
data "huaweicloud_gaussdb_opengauss_flavors" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `version` - (Optional, String) Specifies the version of the GaussDB OpenGauss.

* `spec_code` - (Optional, String) Specifies the specification code.

* `ha_mode` - (Optional, String) Specifies the instance type.
  Value options:
  + **centralization_standard**: Primary/standby
  + **enterprise**: Distributed (independent deployment)

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `flavors` - Indicates the list of the flavors.

  The [flavors](#flavors_struct) structure is documented below.

<a name="flavors_struct"></a>
The `flavors` block supports:

* `name` - Indicates the DB engine.

* `spec_code` - Indicates the specification code.

* `version` - Indicates the DB engine version supported by the specifications.

* `group_type` - Performance specifications.
  The value can be:
  + **normal**: dedicated (1:8)
  + **normal2**: dedicated (1:4)
  + **armFlavors**: Kunpeng dedicated (1:8)
  + **exclusive**: Dedicated (1:4) It is only suitable for primary/standby instances of the basic edition.
  + **armExclusive**: Kunpeng dedicated (1:4) It is only suitable for primary/standby instances of the basic edition.
  + **economical**: Favored (1:4)
  + **economical2**: Favored (1:8)
  + **armFlavors2**: Kunpeng dedicated (1:4)
  + **general**: general-purpose (1:4)

* `vcpus` - Indicates the number of vCPUs.

* `ram` - Indicates the memory size in GB.

* `availability_zone` - Indicates the AZ supported by the specifications.

* `az_status` - Indicates the key/value pairs of the availability zone status.
  **key** indicates the AZ ID, and **value** indicates the specification status in the AZ.
