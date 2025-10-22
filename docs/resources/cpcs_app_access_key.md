---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cpcs_app_access_key"
description: |-
  Manages a CPCS application access key resource within HuaweiCloud.
---

# huaweicloud_cpcs_app_access_key

Manages a CPCS application access key resource within HuaweiCloud.

-> Currently, this resource is valid only in cn-north-9 region.

## Example Usage

```hcl
variable "app_id" {}
variable "key_name" {}

resource "huaweicloud_cpcs_app_access_key" "test" {
  app_id   = var.app_id
  key_name = var.key_name
  status   = "enable"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `app_id` - (Required, String, NonUpdatable) Specifies the application ID to which the access key belongs.

* `key_name` - (Required, String, NonUpdatable) Specifies the access key name. The name must be unique within the
  application.

* `access_key` - (Optional, String, Sensitive, NonUpdatable) Specifies the access key AK. If omitted, the system will
  automatically generate it.

* `secret_key` - (Optional, String, Sensitive, NonUpdatable) Specifies the access key SK. If omitted, the system will
  automatically generate it.

* `status` - (Optional, String) Specifies the status of the access key. Valid values are **enable** and **disable**.
  Defaults to **enable**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (access key ID).

* `app_name` - The name of the application to which the access key belongs.

* `create_time` - The creation time of the access key, UNIX timestamp in milliseconds.

* `download_time` - The time when the access key was downloaded, UNIX timestamp in milliseconds.

* `is_downloaded` - Whether the access key has been downloaded.

* `is_imported` - Whether the access key is imported.

## Import

The CPCS application access key resource can be imported using the `app_id` and `key_name`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_cpcs_app_access_key.test <app_id>/<key_name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `secret_key`.
It is generally recommended running `terraform plan` after importing a cluster.
You can then decide if changes should be applied to the cluster, or the resource definition
should be updated to align with the cluster. Also you can ignore changes as below.

```hcl
resource "huaweicloud_cpcs_app_access_key" "test" {
    ...

  lifecycle {
    ignore_changes = [
      secret_key,
    ]
  }
}
```
