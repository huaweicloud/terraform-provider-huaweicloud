---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_asset_authorization"
description: |-
  Manages a DSC asset authorization resource within HuaweiCloud.
---

# huaweicloud_dsc_asset_authorization

Manages a DSC asset authorization resource within HuaweiCloud.

-> Deleting this resource will not cancel the asset authorization, but will only remove the resource information from
  the tf state file.

## Example Usage

```hcl
resource "huaweicloud_dsc_asset_authorization" "test" {
  type                 = "LTS"
  authorization_status = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `type` - (Required, String, NonUpdatable) Specifies the asset type to be authorized.
  Valid values are **DASHBOARD**, **OBS**, **DATABASE**, **BIGDATA**, **MRS**, **LTS** and **ALL**.

* `authorization_status` - (Required, Bool) Specifies the authorization status.
  The value can be **true** (authorize) or **false** (cancel authorization).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (same as `type`).

* `bigdata_authorization` - Whether the big data asset is authorized.

* `dashboard_authorization` - Whether the dashboard is authorized.

* `database_authorization` - Whether the database asset is authorized.

* `lts_authorization` - Whether the LTS log service is authorized.

* `mrs_authorization` - Whether the MRS big data cluster is authorized.

* `obs_authorization` - Whether the OBS object storage is authorized.

## Import

The DSC asset authorization resource can be imported using the `type`, e.g.

```bash
$ terraform import huaweicloud_dsc_asset_authorization.test <type>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `authorization_status`. It is generally recommended running
`terraform plan` after importing the resource. You can then decide if changes should be applied to the resource, or the
resource definition should be updated to align with the cloud. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_dsc_asset_authorization" "test" {
  ...

  lifecycle {
    ignore_changes = [
      authorization_status,
    ]
  }
}
```
