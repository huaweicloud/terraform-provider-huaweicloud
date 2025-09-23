---
subcategory: "Distributed Database Middleware (DDM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ddm_flavors"
description: ""
---

# huaweicloud_ddm_flavors

Use this data source to get the list of DDM flavors.

## Example Usage

```hcl
data "huaweicloud_ddm_engines" "test" {
  version = "3.0.8.5"
}

data "huaweicloud_ddm_flavors" "test" {
  engine_id = data.huaweicloud_ddm_engines.test.engines[0].id
  cpu_arch  = "X86"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `engine_id` - (Required, String) Specifies the ID of an engine.

* `cpu_arch` - (Optional, String) Specifies the compute resource architecture type. The options are **X86** and **ARM**.

* `code` - (Optional, String) Specifies the VM flavor types recorded in DDM.

* `vcpus` - (Optional, Int) Specifies the number of CPUs.

* `memory` - (Optional, Int) Specifies the memory size. Unit GB.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `flavors` - Indicates the list of DDM compute flavors.
  The [Flavor](#DdmFlavors_Flavor) structure is documented below.

<a name="DdmFlavors_Flavor"></a>
The `Flavor` block supports:

* `id` - Indicates the ID of a flavor.

* `cpu_arch` - Indicates the compute resource architecture type.

* `code` - Indicates the VM flavor types recorded in DDM.

* `vcpus` - Indicates the number of CPUs.

* `memory` - Indicates the memory size.
