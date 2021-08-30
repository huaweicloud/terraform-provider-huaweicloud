---
subcategory: "Elastic Cloud Server (ECS)"
---

# huaweicloud_compute_keypair

Manages a keypair resource within HuaweiCloud. This is an alternative to `huaweicloud_compute_keypair_v2`

## Example Usage

```hcl
resource "huaweicloud_compute_keypair" "test-keypair" {
  name       = "my-keypair"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAlJq5Pu+eizhou7nFFDxXofr2ySF8k/yuA9OnJdVF9Fbf85Z59CWNZBvcAT... root@terra-dev"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the keypair resource. If omitted, the
  provider-level region will be used. Changing this creates a new keypair resource.

* `name` - (Required, String, ForceNew) Specifies a unique name for the keypair. Changing this creates a new keypair.

* `public_key` - (Required, String, ForceNew) Specifies the imported OpenSSH-formatted public key. Changing this creates
  a new keypair.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

## Import

Keypairs can be imported using the `name`, e.g.

```
$ terraform import huaweicloud_compute_keypair.my-keypair test-keypair
```
