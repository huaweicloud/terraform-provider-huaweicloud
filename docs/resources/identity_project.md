---
subcategory: "Identity and Access Management (IAM)"
---

# huaweicloud\_identity\_project

Manages a Project resource within HuaweiCloud Identity And Access 
Management service. This is an alternative to `huaweicloud_identity_project_v3`

Note: You _must_ have security admin privileges in your HuaweiCloud 
cloud to use this resource.

## Example Usage

```hcl
resource "huaweicloud_identity_project" "project_1" {
  name        = "cn-north1_project1"
  description = "This is a test project"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) The name of the project. it must start with 
    ID of an existing region_ and be less than or equal to 64 characters.
    Example: eu-de_project1.

* `description` - (Optional, String) A description of the project.

* `domain_id` - (Optional, String) The domain this project belongs to. Changing this
    creates a new Project.

* `parent_id` - (Optional, String) The parent of this project. Changing this creates
    a new Project.

## Attributes Reference

The following attributes are exported:

## Import

Projects can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_identity_project_v3.project_1 89c60255-9bd6-460c-822a-e2b959ede9d2
```
