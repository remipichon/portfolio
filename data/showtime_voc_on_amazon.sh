#!/bin/bash

# Amazon configuration page
dialog --title "Deploy on Amazon !" \
--ok-label "deploy" --form \
"I need a amazon access key and access token (careful to not copy paste a new line as it would press enter)" 25 60 16 \
"Amazon Access Key:" 1 1 "" 1 25 25 30 \
"Amazon Access Token" 2 1 "" 2 25 25 30


dialog --title "Terraform will create the AWS ressources" \
--timeout 2 --msgbox "Terraform will create the AWS ressources
It will do this
Blablabla" 40 80 2> /dev/null

clear
# Fall back to regular script to get Terraform then Ansible
echo "There
will
be
terraform
logs
followed
by
Ansible
logs
"

echo "******************"
echo "done, press any key continue ?"
read

dialog --title "Ansible will provision" \
--timeout 2 --msgbox "Ansible will provision
It will do this
Blablabla" 40 80 2> /dev/null


clear
# Fall back to regular script to get Terraform then Ansible
cat /var/loogs

echo "******************"
echo "done, press any key continue ?"
read

dialog --title "All good" \
--msgbox "Go to there to enjoy
Blablabla" 40 80


clear

