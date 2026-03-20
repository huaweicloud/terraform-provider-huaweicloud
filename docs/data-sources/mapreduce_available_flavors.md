---
subcategory: "MapReduce Service (MRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_mapreduce_available_flavors"
description: |-
  Use this data source to get the available flavors of MRS cluster within HuaweiCloud.
---

# huaweicloud_mapreduce_available_flavors

Use this data source to get the available flavors of MRS cluster within HuaweiCloud.

## Example Usage

```hcl
variable "version_name" {}

data "huaweicloud_mapreduce_available_flavors" "test" {
  version_name = var.version_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the flavors are located.  
  If omitted, the provider-level region will be used.

* `version_name` - (Required, String) Specifies the version name of the cluster. e.g. `MRS 3.5.0-LTS`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `available_flavors` - The list of available flavors.  
  The [available_flavors](#available_flavors_struct) structure is documented below.

<a name="available_flavors_struct"></a>
The `available_flavors` block supports:

* `az_code` - The availability zone code.

* `az_name` - The availability zone name.

* `master` - The list of flavors supported by master nodes.  
  The [master](#node_available_flavors_struct) structure is documented below.

* `core` - The list of flavors supported by core nodes.  
  The [core](#node_available_flavors_struct) structure is documented below.

* `task` - The list of flavors supported by task nodes.  
  The [task](#node_available_flavors_struct) structure is documented below.

<a name="node_available_flavors_struct"></a>
The `master`, `core` and `task` blocks support:

* `flavor_name` - The flavor name.
