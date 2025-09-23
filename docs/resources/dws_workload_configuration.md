---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_workload_configuration"
description: |-
  Manages a workload configuration resource under specified DWS cluster within HuaweiCloud.
---

# huaweicloud_dws_workload_configuration

Manages a workload configuration resource under specified DWS cluster within HuaweiCloud.

## Example Usage

```hcl
variable "dws_cluster_id" {}

resource "huaweicloud_dws_workload_configuration" "test" {
  cluster_id          = var.dws_cluster_id
  workload_switch     = "on"
  max_concurrency_num = "100"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `cluster_id` - (Required, String, ForceNew) Specifies the DWS cluster ID.
  Changing this creates a new resource.

* `workload_switch` - (Required, String) Specifies the workload management switch.  
  The valid value are as follows:
  + **on**
  + **off**

  -> If this parameter is set to **off**, all resource management functions will be unavailable.

* `max_concurrency_num` - (Optional, String) Specifies the maximum number of concurrent tasks on a single CN.  
  The valid value ranges from `-1` to `2,147,483,647`, `-1` and `0` means unlimited.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also `cluster_id`.

## Import

The resource can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_dws_workload_configuration.test <id>
```
