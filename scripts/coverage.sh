#!/bin/bash

TEST_TARGET_ACCEPTANCE_PKG_NAME=""  
TEST_TARGET_SERVICE_PKG_NAME=""  
TEST_FUNC_PREFIX="TestAcc"  
TEST_PARALLEL_NUM=10  
DESIRED_PATH=""  

# The handle_o function is used to identify the -o option and assign the following parameters to the variables  
# TEST_TARGET_SERVICE_PKG_NAME and TEST_TARGET_SERVICE_PKG_NAME. The two are consistent by default.  
# If the service package name is inconsistent, it can be set through the -s option.  
handle_o() {  
    TEST_TARGET_ACCEPTANCE_PKG_NAME=$1  
    if [ -z "$TEST_TARGET_SERVICE_PKG_NAME" ]; then  
        TEST_TARGET_SERVICE_PKG_NAME=$TEST_TARGET_ACCEPTANCE_PKG_NAME 
    fi 
}  

# The function handle_f is used to recognize the -f option and assign the following parameters to the variable  
# TEST_FUNC_PREFIX.  
handle_f() {  
    TEST_FUNC_PREFIX=$1  
}  

# The handle_n function recognizes the -n option and assigns the following parameter to the variable  
# TEST_PARALLEL_NUM.  
handle_n() {  
    TEST_PARALLEL_NUM=$1  
}  

# The handle_s function recognizes the -s option and assigns the following parameter to the variable  
# TEST_TARGET_SERVICE_PKG_NAME.  
handle_s() {  
    TEST_TARGET_SERVICE_PKG_NAME=$1  
}  

# handle_options is used to process options and their corresponding parameters.  
# Check the parameters and take appropriate action.  
while [[ $# -gt 0 ]]; do  
    case $1 in  
        -o)    
            if [[ $# -gt 1 ]]; then  
                handle_o $2  
                shift  
            else  
                echo "Warning: -o needs to be followed by an expected acceptance package name."  
                exit 1  
            fi  
            ;; 
        -f)  
            if [[ $# -gt 1 ]]; then  
                handle_f $2  
                shift  
            else  
                echo "Warning: -f needs to be followed by a function name input of the acceptance test."  
                exit 1  
            fi  
            ;;  
        -n)    
            if [[ $# -gt 1 ]]; then  
                handle_n $2  
                shift  
            else  
                echo "Warning: -n needs to be followed by a parallel number input."  
                exit 1  
            fi  
            ;;   
        -s)   
            if [[ $# -gt 1 ]]; then  
                handle_s $2  
                shift  
            else  
                echo "Warning: -s needs to be followed by an expected service package name."  
                exit 1  
            fi  
            ;;  
        -h | --help)  
            echo "Usage: $0 [options] [args] ([options] [args]...)  

The available options for execution are listed below.  

  -o            The name of the acceptance package to be tested. If omitted, all acceptance packages will be tested.  
  -f            Uses the custom function name prefix, defaults to 'TestAcc'.  
  -n            The custom parallel number, defaults to 5.  
  -s            The service name corresponding to the acceptance package where the functions are stored. If omitted,  
                the value will be same as the acceptance package name.  
  -h | --help   Show the user guide.  
"  
            exit 0  
            ;;  
        *)   
            echo "Invalid option: $1"  
            exit 1  
            ;;  
    esac  
    shift  
done  

# Iterate through all subdirectory names of a specific directory and store them into a global array.  
get_subdir_names() {  
    local directory_path="$1"  
    local -n subdir_names_ref=$2  

    for item in "$directory_path"/*; do  
        # Check whether the object is a directory.  
        if [ -d "$item" ]; then  
            local subdir_name=$(basename "$item")  
            subdir_names_ref+=("$subdir_name")  
        fi  
    done  
}  

calculate_all_services_coverage() {
    RELATIVE_PATH=$(realpath --relative-to="$(pwd)" "${DESIRED_PATH}/huaweicloud/services")  

    subdir_names=()  

    # Record the name list of the sub directories.  
    get_subdir_names "${RELATIVE_PATH}/acceptance" subdir_names  

    echo "Prepare to execute the acceptance tests for all services"

    for name in "${subdir_names[@]}"; do  
        echo "Prepare to calculate the coverage the following command:
        TF_ACC=1 go test \"./${RELATIVE_PATH}/acceptance/${name}\" -v \
-coverprofile=\"./${RELATIVE_PATH}/acceptance/${name}/${name}_coverage.cov\" \
-coverpkg=\"./${RELATIVE_PATH}/${name}\" -run $TEST_FUNC_PREFIX -timeout 360m -parallel $TEST_PARALLEL_NUM"  

        TF_ACC=1 go test "./${RELATIVE_PATH}/acceptance/${name}" -v -coverprofile="./${RELATIVE_PATH}/acceptance/${name}_coverage.cov" \
            -coverpkg="./${RELATIVE_PATH}/${name}" -run $TEST_FUNC_PREFIX -timeout 360m -parallel $TEST_PARALLEL_NUM  
        go tool cover -html="./${RELATIVE_PATH}/acceptance/${name}/${name}_coverage.cov" -o \
"./${RELATIVE_PATH}/acceptance/${name}/${name}_coverage.html"  
    done  
}  

calculate_special_service_coverage() {
    RELATIVE_PATH=$(realpath --relative-to="$(pwd)" "$DESIRED_PATH/huaweicloud/services")  

    echo "Prepare to calculate the coverage the following command:
TF_ACC=1 go test \"./${RELATIVE_PATH}/acceptance/${TEST_TARGET_ACCEPTANCE_PKG_NAME}\" -v \
-coverprofile=\"./${RELATIVE_PATH}/acceptance/${TEST_TARGET_ACCEPTANCE_PKG_NAME}/${TEST_TARGET_ACCEPTANCE_PKG_NAME}_coverage.cov\" \
-coverpkg=\"./${RELATIVE_PATH}/${TEST_TARGET_SERVICE_PKG_NAME}\" -run $TEST_FUNC_PREFIX -timeout 360m -parallel $TEST_PARALLEL_NUM"  

    TF_ACC=1 go test "./${RELATIVE_PATH}/acceptance/${TEST_TARGET_ACCEPTANCE_PKG_NAME}" -v \
        -coverprofile="./${RELATIVE_PATH}/acceptance/${TEST_TARGET_ACCEPTANCE_PKG_NAME}/${TEST_TARGET_ACCEPTANCE_PKG_NAME}_coverage.cov" \
        -coverpkg="./${RELATIVE_PATH}/${TEST_TARGET_SERVICE_PKG_NAME}" -run $TEST_FUNC_PREFIX -timeout 360m -parallel $TEST_PARALLEL_NUM
    go tool cover -html="./${RELATIVE_PATH}/acceptance/${TEST_TARGET_ACCEPTANCE_PKG_NAME}/${TEST_TARGET_ACCEPTANCE_PKG_NAME}_coverage.cov" \
        -o "./${RELATIVE_PATH}/acceptance/${TEST_TARGET_ACCEPTANCE_PKG_NAME}/${TEST_TARGET_ACCEPTANCE_PKG_NAME}_coverage.html"  
}  

CURRENT_PATH=$(pwd)  
# Make sure the current path is in the terraform-provider-huaweicloud directory or its subdirectory.  
if [[ "$CURRENT_PATH" != *"terraform-provider-huaweicloud"* ]]; then  
    echo "Warning: Your current executant is not in the terraform-provider-huaweicloud directory or its subdirectories.
Please use the 'pwd' command to check it or use the 'cd' command to switch."  
exit 1  
fi  

# Using ##* model to match the current path and check whether in the desired path (terraform-provider-huaweicloud  
# directory or its subdirectory), if so, store the expected path to a environment variable.  
DESIRED_PATH="${CURRENT_PATH%%/terraform-provider-huaweicloud/*}"  

SUFFIX="/terraform-provider-huaweicloud"  
  
# Check if you are currently in the root directory of the repository
if [[ ! "$DESIRED_PATH" == *"$SUFFIX"* ]]; then  
    DESIRED_PATH="$DESIRED_PATH$SUFFIX"  
fi 

if [ -n "$TEST_TARGET_ACCEPTANCE_PKG_NAME" ]; then  
    # Execute coverage tests for specific services based on the parameters entered (-o) by the user.  
    calculate_special_service_coverage  
    exit 1  
fi  

# Execute coverage tests for all services.  
calculate_all_services_coverage  
