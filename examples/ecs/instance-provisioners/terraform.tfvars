vpc_name            = "tf_test-vpc"
subnet_name         = "tf_test-subnet"
security_group_name = "tf_test-security-group"
keypair_name        = "tf_test-keypair"
private_key_path    = "./id_rsa"
bandwidth_name      = "tf_test_for_instance"
instance_name       = "tf_test_instance_provisioner"
instance_user_data  = <<EOF
#!/bin/bash
echo "Hello, World!" > /home/test.txt
EOF

instance_remote_exec_inline = [
  "cat /home/test.txt"
]
