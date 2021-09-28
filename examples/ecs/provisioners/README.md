# Using Provisioner over SSH to an ECS instance

This example provisions an ECS instance with a Public IP address and runs a `remote-exec` provisioner over SSH.

## Using private key to remote ECS instance

For `huaweicloud_compute_instance` resoruce, if `user_data` specified, `admin_pass` setting will not take effect.
Therefore, we need to use a key-pair to log in to the ECS instance.
For this example, we create a key-pair resource and specify the private key path on local host.
For how to create a key-pair, please refer to HuaweiCloud
[documentation](https://support.huaweicloud.com/intl/en-us/usermanual-ecs/en-us_topic_0014250631.html).

Also, you can create a key-pair on the console and use HCL command
`terraform import huaweicloud_compute_keypair.default {your keypair name}` to import the key-pair into terraform state.
The [`remote-exec` provisioner](https://www.terraform.io/docs/provisioners/remote-exec.html) provide the `private_key`
to log in to ECS instance and apply some commands.
The configuration of the `private_key` is as below:

```hcl
resource "null_resource" "provision" {
  ...

  provisioner "remote-exec" {
    connection {
      ...
      private_key = "your private key"
    }

    inline = [
      ...
    ]
  }
}
```

## Using password to remote ECS instance

If you just want to log in to an ECS instance and the instance is constructed with a password, please use the following
configuration:

```hcl
resource "null_resource" "provision" {
  ...

  provisioner "remote-exec" {
    connection {
      ...
      password = "your password"
    }

    inline = [
      ...
    ]
  }
}
```

## Usage

```shell
terraform init
terraform plan
terraform apply
terraform destroy
```

The creation of the ECS instance takes about few minutes. After the creation is successful, the provisioner starts to run.

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.26.0 |
