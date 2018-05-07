#!/bin/bash

readApiToken(){
    gitlab_root_username=$1
    gitlab_root_password=$2
    token_request=$(curl --show-error --silent --request POST -H "Content-Type: application/json" \
    --data "{\"grant_type\":\"password\",\"username\":\"$gitlab_root_username\",\"password\": \"$gitlab_root_password\"}" \
    $gitlab_host/oauth/token)

    token=$(echo $token_request | jq -r '.access_token')
    if [ "$token" = "null" ]; then
        echo $token_request
        exit 1
    fi;
    echo $token
}

queryGitlabApi(){
    api_token=$1
    endpoint=$2
    curl --silent --show-error --header "Authorization: Bearer ${api_token}" "http://${gitlab_host}/api/v4/${endpoint}"
}

createGitlabUser(){
    username=$1
    name=$1
    password=$2
    email=$3
    curl --silent --show-error --header "Authorization: Bearer ${root_api_token}" \
        --request POST -H "Content-Type: application/json" \
        --data "{\"email\":\"$email\",\"password\":\"$password\",\"username\":\"$username\",\"name\":\"$name\",\"skip_confirmation\":\"true\"}" \
        "http://${gitlab_host}/api/v4/users"
}

createGitlabGroup(){
    user_token=$1
    name=$2
    path=$2
    curl --silent --show-error --header "Authorization: Bearer ${user_token}" \
        --request POST -H "Content-Type: application/json" \
        --data "{\"name\":\"$name\",\"path\":\"$path\",\"visibility\":\"private\"}" \
        "http://${gitlab_host}/api/v4/groups"
}

createGitlabProject(){
    user_token=$1
    name=$2
    curl --silent --show-error --header "Authorization: Bearer ${user_token}" \
        --request POST -H "Content-Type: application/json" \
        --data "{\"name\":\"$name\",\"visibility\":\"private\"}" \
        "http://${gitlab_host}/api/v4/projects"
}

instantiateVocConfigurationRepo(){
    username=$1
    password=$2
    email=$3
    repo=$4

    git config --global user.name "$username"
    git config --global user.email "$email"

    git clone "http://${username}:${password}@gitlab.remip.eu/${username}/${repo}.git"
    cd ${repo}
    rm -rf *
    cp -af ${datadir}/wae_voc/. .
    git add --all
    git commit -m "add VOC configuration to build and deploy WAE"
    git push -u origin master
}

function dialog_message(){
    dialog --title "Showtime ! Build and deploy WhatStat on VOC" \
    --timeout 2 --msgbox "$1" 40 80 2> /dev/null
}
############################
##  Main
############################
clear

wordir=~/sandbox/
datadir=$(pwd)

cd $wordir
pwd

gitlab_root_password=rootroot
gitlab_host=gitlab.remip.eu

new_user_name=newuser
new_user_password=newuserpwd
new_user_email=newuser@mail.com
new_user_repo=awesome_group

root_api_token=$(readApiToken 'root' 'rootroot')


echo "Root api token is $root_api_token"

#echo "          $(createGitlabUser $new_user_name $new_user_password $new_user_email)"
echo "existing users: "
#echo $(queryGitlabApi $root_api_token '/users') | jq '.[].name'

#user_api_token=$(readApiToken $new_user_name $new_user_password)
echo "New user api token is $user_api_token"

#echo "          $(createGitlabProject $user_api_token $new_user_repo)"

#instantiateVocConfigurationRepo $new_user_name $new_user_password $new_user_email $new_user_repo

echo "TODO using dialog explain what does it do with WAE"

echo "TODO using dialog explain how to use WAE"