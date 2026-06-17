#!/bin/bash

all_files=$(git diff $1 $2 --name-only huaweicloud | grep -v "_test.go")

failed=0
total=0

coveraged=0
cur_time=$(date +%Y%m%d%H%M)
cov_one_file="coverage_one_${cur_time}.cov"
cov_all_file="coverage_${cur_time}.cov"
cov_html="coverage_${cur_time}.html"
echo 'mode: set' > $cov_all_file

for f in $all_files; do
    extArgs=""

    path=${f%/*}
    ## Skip directories without resource files
    if [ "$path" != "huaweicloud" ] && [[ $path != huaweicloud/services/* ]]; then
        continue
    fi

    if [ "$path" != "huaweicloud" ]; then
        coveraged=1
        extArgs="-coverprofile=${cov_one_file} -coverpkg=./${path}"
        # update path to "huaweicloud/services/acceptance/xxx"
        path=${path/"services"/"services/acceptance"}
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
        TF_ACC=1 go test -v ${extArgs} "./$path" -run "${S}" -timeout 360m -parallel 4
        if [ $? -ne 0 ]; then
            failed=$(expr $failed + 1)
        fi

        if [ -f ${cov_one_file} ]; then
            tail -q -n +2 ${cov_one_file} >> $cov_all_file
            rm -f ${cov_one_file}
        fi
        echo "\`\`\`"
    else
        echo -e "\n[skipped] ./${test_file} does not exist"
    fi
done

go tool cover -html=$cov_all_file -o $cov_html
if [ $coveraged -eq 1 ]; then
    echo -e "\n### coverage of files in huaweicloud/services:\n"
    for f in $all_files; do
        path=${f%/*}
        if [[ $path == huaweicloud/services/* ]]; then
            grep "$f" $cov_html | awk -F "<|>" '{printf "- %s\n", $3}'
        fi
    done
fi

echo -e "\n### [summary] $failed failed in $total resource acceptance tests"
