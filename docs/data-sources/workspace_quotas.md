---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_quotas"
description: |-
  Use this data source to get the tenant quotas of the Workspace within HuaweiCloud.
---

# huaweicloud_workspace_quotas

Use this data source to get the tenant quotas of the Workspace within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_workspace_quotas" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the quotas are located.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - The list of common quotas.  
  The [quotas](#workspace_quotas_attr) structure is documented below.

* `site_quotas` - The list of site quotas.  
  The [site_quotas](#workspace_site_quotas_attr) structure is documented below.

<a name="workspace_quotas_attr"></a>
The `quotas` block supports:

* `resources` - The list of quota resources.  
  The [resources](#quota_resources_attr) structure is documented below.

<a name="quota_resources_attr"></a>
The `resources` block supports:

* `type` - The resource type.
  + **general_instances**: Number of general desktops.
  + **Memory**: Memory capacity.
  + **cores**: Number of CPUs.
  + **volumes**: Number of disks.
  + **volume_gigabytes**: Disk capacity.
  + **gpu_instances**: Number of GPU desktops.
  + **deh**: Cloud office host.
  + **users**: Number of users.
  + **policy_groups**: Number of policy groups.
  + **Cores**: Number of CPUs (used by quota tool).

* `quota` - The quota value.

* `used` - The used quota value.

* `unit` - The quota unit.

<a name="workspace_site_quotas_attr"></a>
The `site_quotas` block supports:

* `site_id` - The site ID.

* `resources` - The list of quota resources.  
  The [resources](#quota_resources_attr) structure is documented above.
