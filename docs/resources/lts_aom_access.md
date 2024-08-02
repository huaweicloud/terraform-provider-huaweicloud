---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_aom_access"
description: |-
  Manages an AOM to LTS log mapping rule resource within HuaweiCloud.
---

# huaweicloud_lts_aom_access

Manages an AOM to LTS log mapping rule resource within HuaweiCloud.

-> The resource of connecting AOM logs to LTS is currently restricted. Please submit a service ticket to open this
feature for you. Refer to
[How to submit a service ticket](https://support.huaweicloud.com/intl/en-us/usermanual-ticket/topic_0065264094.html).

## Example Usage

### Creating with all workloads

```hcl
variable "cluster_id" {}
variable "cluster_name" {}
variable "log_group_id" {}
variable "log_group_name" {}
variable "log_stream_id" {}
variable "log_stream_name" {}

resource "huaweicloud_lts_aom_access" "test" {
  name         = "test_name"
  cluster_id   = var.cluster_id
  cluster_name = var.cluster_name
  namespace    = "default"
  workloads    = ["__ALL_DEPLOYMENTS__"]
  
  access_rules {
    file_name       = "/test/*"
    log_group_id    = var.log_group_id
    log_group_name  = var.log_group_name
    log_stream_id   = var.log_stream_id
    log_stream_name = var.log_stream_name
  }

  access_rules {
    file_name       = "/demo/demo.log"
    log_group_id    = var.log_group_id
    log_group_name  = var.log_group_name
    log_stream_id   = var.log_stream_id
    log_stream_name = var.log_stream_name
  }
}
```

### Creating with specify workloads

```hcl
variable "cluster_id" {}
variable "cluster_name" {}
variable "log_group_id" {}
variable "log_group_name" {}
variable "log_stream_id" {}
variable "log_stream_name" {}
variable "workload" {}

resource "huaweicloud_lts_aom_access" "test" {
  name         = "test_name"
  cluster_id   = var.cluster_id
  cluster_name = var.cluster_name
  namespace    = "default"
  workloads    = [var.workload]

  access_rules {
    file_name       = "__ALL_FILES__"
    log_group_id    = var.log_group_id
    log_group_name  = var.log_group_name
    log_stream_id   = var.log_stream_id
    log_stream_name = var.log_stream_name
  }
}
```

### Creating with CCI cluster

```hcl
variable "log_group_id" {}
variable "log_group_name" {}
variable "log_stream_id" {}
variable "log_stream_name" {}

resource "huaweicloud_lts_aom_access" "test" {
  name         = "test_name"
  cluster_id   = "CCI-ClusterID"
  cluster_name = "CCI-Cluster"
  namespace    = "default"
  workloads    = ["__ALL_DEPLOYMENTS__"]

  access_rules {
    file_name       = "__ALL_FILES__"
    log_group_id    = var.log_group_id
    log_group_name  = var.log_group_name
    log_stream_id   = var.log_stream_id
    log_stream_name = var.log_stream_name
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the unique rule name. The name consists of `1` to `100` characters,
  including letters, digits, underscores (_), hyphens (-) and dots (.).

* `cluster_id` - (Required, String) Specifies the CCE or CCI cluster ID. It is fixed to **CCI-ClusterID** for CCI.

* `cluster_name` - (Required, String) Specifies the CCE or CCI cluster name. It is fixed to **CCI-Cluster** for CCI.

* `namespace` - (Required, String) Specifies the namespace.

* `workloads` - (Required, List) Specifies the workloads.
  + When creating with all workloads, this field should be `["__ALL_DEPLOYMENTS__"]`.
  + When creating with specify workloads, this field should be the list of workloads.

* `access_rules` - (Required, List) Specifies the access log details.
The [access_rules](#AOMAccess_access_rules) structure is documented below.

* `container_name` - (Optional, String) Specifies the container name.

<a name="AOMAccess_access_rules"></a>
The `access_rules` block supports:

* `file_name` - (Required, String) Specifies the path name.
  + When collecting access all logs, set this field to `__ALL_FILES__`.
  + When collecting specify log paths, the matching rule should be `^\/[A-Za-z0-9.*_\/-]+|stdout\.log|`, such as
  `/test/*` or `/test/demo.log`. Up to two asterisks (*) are allowed.

* `log_group_id` - (Required, String) Specifies the log group ID.

* `log_group_name` - (Required, String) Specifies the log group name.

* `log_stream_id` - (Required, String) Specifies the log stream ID.

* `log_stream_name` - (Required, String) Specifies the log stream name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The AOM to LTS log mapping rule can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_lts_aom_access.test <id>
```
