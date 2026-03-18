#!/usr/bin/env bash

# Check gofmt
echo "==> Checking for unchecked errors..."

if ! which errcheck > /dev/null; then
    echo "==> Installing errcheck..."
    go get -u github.com/kisielk/errcheck
fi

# Parse command line arguments
TARGET_FILE=""

while [[ $# -gt 0 ]]; do
    case $1 in
        -f | --file)
            if [[ $# -gt 1 ]]; then
                TARGET_FILE=$2
                shift 2
            else
                echo "Error: -f/--file requires a file path argument"
                echo "Usage: $0 [-f|--file <file_path>]"
                exit 1
            fi
            ;;
        -h | --help)
            echo "Usage: $0 [-f|--file <file_path>]"
            echo ""
            echo "Options:"
            echo "  -f, --file    Check specific file (schema or test file)"
            echo "  -h, --help    Show this help message"
            echo ""
            echo "Examples:"
            echo "  $0                                    # Check all files"
            echo "  $0 -f huaweicloud/services/dws/data_source_huaweicloud_dws_host_nets.go"
            echo "  $0 --file huaweicloud/services/acceptance/dws/data_source_huaweicloud_dws_host_nets_test.go"
            exit 0
            ;;
        *)
            echo "Error: Unknown option $1"
            echo "Usage: $0 [-f|--file <file_path>]"
            exit 1
            ;;
    esac
done

if [[ -n "$TARGET_FILE" ]]; then
    echo "==> Checking specific file: $TARGET_FILE"
    err_files=$(errcheck -ignoretests \
                         -ignore 'github.com/hashicorp/terraform/helper/schema:Set' \
                         -ignore 'bytes:.*' \
                         -ignore 'io:Close|Write' \
                         "$TARGET_FILE")
else
    echo "==> Checking all files..."
    err_files=$(errcheck -ignoretests \
                         -ignore 'github.com/hashicorp/terraform/helper/schema:Set' \
                         -ignore 'bytes:.*' \
                         -ignore 'io:Close|Write' \
                         $(go list ./...| grep -v /vendor/))
fi

if [[ -n ${err_files} ]]; then
    echo 'Unchecked errors found in the following places:'
    echo "${err_files}"
    echo "Please handle returned errors. You can check directly with \`make errcheck\`"
    exit 1
fi

echo "Check completed!"
exit 0
