---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_quotas"
description: |-
  Use this data source to get the list of enterprise project quota.
---

# huaweicloud_gaussdb_opengauss_quotas

Use this data source to get the list of enterprise project quota.

## Example Usage

```hcl
data "huaweicloud_gaussdb_opengauss_quotas" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `eps_quotas` - Indicates the list of enterprise project quota.

  The [eps_quotas](#eps_quotas_struct) structure is documented below.

<a name="eps_quotas_struct"></a>
The `eps_quotas` block supports:

* `enterprise_project_id` - Indicates the enterprise project ID.

* `enterprise_project_name` - Indicates the enterprise project name.

* `instance_eps_quota` - Indicates the instance quantity quota.

* `vcpus_eps_quota` - Indicates the vCPU quota.

* `ram_eps_quota` - Indicates the memory quota in GB.

* `volume_eps_quota` - Indicates the storage quota in GB.

* `instance_used` - Indicates the used EPS instance quota.

* `vcpus_used` - Indicates the used EPS compute quota.

* `ram_used` - Indicates the used EPS memory quota in GB.

* `volume_used` - Indicates the used EPS storage quota, in GB.
