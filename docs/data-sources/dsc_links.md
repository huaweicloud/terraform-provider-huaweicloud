---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_links"
description: |-
  Use this data source to get the data link information within HuaweiCloud.
---

# huaweicloud_dsc_links

Use this data source to get the data link information within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_links" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `db_name` - (Optional, String) Specifies the database name used to filter links.

* `dbss_instance_id` - (Optional, String) Specifies the DBSS instance ID used to filter links.

* `ecs_name` - (Optional, String) Specifies the ECS instance name used to filter links.

* `egress_type` - (Optional, String) Specifies the egress type used to filter links.

* `end_time` - (Optional, Int) Specifies the end time used to filter links.

* `internet_ip` - (Optional, String) Specifies the internet IP address used to filter links.

* `labels` - (Optional, List) Specifies the label list used to filter links.

* `oem_instance_id` - (Optional, String) Specifies the ADG instance ID used to filter links.

* `sensitive_level` - (Optional, String) Specifies the sensitive level used to filter links.

* `start_time` - (Optional, Int) Specifies the start time used to filter links.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `links` - The link information list.

  The [links](#links_struct) structure is documented below.

* `nodes` - The node information list.

  The [nodes](#nodes_struct) structure is documented below.

<a name="links_struct"></a>
The `links` block supports:

* `access_times` - The number of access times.

* `encrypt_status` - The encryption status.

* `source_node_id` - The source node ID.

* `ssl_access_times` - The number of SSL access times.

* `target_node_id` - The target node ID.

<a name="nodes_struct"></a>
The `nodes` block supports:

* `encrypt_status` - The encryption status.

* `fixed_ip` - The fixed IP address.

* `floating_ip` - The floating IP address.

* `id` - The node ID.

* `name` - The node name.

* `node_type` - The node type.

* `sensitive_infos` - The sensitive information list.

  The [sensitive_infos](#sensitive_infos_struct) structure is documented below.

<a name="sensitive_infos_struct"></a>
The `sensitive_infos` block supports:

* `label_name` - The label name.

* `sensitive_level` - The sensitive level.
