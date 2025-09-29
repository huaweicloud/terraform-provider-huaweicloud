---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_cicd_configuration"
description: |-
  Manages an HSS CiCd configuration resource within HuaweiCloud.
---

# huaweicloud_hss_cicd_configuration

Manages an HSS CiCd configuration resource within HuaweiCloud.

## Example Usage

```hcl
variable "cicd_name" {}

resource "huaweicloud_hss_cicd_configuration" "test" {
  cicd_name = var.cicd_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `cicd_name` - (Required, String, NonUpdatable) Specifies the name of the CiCd configuration.

* `vulnerability_whitelist` - (Optional, List) Specifies the vulnerability whitelist list.

* `vulnerability_blocklist` - (Optional, List) Specifies the vulnerability blacklist list.

* `image_whitelist` - (Optional, List) Specifies the whitelist of images

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `associated_images_num` - The number of associated mirror scans.

## Import

The HSS CiCd configuration can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_hss_cicd_configuration.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `enterprise_project_id`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition
should be updated to align with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_hss_cicd_configuration" "test" { 
  ...

  lifecycle {
    ignore_changes = [
      enterprise_project_id,
    ]
  }
}
```
