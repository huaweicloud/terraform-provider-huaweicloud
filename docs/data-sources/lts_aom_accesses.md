---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_aom_accesses"
description: |-
  Use this data source to get the list of AOM accesses.
---

# huaweicloud_lts_aom_accesses

Use this data source to get the list of AOM accesses.

## Example Usage

```hcl
data "huaweicloud_lts_aom_accesses" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `log_group_name` - (Optional, String) Specifies the ID of the log group name to be queried.

* `log_stream_name` - (Optional, String) Specifies the log stream name to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `accesses` - All AOM access rules that match the filter parameters.
  The [accesses](#accesses_struct) structure is documented below.

<a name="accesses_struct"></a>
The `accesses` block supports:

* `access_rules` - The AOM access log details.
  The [access_rules](#access_rules_struct) structure is documented below.

* `id` - The ID of the AOM access rule.

* `name` - The name of the AOM access rule.

* `cluster_id` - The cluster ID corresponding to the AOM access rule.

* `cluster_name` - The cluster name corresponding to the AOM access rule.

* `namespace` - The namespace corresponding to the AOM access rule.

* `workloads` - The list of the workloads corresponding to AOM access rule.

* `container_name` - The name of the container corresponding to AOM access rule.

<a name="access_rules_struct"></a>
The `access_rules` block supports:

* `file_name` - The name of the log path.

* `log_group_id` - The ID of the log group.

* `log_group_name` - The name of the log group.

* `log_stream_id` - The ID of the log stream.

* `log_stream_name` - The name of the stream.
