---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_cluster_log_config"
description: |-
  Use this resource to manage the log config of a CCE cluster within HuaweiCloud.
---

# huaweicloud_cce_cluster_log_config

Use this resource to manage the log config of a CCE cluster within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "cluster_id" {}

resource "huaweicloud_cce_cluster_log_config" "test" {
  cluster_id  = var.cluster_id
  ttl_in_days = 3

  log_configs {
    name   = "kube-apiserver"
    enable = true
  }

  log_configs {
    name   = "kube-controller-manager"
    enable = false
  }

  log_configs {
    name   = "kube-scheduler"
    enable = false
  }

  log_configs {
    name   = "audit"
    enable = true
  }
}

```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the cluster log config resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the cluster ID.

* `ttl_in_days` - (Optional, Int) Specifies the log keeping days, default to `7`.

* `log_configs` - (Optional, List) Specifies the list of log configs.
  The [log_configs](#log_configs) structure is documented below.

<a name="log_configs"></a>
The `log_configs` block supports:

* `name` - (Optional, String) Specifies the log type.
  
* `enable` - (Optional, Bool) Specifies whether to collect the log.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, the value is the cluster ID.

## Import

The cluster log config can be imported using the cluster ID, e.g.

```bash
$ terraform import huaweicloud_cce_cluster_log_config.test <cluster_id>
```
