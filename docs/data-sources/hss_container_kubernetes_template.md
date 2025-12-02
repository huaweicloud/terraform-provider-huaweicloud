---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_container_kubernetes_template"
description: |-
  Use this data source to get the container kubernetes template of HSS within HuaweiCloud.
---

# huaweicloud_hss_container_kubernetes_template

Use this data source to get the container kubernetes template of HSS within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_container_kubernetes_template" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter only needs to be configured after the Enterprise Project feature is enabled.
  For enterprise users, if omitted, default enterprise project will be used.
  Value **0** means default enterprise project.
  Value **all_granted_eps** means all enterprise projects to which the user has been granted access.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `is_default` - Whether the template is a default template.

* `runtime_info` - The container runtime configuration.

  The [runtime_info](#runtime_info_struct) structure is documented below.

* `schedule_info` - The node scheduling information.

  The [schedule_info](#schedule_info_struct) structure is documented below.

<a name="runtime_info_struct"></a>
The `runtime_info` block supports:

* `runtime_name` - The runtime name. Valid values are:
  + **crio_endpoint**: CRIO
  + **containerd_endpoint**: Containerd
  + **docker_endpoint**: Docker
  + **isulad_endpoint**: Isulad
  + **podman_endpoint**: Podman

* `runtime_path` - The runtime path.

<a name="schedule_info_struct"></a>
The `schedule_info` block supports:

* `node_selector` - The node selector.

* `pod_tolerances` - The pod tolerance.
