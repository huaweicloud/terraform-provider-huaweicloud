---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_architecture_subject"
description: ""
---

# huaweicloud_dataarts_architecture_subject

Manages DataArts Architecture subject resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "name" {}
variable "code" {}
variable "owner" {}
variable "subject_id" {}

resource "huaweicloud_dataarts_architecture_subject" "test-L1" {
  workspace_id = var.workspace_id
  name         = var.name
  code         = var.code
  owner        = var.owner
  level        = 1
}

resource "huaweicloud_dataarts_architecture_subject" "test-L2" {
  workspace_id = var.workspace_id
  name         = var.name
  code         = var.code
  owner        = var.owner
  parent_id    = var.subject_id
  level        = 2
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to manage the subject.
  If omitted, the provider-level region will be used. Changing this creates a new subject.

* `workspace_id` - (Required, String, ForceNew) Specifies the workspace ID which the subject in.
  Changing this creates a new subject.

* `code` - (Required, String, ForceNew) Specifies the subject code.

* `name` - (Required, String) Specifies the subject name.

* `owner` - (Required, String) Specifies the owner of the subject.

* `level` - (Required, Int) Specifies the level of subject. The valid values are `1`, `2` and `3`.

* `parent_id` - (Optional, String) Specifies the parent ID of the subject.
  It's **Required** when you created a **L2** or **L3** subject.

* `description` - (Optional, String) Specifies the description of subject.
  It's **Required** when you created a **L3** subject.

* `department` - (Optional, String) Specifies the department of subject.
  It's **Required** when you created a **L3** subject.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `path` - The subject path. Format is `<L1_subject_name>.<L2_subject_name>.<L3_subject_name>`

* `created_at` - The creating time of the subject.

* `updated_at` - The updating time of the subject.

* `created_by` - The person creating the subject.

* `updated_by` - The person updating the subject.

* `status` - The status of the subject.

* `guid` - The globally unique ID of the subject, generating when the subject was published.

## Import

DataArts Architecture subject can be imported using `<workspace_id>/<path>`, e.g.

```sh
terraform import huaweicloud_dataarts_architecture_subject.test b606cd4a47b645108a122857204b360f/test-L1.test-L2
```
