#!/bin/bash

# define data
declare -A techno_projects
declare -A desc_projects

techno_projects[voc]="docker gitlab swarm node"
desc_projects[voc]="VOC: automatic build and deployment for Swarm based on Gitlab CI"

techno_projects[manifmaker]="node front"
desc_projects[manifmaker]="ManifMaker: stuff to assign"

techno_projects[cloudbase]="docker swarm groovy"
desc_projects[cloudbase]="CloudbaseV2: custom orchestration tool for Swarm"

techno_projects[wae]="front elastic spring"
desc_projects[wae]="WhatStat: web app to display statistics for Whatsapp chats"

techno_projects[app]="android"
desc_projects[app]="24Heure Android: application to provide live information to festival attendees"

techno_projects[toast]="android python"
desc_projects[toast]="ToastShooter: my first software, an Android replica of an iOs exclusive game"

# Define the dialog exit status codes
: ${DIALOG_OK=0}
: ${DIALOG_CANCEL=1}
: ${DIALOG_ESC=255}

function technoView(){
    # Duplicate file descriptor 1 on descriptor 3
    exec 3>&1

    result=$(dialog --clear --title " Bunch of technologies" \
    --ok-label "See related projects" --cancel-label "Go back home" --checklist \
    "Following is a list of some of the most vibrant technologies used amongst successful companies making earning unicorns per quarters.
    Please select at least one to see details about some projects I work on using those mind blowing state of art nicely packaged computer binaries.
    " 30 100 40 \
        docker "Docker" off \
        gitlab "Gitlab 8 with CI" off \
        swarm "Docker Swarm" off \
        node "NodeJs" off \
        groovy "Groovy" off \
        spring "Java with Spring" off \
        android "Android" off \
        elastic "ElasticSearch" off \
        python "Python" off \
        front "Frontend JS" off \
        all "Or go to the projects list" off 2>&1 1>&3)

    # Get dialog's exit status
    return_value=$?

    # Close file descriptor 3
    exec 3>&-
}

function projectView(){

    # if none or at least 'all' technos selected, same as 'all' selected alone
    if [[ "${selected_techos}" == "" ]] || [[ "${selected_techos}" == *"all"* ]] ; then
        echo "selected_techos either empty or contains 'all'"
        selected_technos_array=("all")
    else
       IFS=' ' read -r -a selected_technos_array <<< "$selected_techos"
    fi

    # compute projects to display according to selected technos
    declare -A selected_projects
    for i in "${!techno_projects[@]}"
    do
      for techno in "${selected_technos_array[@]}"
      do
        if [[ "${techno_projects[$i]}" == *"$techno"* ]] || [[ "$techno" == "all" ]];then
            selected_projects["$i"]="${desc_projects[$i]}"
            break
        fi
      done
    done

    result=""
    first_time=true
    while [ -z "$result" ]
    do
        exec 3>&1

        # compute title
        if [[ "$selected_technos_array" != "" ]] && [[ "$selected_technos_array" != "all" ]]; then
            title="Displayed are all the projects that used one of the selected technos:
      $selected_techos" #this is text formatting
        else
           selected_technos_array=("all")
           title="Displayed are all the projects"
        fi

        if $first_time; then
            title="$title
Please select one to see the details"
        else
            title="$title
!! Hunhun, you need to select at least one project !!"
        fi
        first_time=false

        # generate dialog command
        command="dialog --clear --title \"I made this !\" --ok-label \"Project detail\" --cancel-label \"Go back to technos\" --radiolist "
        command="$command \"$title\""
        command="$command 30 100 40"
        for i in "${!selected_projects[@]}"
        do
          echo "project  : $i"
          echo "desc: ${selected_projects[$i]}"
          command="$command $i \"${selected_projects[$i]}\" off"
        done
        result=$(eval "${command}" 2>&1 1>&3)

        # partly handle result
        return_value=$?
        exec 3>&-
        case $return_value in
        $DIALOG_CANCEL)
            echo "Cancel pressed (displayProjects)"
            return;;
        $DIALOG_ESC)
            echo "ESC pressed (displayProjects)"
            return;;
        esac
    done
}


function technoController(){
    technoView

    case $return_value in
      $DIALOG_OK)
        echo "selected technos $result"
        selected_techos=$result
        projectController
        ;;
      $DIALOG_CANCEL)
        echo "Cancel pressed (displayTechos)"
        exit 0;;
      $DIALOG_ESC)
          echo "ESC pressed (displayTechos)"
          exit 0;;
    esac
}

projectController(){
    projectView

    case $return_value in
      $DIALOG_OK)
        echo "selected project $result"
        selected_project=$result
        #TODO
        echo "TODO go to project detail"
        ;;
      $DIALOG_CANCEL)
        technoController
        ;;
      $DIALOG_ESC)
          echo "ESC pressed (displayProjects)"
          exit 0;;
    esac
}



technoController

#TODO projectDetailView and projectDetailController


