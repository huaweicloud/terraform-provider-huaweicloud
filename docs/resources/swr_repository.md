---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_repository"
description: ""
---

# huaweicloud_swr_repository

Manages a SWR repository resource within HuaweiCloud.

## Example Usage

```hcl
variable "organization_name" {} 

resource "huaweicloud_swr_repository" "test" {
  organization = var.organization_name
  name         = "%s"
  description  = "Test repository"
  category     = "linux"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `organization` - (Required, String, ForceNew) Specifies the name of the organization (namespace) the repository
  belongs. Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the repository.  
  The valid length is limited from `1` to `128`, only lowercase letters, digits, periods (.), hyphens (-) and
  underscores (_) are allowed. Periods, underscores, and hyphens cannot be placed next to each other.
  A maximum of two consecutive underscores are allowed. Changing this creates a new resource.

* `is_public` - (Optional, Bool) Specifies whether the repository is public. Default is false.
  + `true` - Indicates the repository is public.
  + `false` - Indicates the repository is private.

* `description` - (Optional, String) Specifies the description of the repository.

* `category` - (Optional, String) Specifies the category of the repository.
  The value can be `app_server`, `linux`, `framework_app`, `database`, `lang`, `other`, `windows`, `arm`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the repository. The value is the name of the repository.

* `repository_id` - Numeric ID of the repository

* `path` - Image address for docker pull.

* `internal_path` - Intra-cluster image address for docker pull.

* `num_images` - Number of image tags in a repository.

* `size` - Repository size.

## Import

Repository can be imported using the organization name and repository name separated by a slash or a comma, e.g.:

Only when repository name is with no slashes, can use a slash to separate the organization name and repository name.

```bash
$ terraform import huaweicloud_swr_repository.test org-name/repo-name
```

Using comma to separate the organization name and repository name is available for repository name with slashes or not.

```bash
$ terraform import huaweicloud_swr_repository.test org-name,repo-name
```
