#!/bin/bash
# need bash 4


# snif, I had nice Regex but bash =~ doesn't support 'g' flag, grep doesn't support look ahead '?=' and sed was a mess with all my parenthesis
# \[([a-z]*)(?=[a-z ]*(?:docker|gitlab)+[a-z ]*)
# [voc docker gitlab swarm node][mm node front][cloudbase docker swarm groovy][toast android python][wae front elastic spring][24happ android]

declare -A techno_projects
declare -A desc_projects

techno_projects[voc]="docker gitlab swarm node"
desc_projects[voc]="VOC: automatic build and deployment for Swarm based on Gitlab CI"

techno_projects[mm]="node front"
desc_projects[mm]="ManifMaker: stuff to assign"



selected_techos="node gitlab"

IFS=' ' read -r -a selected_technos_array <<< "$selected_techos"
#TODO support 'all', if 'all' copy techno_projects keys into selected_projects
declare -A selected_projects
for i in "${!techno_projects[@]}"
do
  for techno in "${selected_technos_array[@]}"
  do
    if [[ "${techno_projects[$i]}" == *"$techno"* ]];then
        selected_projects["$i"]="${desc_projects[$i]}"
        break
    fi
  done
done






### debug
for i in "${!techno_projects[@]}"
do
  echo "project  : $i"
  echo "technos: ${techno_projects[$i]}"
done

echo "selected technos"
for techno in "${selected_technos_array[@]}"
do
    echo "$techno"
done

echo "selected project"
for i in "${!selected_projects[@]}"
do
  echo "project  : $i"
  echo "desc: ${selected_projects[$i]}"
done




#options=(${selected_projects[@]})
#
#cmd=(dialog --keep-tite --radiolist "Select options:" 22 76 16)
#cmd+=( mm "ManifMaker: stuff to assign" off)
#cmd+=( mm "ManifMaker: stuff to assign" off)
##
#choices=$("${cmd[@]}" 2>&1 >/dev/tty)

#debug
title="Please select one to see the details"


#
command="dialog --clear --title \"I made this !\" --ok-label \"Project detail\" --cancel-label \"Go back to technos\" --radiolist "
command="$command \"$title\""
command="$command 30 100 40"
for i in "${!selected_projects[@]}"
do
  echo "project  : $i"
  echo "desc: ${selected_projects[$i]}"
  command="$command $i \"${selected_projects[$i]}\" off"
done

echo "==="
echo "${command}"
eval "${command}"






#        result=$(dialog --clear --title "I made this !" \
#        --ok-label "Project detail" --cancel-label "Go back to technos" --radiolist \
#        "$title" 30 100 40 \
#            voc "VOC: automatic build and deployment for Swarm based on Gitlab CI" off \
#            mm "ManifMaker: stuff to assign" off \
#            cloudbase "CloudbaseV2: custom orchestration tool for Swarm" off \
#            wae "WhatStat: web app to display statistics for Whatsapp chats" off \
#            24happ "24Heure Android: application to provide live information to festival attendees" off \
#            toast "ToastShooter: my first software, an Android replica of an iOs exclusive game" off 2>&1 1>&3)

















