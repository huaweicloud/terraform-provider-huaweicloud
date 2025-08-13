---
subcategory: "CodeArts"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_repository"
description: ""
---

# huaweicloud_codearts_repository

Manages a CodeArts repository resource within HuaweiCloud.

## Example Usage

```hcl
variable "project_id" {}
variable "repository_name" {}
variable "repository_description" {}

resource "huaweicloud_codearts_repository" "test" {
  project_id = var.project_id // You can use project resource to generate it

  name             = var.repository_name
  description      = var.repository_description
  gitignore_id     = "Go"
  enable_readme    = 0  // Do not auto generate README.md
  visibility_level = 20 // Public read-only
  license_id       = 2  // MIT License
  import_members   = 0  // Do not import members of the project into this repository when creating
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) The repository name.
  Changing this parameter will create a new resource.

* `project_id` - (Required, String, ForceNew) The project ID for CodeArts service.
  Changing this parameter will create a new resource.

* `visibility_level` - (Optional, Int, ForceNew) The visibility level.  
  The valid values are as follows:
  + **0**: Private.
  + **20**: Public read-only.

  Defaults to `0`. Changing this parameter will create a new resource.

* `description` - (Optional, String, ForceNew) The repository description.
  Changing this parameter will create a new resource.

* `import_url` - (Optional, String, ForceNew) The HTTPS address of the template repository encrypted using Base64.
  Changing this parameter will create a new resource.

* `gitignore_id` - (Optional, String, ForceNew) The program language type for generating `.gitignore` files.
  Changing this parameter will create a new resource.

* `license_id` - (Optional, Int, ForceNew) The license ID for public repository. The valid values are as follows:
  + **1**: Apache License v2.0
  + **2**: MIT License
  + **3**: BSD 2-clause
  + **4**: BSD 3-clause
  + **5**: Eclipse Public License v1.0
  + **6**: GNU General Public License v2.0
  + **7**: GNU General Public License v3.0
  + **8**: GNU Afferent General Public License v3.0
  + **9**: GNU Lesser General Public License v2.1
  + **10**: GNU Lesser General Public License v3.0
  + **11**: Mozilla Public License v2.0
  + **12**: The Unlicense

  Defaults to `1`. Changing this parameter will create a new resource.

* `enable_readme` - (Optional, Int, ForceNew) Whether to generate the `README.md` file.  
  The valid values are as follows:
  + **0**: Disable.
  + **1**: Enable.

  Defaults to `1`. Changing this parameter will create a new resource.

* `import_members` - (Optional, Int, ForceNew) Whether to import the project members.  
  The valid values are as follows:
  + **0**: Do not import members.
  + **1**: Import members.

  Defaults to `1`. Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `https_url` - The HTTPS URL that used to the fork repository.

* `ssh_url` - The SSH URL that used to the fork repository.

* `web_url` - The web URL, accessing this URL will redirect to the repository detail page.

* `lfs_size` - The LFS capacity, in MB. If the capacity is greater than `1,024`M, the unit is GB.

* `capacity` - The total size of the repository, in MB. If the capacity is greater than `1,024`M, the unit is GB.

* `status` - The repository status.  
  The valid values are as follows:
  + **0**: Normal.
  + **3**: Frozen.
  + **4**: Closed.

* `create_at` - The creation time.

* `update_at` - The last update time.

* `repository_id` - The repository primart key ID.

## Import

The repository can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_codearts_repository.test 0ce123456a00f2591fabc00385ff1234
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `name`, `description`, `gitignore_id`, `enable_readme`, `license_id` and
`import_members`. It is generally recommended running `terraform plan` after importing the repository.
You can then decide if changes should be applied to the repository, or the resource definition should be updated to
align with the repository. Also you can ignore changes as below.

```hcl
resource "huaweicloud_codearts_repository" "test" {
  ...

  lifecycle {
    ignore_changes = [
      name, license_id,
    ]
  }
}
```
