---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_application_authorizations"
description: |-
  Use this data source to get the list of the application authorizations within HuaweiCloud.
---

# huaweicloud_workspace_application_authorizations

Use this data source to get the list of the application authorizations within HuaweiCloud.

## Example Usage

### Query all application authorizations

```hcl
variable "application_id" {}

data "huaweicloud_workspace_application_authorizations" "test" {
  app_id = var.application_id
}
```

### Filter authorizations by name and target type

```hcl
variable "application_id" {}

data "huaweicloud_workspace_application_authorizations" "test" {
  app_id      = var.application_id
  name        = "user_name"
  target_type = "SIMPLE"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the application authorizations are located.  
  If omitted, the provider-level region will be used.

* `app_id` - (Required, String) Specifies the ID of the application.

* `name` - (Optional, String) Specifies the username or user group name.  
  Fuzzy search is supported.

* `target_type` - (Optional, String) Specifies the type of the target.  
  The valid values are as follows:
  + **SIMPLE**: Normal user.
  + **USER_GROUP**: User group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `authorizations` - The list of application authorizations that matched filter parameters.  
  The [authorizations](#workspace_application_authorizations_attr) structure is documented below.

<a name="workspace_application_authorizations_attr"></a>
The `authorizations` block supports:

* `account_type` - The account type.
  + **SIMPLE**: Normal user.
  + **USER_GROUP**: User group.

* `account` - The account information.

* `domain` - The domain name.

* `platform_type` - The platform type.
  + **AD**: AD domain.
  + **LOCAL**: LiteAs.
