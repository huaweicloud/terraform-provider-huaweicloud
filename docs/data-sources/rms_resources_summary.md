---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_resources_summary"
description: |-
  Use this data source to get the list of RMS resources summary.
---

# huaweicloud_rms_resources_summary

Use this data source to get the list of RMS resources summary.

## Example Usage

```hcl
data "huaweicloud_rms_resources_summary" "test" {
  type = "vpc.vpcs"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, String) Specifies the resource name.

* `type` - (Optional, String) Specifies the resource type. For example, **vpc.vpcs** and **rds.instances**.

* `region_id` - (Optional, String) Specifies the region to which the resource belongs.

* `enterprise_project_id` - (Optional, String) Specifies the ID of enterprise project to which the resource belongs.

* `project_id` - (Optional, String) Specifies the ID of project to which the resource belongs.

* `tracked` - (Optional, Bool) Specifies whether the resource is tracked.

* `resource_deleted` - (Optional, Bool) Specifies whether the query the deleted resources.

* `tags` - (Optional, Map) Specifies the tags to filter the resources.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resources_summary` - The service details list.

  The [resources_summary](#resources_summary) structure is documented below.

<a name="resources_summary"></a>
The `resources_summary` block supports:

* `service` - The service name.

* `types` - The resource type list.
  The [types](#types) structure is documented below.

<a name="types"></a>
The `types` block supports:

* `type` - The resource type.

* `regions` - The list of supported regions.
  The [regions](#regions) structure is documented below.

<a name="regions"></a>
The `regions` block supports:

* `region` - The region name.

* `count` - The number of resource in this region.
