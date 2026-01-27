---
subcategory: "Cloud Container Engine Autopilot (CCE Autopilot)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_autopilot_cluster_log_configs"
description: |-
  Use this data source to get the list of CCE Autopilot clusters log configs within huaweicloud.
---

# huaweicloud_cce_autopilot_cluster_log_configs

Use this data source to get the list of CCE Autopilot clusters log configs within huaweicloud.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_cce_autopilot_cluster_log_configs" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the cluster ID.

* `type` - (Optional, String) Specifies the type of log config. Value options:
  + **control**: specifies the logs of the control plane components.
  + **audit**: specifies the audit logs on the control plane.
  + **system-addon**: s ecifies the logs of the system add-ons.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `ttl_in_days` - The storage duration.

* `log_configs` - Log config information.

  The [log_configs](#log_configs_struct) structure is documented below.

<a name="log_configs_struct"></a>
The `log_configs` block supports:

* `name` - Last config server name.

* `enable` - Whether to enable log collection.

* `type` - The type of log config.
