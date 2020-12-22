---
subcategory: "Intelligent EdgeCloud (IEC)"
---

# huaweicloud\_iec\_keypair

Manages a keypair resource within HuaweiCloud IEC.

## Example Usage

```hcl
resource "huaweicloud_iec_keypair" "test_keypair" {
  name       = "iec-keypair-test"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDy+49hbB9Ni2SttHcbJU+ngQXUhiGDVsflp2g5A3tPrBXq46kmm/nZv9JQqxlRzqtFi9eTI7OBvn2A34Y+KCfiIQwtgZQ9LF5ROKYsGkS2o9ewsX8Hghx1r0u5G3wvcwZWNctgEOapXMD0JEJZdNHCDSK8yr+btR4R8Ypg0uN+Zp0SyYX1iLif7saiBjz0zmRMmw5ctAskQZmCf/W5v/VH60fYPrBU8lJq5Pu+eizhou7nFFDxXofr2ySF8k/yuA9OnJdVF9Fbf85Z59CWNZBvcTMaAH2ALXFzPCFyCncTJtc/OVMRcxjUWU1dkBhOGQ/UnhHKcflmrtQn04eO8xDr root@terra-dev"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) A unique name for the keypair. Changing this parameter creates a new keypair resource.

* `public_key` - (Optional, String, ForceNew) A pregenerated OpenSSH-formatted public key. Changing this parameter creates a new keypair resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.

* `fingerprint` - Specifies a resource ID in UUID format.

## Import

Keypairs can be imported using the `name`, e.g.

```
$ terraform import huaweicloud_iec_keypair.test_keypair iec-keypair-test

```
