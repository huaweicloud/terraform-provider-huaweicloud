---
subcategory: "EventGrid (EG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_eg_eventrouter_clusters"
description: |-
  Use this data source to query EG event router clusters within HuaweiCloud.
---

# huaweicloud_eg_eventrouter_clusters

Use this data source to query EG event router clusters within HuaweiCloud.

## Example Usage

```hcl
variable "fuzzy_match_keywrod" {}

data "huaweicloud_eg_eventrouter_clusters" "test" {
  name = var.fuzzy_match_keywrod
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the event router clusters are located.  
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the event router cluster to be queried for fuzzy matching.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `clusters` - The list of the event router clusters that matched filter parameters.  
  The [clusters](#eg_eventrouter_clusters_attr) structure is documented below.

<a name="eg_eventrouter_clusters_attr"></a>
The `clusters` block supports:

* `id` - The ID of the event router cluster.

* `name` - The name of the event router cluster.

* `description` - The description of the event router cluster.

* `source_type` - The source type of the event router cluster.

* `sink_type` - The sink type of the event router cluster.

* `vpc_id` - The VPC ID to which the event router cluster belongs.

* `subnet_id` - The subnet ID to which the event router cluster belongs.

* `availability_zones` - The availability zone names of the event router cluster.

* `flavor` - The flavor of the event router cluster.

* `charging_mode` - The charging mode of the event router cluster.

* `status` - The status of the event router cluster.

* `job_count` - The number of jobs running in the event router cluster.

* `err_code` - The error code of the event router cluster.

* `err_message` - The error message of the event router cluster.

* `public_access_enabled` - Whether public access is enabled for the event router cluster.

* `nat_id` - The NAT gateway ID of the event router cluster.

* `eip_id` - The EIP ID of the event router cluster.

* `created_at` - The creation time of the event router cluster, in RFC3339 format.

* `updated_at` - The latest update time of the event router cluster, in RFC3339 format.
