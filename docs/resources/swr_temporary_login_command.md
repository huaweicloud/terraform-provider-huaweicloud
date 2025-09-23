---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_temporary_login_command"
description: |-
  Manages a SWR temporary login command resource within HuaweiCloud.
---

# huaweicloud_swr_temporary_login_command

Manages a SWR temporary login command resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_swr_temporary_login_command" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `enhanced` - (Optional, Bool, NonUpdatable) Specifies whether to create enhanced login command. Default to **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `auths` - Indicates the authentication information.
  The [auths](#attrblock--auths) structure is documented below.

* `x_swr_docker_login` - Indicates the docker login command.

* `x_expire_at` - Indicates the expiration time of the login command.

<a name="attrblock--auths"></a>
The `auths` block supports:

* `key` - Indicates the authentication information key.

* `auth` - Indicates the base64-encoded authentication information.
