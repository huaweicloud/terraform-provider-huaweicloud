vpc_name             = "tf_test_vpc_with_ecs_instance"
subnet_name          = "tf_test_subnet_with_ecs_instance"
security_group_names = ["tf_test_seg1_with_ecs_instance", "tf_test_seg2_with_ecs_instance"]
instance_name        = "tf_test_with_userdata"
keypair_name         = "tf_test_keypair_with_ecs_instance"
instance_user_data   = <<EOF
#!/bin/bash
echo "Hello, World!" > /home/terraform.txt
EOF
