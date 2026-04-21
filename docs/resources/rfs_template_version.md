---
subcategory: "Resource Formation (RFS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rfs_template_version"
description: |-
  Manages a RFS template version resource within HuaweiCloud.
---

# huaweicloud_rfs_template_version

Manages a RFS template version resource within HuaweiCloud.

-> Template versions are immutable and append-only. Once a version is created, it cannot be modified.
To update a template version, you must create a new version with the desired changes.

## Example Usage

```hcl
variable "template_name" {} 
variable "template_id" {} 
variable "version_description" {}
variable "template_body" {}

resource "huaweicloud_rfs_template_version" "test" {
  template_name       = var.template_name 
  template_id         = var.template_id 
  version_description = var.version_description 
  template_body       = var.template_body
}
```

## Argument Reference

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this will create a new resource.

* `template_name` - (Required, String, NonUpdatable) Specifies the template name. The name must be unique and can
  contain only letters, digits, hyphens (-), and underscores (_). This field refers to an existing template.

* `template_id` - (Optional, String, NonUpdatable) Specifies the template ID used as a query parameter when creating
  the version.

* `template_body` - (Optional, String, NonUpdatable) Specifies the Terraform template body. This field
  and `template_uri` are mutually exclusive, one of them must be specified. The template body must be a valid Terraform
  configuration in HCL format.

* `template_uri` - (Optional, String, NonUpdatable) Specifies the OBS URL where the Terraform template is stored. This
  field and `template_body` are mutually exclusive, one of them must be specified. The corresponding file should be a
  pure Terraform file or a ZIP archive:
  + For pure Terraform files, the file name must end with `.tf` or `.tf.json` and comply with HCL syntax.
  + For archives, only ZIP format is supported. The file name must end with `.zip`. After decompression, the archive
    must not contain `.tfvars` files.

* `version_description` - (Optional, String, NonUpdatable) Specifies the description of this template version.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (also version ID).

* `create_time` - The time when the template was created, in RFC3339 format (yyyy-mm-ddTHH:MM:SSZ),
  such as **1970-01-01T00:00:00Z**.

## Import

The template version can be imported using their `template_name` and `id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_rfs_template_version.test <template_name>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include `template_body`, `template_uri` and `version_description`. It is generally
recommended running `terraform plan` after importing a resource. You can then decide if changes should be applied to the
resource, or the resource definition should be updated to align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_rfs_template_version" "test" {
  ...

  lifecycle {
    ignore_changes = [
      template_body,
      template_uri,
      version_description,
    ]
  }
}
```
