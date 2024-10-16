---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_central_network_policies"
description: |-
  Use this data source to get the list of CC central network policies.
---

# huaweicloud_cc_central_network_policies

Use this data source to get the list of CC central network policies.

## Example Usage

```hcl
variable "central_network_id" {}
variable "central_network_policy_id" {}

data "huaweicloud_cc_central_network_policies" "test" {
  central_network_id = var.central_network_id
  policy_id          = var.central_network_policy_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `central_network_id` - (Required, String) Specifies the ID of central network.

* `policy_id` - (Optional, String) Specifies the ID of central network policy.

* `status` - (Optional, String) Specifies the status of central network policy.
  The valid values can be **AVAILABLE**, **CANCELING**, **APPLYING**, **FAILED** and **DELETED**.

* `is_applied` - (Optional, String) Specifies whether the central network policy is applied or not.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `central_network_policies` - The list of the central network policies.

  The [central_network_policies](#central_network_policies_struct) structure is documented below.

<a name="central_network_policies_struct"></a>
The `central_network_policies` block supports:

* `id` - The ID of the central network policies.

* `status` - The status of the central network policy.

* `central_network_id` - The ID of the central network.

* `is_applied` - Whether the policy is applied or not.

* `version` - The version of the central network policy.

* `document_template_version` - The document template version of the central network policy.

* `document` - The document of the central network policy.

  The [document](#central_network_policies_document_struct) structure is documented below.

* `created_at` - The creation time of the central network policy.
  The time is in the **yyyy-MM-ddTHH:mm:ss** format.

<a name="central_network_policies_document_struct"></a>
The `document` block supports:

* `default_plane` - The name of the default central network plane.

* `planes` - The list of the central network planes.

  The [planes](#document_planes_struct) structure is documented below.

* `er_instances` - The list of the enterprise routers instances.

  The [er_instances](#document_er_instances_struct) structure is documented below.

<a name="document_planes_struct"></a>
The `planes` block supports:

* `name` - The name of the central network plane.

* `associate_er_tables` - The list of the enterprise router tables on the central network.

  The [associate_er_tables](#planes_associate_er_tables_struct) structure is documented below.

* `exclude_er_connections` - The list of the enterprise router connections excluded from the central network policy.

  The [exclude_er_connections](#planes_exclude_er_connections_struct) structure is documented below.

<a name="planes_associate_er_tables_struct"></a>
The `associate_er_tables` block supports:

* `project_id` - The project ID of the enterprise router table on the central network.

* `region_id` - The region ID of the enterprise router table on the central network.

* `enterprise_router_table_id` - The ID of the enterprise router table on the central network.

* `enterprise_router_id` - The ID of the enterprise router table on the central network.

<a name="planes_exclude_er_connections_struct"></a>
The `exclude_er_connections` block supports:

* `exclude_er_instances` - The list of enterprise routers that will not establish a connection.
  The [exclude_er_instances](#planes_exclude_er_instances_struct) structure is the same as `er_instances`.

<a name="planes_exclude_er_instances_struct"></a>
The `exclude_er_instances` block supports:

* `project_id` - The project ID of the exclude enterprise router instance.

* `region_id` - The region ID of the exclude enterprise router instance.

* `enterprise_router_id` - The ID of the exclude enterprise router instance.

<a name="document_er_instances_struct"></a>
The `er_instances` block supports:

* `project_id` - The project ID of the enterprise router on the central network.

* `region_id` - The region ID of the enterprise router on the central network.

* `enterprise_router_id` - The ID of the enterprise router on the central network.
