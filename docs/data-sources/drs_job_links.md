---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_job_links"
description: |-
  Use this data source to get a list of job links for DRS.
---

# huaweicloud_drs_job_links

Use this data source to get a list of job links for DRS.

## Example Usage

```hcl
variable "job_type" {}

data "huaweicloud_drs_job_links" "test" { 
  job_type = var.job_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `job_type` - (Required, String) Specifies the DRS job type.
  The valid values are as follows:
  + **migration**: Online migration.
  + **sync**: Data synchronization.
  + **cloudDataGuard**: Disaster recovery.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `job_links` - The list of job links.

The [job_links](#job_links_struct) structure is documented below.

<a name="job_links_struct"></a>
The `job_links` block supports:

* `job_type` - The job type.

* `engine_type` - The engine type of the job.

* `source_endpoint_type` - The source endpoint type.

* `target_endpoint_type` - The target endpoint type.

* `job_direction` - The direction of data flow.

* `net_type` - The network type of the job.

* `task_types` - The list of task types.

* `cluster_modes` - The list of cluster modes.
