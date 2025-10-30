---
subcategory: "Application Operations Management (AOM 2.0)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_prom_instances"
description: |-
  Use this data source to get the list of AOM prometheus instances.
---

# huaweicloud_aom_prom_instances

Use this data source to get the list of AOM prometheus instances.

## Example Usage

```hcl
data "huaweicloud_aom_prom_instances" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the prometheus instance belongs.
  If specifies it as **all_granted_eps**, means to query instances in all enterprise projects.

* `prom_id` - (Optional, String) Specifies the prometheus instance ID.
  If both **prom_id** and **prom_type** exist, only **prom_id** takes effect.

* `prom_type` - (Optional, String) Specifies the prometheus instance type. Valid values are **default**, **ECS**,
  **VPC**, **CCE**, **REMOTE_WRITE**, **KUBERNETES**, **CLOUD_SERVICE** and **ACROSS_ACCOUNT**.

* `cce_cluster_enable` - (Optional, String) Specifies whether to enable a CCE cluster.
  Valid values are **true** and **false**.

* `prom_status` - (Optional, String) Specifies the prometheus instance status.
  Valid values are **DELETED**, **NORMAL** and **ALL**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - Indicates the prometheus instances list.

  The [instances](#instances_struct) structure is documented below.

<a name="instances_struct"></a>
The `instances` block supports:

* `id` - Indicates the prometheus instance ID.

* `prom_name` - Indicates the prometheus instance name.

* `prom_type` - Indicates the prometheus instance type.

* `prom_version` - Indicates the prometheus instance version.

* `enterprise_project_id` - Indicates the enterprise project ID to which the prometheus instance belongs.

* `prom_http_api_endpoint` - Indicates the URL for calling the prometheus instance.

* `remote_read_url` - Indicates the remote read address of the prometheus instance.

* `remote_write_url` - Indicates the remote write address of the prometheus instance.

* `prom_status` - Indicates the prometheus instance status.

* `is_deleted_tag` - Indicates whether the prometheus is deleted.

* `created_at` - Indicates the create time of the prometheus instance.

* `updated_at` - Indicates the update time of the prometheus instance.

* `deleted_at` - Indicates the delete time of the prometheus instance.
