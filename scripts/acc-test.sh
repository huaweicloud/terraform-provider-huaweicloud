#!/bin/bash

all_files=$(git diff $1 $2 --name-only huaweicloud | grep -v "_test.go")

failed=0
total=0

for f in $all_files; do

    path=${f%/*}

    ## Skip directories without resource files
    if [ "$path" != "huaweicloud" ] && [[ $path != huaweicloud/services/* ]]; then
        continue
    fi

    if [ "$path" != "huaweicloud" ]; then
        # update path to "huaweicloud/services/acceptance/xxx"
        path=${path/services/services\/acceptance}
    fi
    S=""
    org_file=${f##*/}
    test_file=$path/${org_file/%.go/_test.go}
    if [ -f "./${test_file}" ]; then
        total=$(expr $total + 1)
        cases=$(grep "^func TestAcc" ./${test_file} | awk '{print $2}' | awk -F '(' '{print $1}')
        for case in ${cases[@]}; do
            S="${S:+${S}|}${case}"
        done
        echo -e "\n\`\`\` \nrun acceptance tests of ${test_file}:"
        TF_ACC=1 go test "./$path" -v -run "${S}" -timeout 360m -parallel 4
        if [ $? -ne 0 ]; then
            failed=$(expr $failed + 1)
        fi
        echo "\`\`\`"
    else
        echo -e "\n[skipped] ./${test_file} does not exist"
    fi
done

echo -e "\n### [summary] $failed failed in $total resource acceptance tests"
