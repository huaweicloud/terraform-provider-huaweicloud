---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_ad_ous"
description: |-
  Use this data source to get organizational units (OUs) list in Active Directory (AD) within HuaweiCloud.
---

# huaweicloud_workspace_ad_ous

Use this data source to get organizational units (OUs) list in Active Directory (AD) within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_workspace_ad_ous" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the OUs are located.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `ous` - The list of OUs.  
  The [ous](#workspace_ad_ous) structure is documented below.

<a name="workspace_ad_ous"></a>
The `ous` block supports:

* `name` - The name of the OU.

* `ou_dn` - The distinguished name (DN) of the OU.
