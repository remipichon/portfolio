#!/bin/bash

readApiToken(){
    gitlab_root_username=$1
    gitlab_root_password=$2
    token_request=$(curl --show-error --silent --request POST -H "Content-Type: application/json" \
    --data "{\"grant_type\":\"password\",\"username\":\"$gitlab_root_username\",\"password\": \"$gitlab_root_password\"}" \
    $gitlab_host/oauth/token || exit 1)

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
    curl --silent --show-error --header "Authorization: Bearer ${api_token}" "http://${gitlab_host}/api/v4/${endpoint}" || exit 1
}

createGitlabUser(){
    username=$1
    name=$1
    password=$2
    email=$3
    curl --silent --show-error --header "Authorization: Bearer ${root_api_token}" \
        --request POST -H "Content-Type: application/json" \
        --data "{\"email\":\"$email\",\"password\":\"$password\",\"username\":\"$username\",\"name\":\"$name\",\"skip_confirmation\":\"true\"}" \
        "http://${gitlab_host}/api/v4/users" || exit 1
}

createGitlabGroup(){
    user_token=$1
    name=$2
    path=$2
    curl --silent --show-error --header "Authorization: Bearer ${user_token}" \
        --request POST -H "Content-Type: application/json" \
        --data "{\"name\":\"$name\",\"path\":\"$path\",\"visibility\":\"public\"}" \
        "http://${gitlab_host}/api/v4/groups" || exit 1
}

createGitlabProject(){
    user_token=$1
    name=$2
    curl --silent --show-error --header "Authorization: Bearer ${user_token}" \
        --request POST -H "Content-Type: application/json" \
        --data "{\"name\":\"$name\",\"visibility\":\"private\"}" \
        "http://${gitlab_host}/api/v4/projects" || exit 1
}

instantiateVocConfigurationRepo(){
    username=$1
    password=$2
    email=$3
    repo=$4

    git config --global user.name "$username"
    git config --global user.email "$email"

    rm -rf $repo
    git clone "http://${username}:${password}@gitlab.remip.eu/${username}/${repo}.git" || exit 1
    cd ${repo}
    cp -af ${datadir}/wae_voc/. .
    sed -i "s/PUBLIC_ACCESS_VALUE/$virtual_host/g" simple-stack-instance.remote-repo.wae.wae.json
    sed -i "s/VIRTUAL_HOST_VALUE/http://$virtual_host/g" simple-stack-instance.remote-repo.wae.wae.json
    if [ -z "$public_port" ]; then
        sed -i "s/PUBLIC_PORT_VALUE//g" simple-stack-instance.remote-repo.wae.wae.json
    else
        sed -i "s/PUBLIC_PORT_VALUE/$public_port:80/g" simple-stack-instance.remote-repo.wae.wae.json
    fi
    for f in *STACKNAME*; do mv "$f" "`echo $f | sed s/STACKNAME/$stack_name/`"; done
    git add --all
    git commit -m "[do-all] add VOC configuration to build and deploy WAE"
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

default_gitlab_host=gitlab.remip.eu

## deploy the real WAE
#gitlab_root_password=rootroot
#new_user_name=wae_user
#new_user_password=wae_user
#new_user_email=wae_user@mail.com
#new_user_repo=wae_config
#gitlab_host=$default_gitlab_host
#virtual_host=whatstat.remip.eu   # for the proxy, to access wae
#public_port=    #if no DNS records, use public port
#stack_name=wae
#
## deploy the demo WAE
gitlab_root_password=rootroot
new_user_name=demo_user3
new_user_password=demo_user3
new_user_email=demo_user@mail.com
new_user_repo=demo_wae_config3
gitlab_host=$default_gitlab_host
virtual_host=
public_access=
public_port=30200
stack_name=demo_wae3


#echo "Welcome to the wizard that will guide you through the process on building and deploying WhatStat, a Java based web app.
#If I gave you the root password to my ${default_gitlab_host}, then you can proceed.
#Else, you might want to deploy VOC on Amazon Web Service and use it to deploy WhatStat.
#You can also use the build it and deploy it yourself using the given Dockerfiles and Docker Composes. "
#echo ""
#echo "To start, I need a few parameters: "
#read -p "   confirm the hostname ${default_gitlab_host} (press enter) or override it with your own VOC (type it!): " given
#gitlab_host=${given:-$default_gitlab_host}
#
#printf "   Gitlab root password to create a user: "
#read -s gitlab_root_password
#echo ""
#
#echo "A demo user will be created, you are free to name it the way you want"
#read -p "   Demo user name: " new_user_name
#printf "   Demo user password: "
#read -s new_user_password
#echo ""
#read -p "   Demo user email: " new_user_email
#read -p "   Demo repository: " new_user_repo
#echo ""
#
#echo "Now, I would like some details about the WhatStat you are going to deploy:"
#read -p "   Deployed WhatStat stack name (given to 'docker stack deploy) ? " stack_name
#read -p "   Do you have a A DNS record pointing to '164.132.42.48' or to 'gitlab.remip.eu' ? Y\N " yesno
#if [ "$yesno" == "Y" ] || [ "$yesno" == "y" ]; then
#    read -p "       DNS record that will be redirected to the WhatStat you are deploying: " virtual_host
#else
#    public_port=$(shuf -i 30000-34000 -n 1)
#    echo "        Without DNS records, the only way to access WhatStat is via a port mapping, a random one has been choosen for you"
#fi
#
#if [[ "$stack_name" == "wae" ]]; then
#    echo "Sorry, 'wae' is not accepted as it would remove the live demo of WhatStat"
#    exit 1
#fi

echo ""
echo "Configuration is done, here is a summary"
echo "      Using VOC running on '${gitlab_host}', WhatStat will be deployed via a repository '${new_user_repo}' linked to the user '${new_user_name}' ('${new_user_email}')."
echo ""
echo "Now, lets appreciate..."
echo ""

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

    A repo have been created with a few files: http://${gitlab_host}/${new_user_name}/${new_user_repo}/tree/master (login with ${new_user_name} and the password you gave earlier)
    Only 3 files are needed to trigger the build and deploy of WhatStat which is hosted on Github
      * .gitlab-ci.yml: standard Gitlab8 configuration file
         in 'script', 'node /root/app/app.js' which start the VOC Runner App
         in 'tags', the Voc Runner defined in Gitlab has the same tag, go check http://gitlab.remip.eu/admin/runners (logout then login as the Gitlab root user)
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
    from https://github.com/jwilder/nginx-proxy"
if [ -z "$virtual_host" ]; then
    echo "    You can access your WhatStat on the DNS record you provided ${virtual_host}/${stack_name}/legacy."
else
    echo "    You can access your WhatStat on the public port generated for you ${gitlab_host}:${public_port}"
fi
echo "
    It also provides forwarding incoming emails to the HTTP endpoint of your choice.
    It used the npm package mailin: https://www.npmjs.com/package/mailin

    If you create the correct DNS records, you can start sending emails to WhatStat !
        * DNS A to gitlab.remip.eu
        * DNS MX record see DNS configuration https://www.npmjs.com/package/mailin#the-crux-setting-up-your-dns-correctly

    If not, WhatStat is running configured at http://whatstat.remip.eu/, go take a look an try out sending a chat.

    Hope it worked for you, hope you liked it and are now interested in either VOC or WhatStat.

    If you want to have fun, please edit the VOC configuration to 'enable: false' in order to release my really small server resources.
    Gitlab might crash if you deploy to much because my server is a patato really small VPS.
"
