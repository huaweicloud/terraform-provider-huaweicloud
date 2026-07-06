---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_enterprise_project_remaining_quotas"
description: |-
  Use this data source to query the remaining quotas of enterprise projects for GaussDB within HuaweiCloud.
---

# huaweicloud_gaussdb_enterprise_project_remaining_quotas

Use this data source to query the remaining quotas of enterprise projects for GaussDB within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_gaussdb_enterprise_project_remaining_quotas" "test" {
  eps_tags = ["0"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the remaining quotas.
  If omitted, the provider-level region will be used.

* `eps_tags` - (Required, List) Specifies the list of enterprise project IDs to query.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `job_id` - The job ID.

* `eps_remaining_quotas` - The list of remaining enterprise project quotas.
  The [eps_remaining_quotas](#eps_remaining_quotas) structure is documented below.

<a name="eps_remaining_quotas"></a>
The `eps_remaining_quotas` block supports:

* `eps_tag` - The enterprise project ID.

* `instance_eps_quota` - The instance quota.

* `cpu_eps_quota` - The CPU quota.

* `mem_eps_quota` - The memory quota.

* `volume_eps_quota` - The volume quota.

* `instance_eps_remaining_quota` - The remaining instance quota.

* `cpu_eps_remaining_quota` - The remaining CPU quota.

* `mem_eps_remaining_quota` - The remaining memory quota.

* `volume_eps_remaining_quota` - The remaining volume quota.
