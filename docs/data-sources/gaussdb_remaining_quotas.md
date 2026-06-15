---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_remaining_quotas"
description: |-
  Use this data source to query the remaining quotas of GaussDB enterprise projects within HuaweiCloud.
---

# huaweicloud_gaussdb_remaining_quotas

Use this data source to query the remaining quotas of GaussDB enterprise projects within HuaweiCloud.

## Example Usage

```hcl
variable "eps_id" {}

data "huaweicloud_gaussdb_remaining_quotas" "test" {
  eps_tags = [var.eps_id]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `eps_tags` - (Required, List) Specifies the list of enterprise project IDs to query the remaining quotas.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is a UUID.

* `eps_remaining_quotas` - The list of remaining quotas for the enterprise projects.
  The [eps_remaining_quotas](#eps_remaining_quotas_struct) structure is documented below.

<a name="eps_remaining_quotas_struct"></a>
The `eps_remaining_quotas` block supports:

* `eps_tag` - The enterprise project ID.

* `instance_eps_quota` - The instance quota. The value **-1** indicates unlimited quota.

* `cpu_eps_quota` - The CPU quota. The value **-1** indicates unlimited quota.

* `mem_eps_quota` - The memory quota. The value **-1** indicates unlimited quota.

* `volume_eps_quota` - The storage quota. The value **-1** indicates unlimited quota.

* `instance_eps_remaining_quota` - The remaining instance quota.

* `cpu_eps_remaining_quota` - The remaining CPU quota.

* `mem_eps_remaining_quota` - The remaining memory quota.

* `volume_eps_remaining_quota` - The remaining storage quota.
