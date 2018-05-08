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
    cp -af ${datadir}/wae_voc/. .
    git add --all
    git commit -m "add VOC configuration to build and deploy WAE"
    git push -u origin master
}

function dialog_message(){
    dialog --title "Showtime ! Build and deploy WhatStat on VOC" \
    --timeout 200 --msgbox "$1" 40 80 2> /dev/null
}
############################
##  Main
############################
clear

wordir=/sandbox/
datadir=$(pwd)

mkdir -p $wordir
cd $wordir
pwd

gitlab_root_password=rootroot
gitlab_host=gitlab.remip.eu

new_user_name=newuser
new_user_password=newuserpwd
new_user_email=newuser@mail.com
new_user_repo=awesome_group


echo "Using the root password to retrieve root token"
root_api_token=$(readApiToken 'root' 'rootroot')
echo "      root api token is $root_api_token"
echo ""
echo "Using root token to create demo user ${new_user_name}"
echo "      response from Gitlab: $(createGitlabUser $new_user_name $new_user_password $new_user_email)"
echo "      existing users: "
echo $(queryGitlabApi $root_api_token '/users') | jq '.[].name'
echo ""
echo "Using demo user password to retrieve demo user token"
user_api_token=$(readApiToken $new_user_name $new_user_password)
echo "      new user api token is $user_api_token"
echo ""
echo "Using demo user token to create empty repository ${new_user_repo}"
echo "      response from Gitlab: $(createGitlabProject $user_api_token $new_user_repo) | jq '.'"
echo ""
echo "Git commit to ${new_user_repo} the VOC configuration files to build and deploy Whatstat"
instantiateVocConfigurationRepo $new_user_name $new_user_password $new_user_email $new_user_repo
echo ""
echo "*******************************************************************************"
echo "Everything is configured, time to chill out and appreciate what have been done.

    A repo have been created with a few files: http://${gitlab_host}/${new_user_name}/${new_user_repo}/tree/master
    Only 3 files are needed to trigger the build and deploy of WhatStat which is hosted on Github
      * .gitlab-ci.yml: standard Gitlab8 configuration file
         in 'script', 'node /root/app/app.js' which start the VOC Runner App
         in 'tags', the Voc Runner defined in Gitlab has the same tag, go check http://gitlab.remip.eu/admin/runners login as the root user
         finally the artifact to retrieve the result in json format

      * repo.whatsappelastic.json: defines a VOC resource to explain where to find the code to build. It supports SSH credentials as well
      * simple-stack-instance.remote-repo.wae.wae.json: defines a VOC resource to build and instantiate a stack. Let's break it down
         simple-stack-instance      this instance directly refers a docker compose, VOC support stack definitions to assemble several docker compose in one stack
         .remote-repo               the Dockerfile, Compose files and potentially context for the build are to be found on a separate repo, not the VOC one. VOC supports code coming from itself or from an external repo.
         .wae                       the Docker Compose to look for is named docker-compose.wae.yml and is in the remote repo. Which repo is defined in the file
         .wae                       the stack deployed from the compose file with 'docker stack deploy' will be named 'wae', Docker services will get names like 'wae_elasticsearch'
         .json                      yeah, Json !

    Now, let's see it in actions (well, it might have already finished by the time you are done reading, but still)
        http://${gitlab_host}/${new_user_name}/${new_user_repo}/-/jobs
        click on the most recent one, which should be 'running' or 'passed' and read the logs to see what VOC did for you. You can as well download the artifact to act on it with an external tool.

    And it's not done, VOC doesn't only build and deploy, it comes with some other features critical for hosting web app.

    It provides an Nginx with service discovery to proxy requests to your service.
    \"Provided your DNS is setup to forward foo.bar.com to the host running nginx-proxy, the request will be routed to a container with the VIRTUAL_HOST env var set.\"
    from https://github.com/jwilder/nginx-proxy

    ??? TODO how to access the non DNS deployed app ??? open a port for demo purposes ????

    It also provides forwarding incoming emails to the HTTP endpoint of your choice.
    It used the npm package mailin: https://www.npmjs.com/package/mailin

    If you create the correct DNS records, you can start sending emails to WhatStat !
        * DNS A to gitlab.remip.eu
        * DNS MX record see DNS configuration https://www.npmjs.com/package/mailin#the-crux-setting-up-your-dns-correctly

    If not, WhatStat is running configured at http://whatstat.remip.eu/, go take a look an try out sending a chat.

    Hope it worked for you, hope you liked it and are now interested in either VOC or WhatStat.
"
