#!/bin/bash

# Define data
declare -A techno_projects
declare -A desc_projects
declare -A detail_projects

techno_projects[voc]="docker gitlab swarm node"
desc_projects[voc]="VOC: automatic build and deployment for Swarm based on Gitlab CI"
detail_projects[voc]="
https://github.com/remipichon/voc
Voc is Docker close friend, it can build and deploy on Docker Swarm

Voc is so easy to enable. If you already have a Gitlab 8 and a Gitlab Runner defined, you just need to create a new runner using the Voc Runner NodeJs App available on Docker Hub and voilà ! If you don't have Gitlab already, which is bad, you can easily have it with Voc all configured via the given Ansible playbook that does the heavy lifting.

Voc is stunning by its simplicity. With a few configuration in json format you can build and push images and deploy from docker compose file. Voc is awesome, yes, but it's partly because it uses another wonderful tool; Gitlab with CI enable. Each time a config or any Docker related files is committed, the Voc Runner App is triggered to apply just what it needs, smartly detecting with files are related to which Voc instance. So smart, it also monitor the context or your configured images via Dockerfile and can build and deploy your master branch each time a merge request is accepted.
Voc should start making the coffee, let's create a feature request...

But first, would you like to see it live? Go back to 'technos list' and hit 'Showtime !'

"

techno_projects[manifmaker]="node front"
desc_projects[manifmaker]="ManifMaker: stuff to assign"
detail_projects[manifmaker]="
"

techno_projects[cloudbase]="docker swarm groovy"
desc_projects[cloudbase]="CloudbaseV2: custom orchestration tool for Swarm"
detail_projects[cloudbase]="
"

techno_projects[wae]="front elastic spring"
desc_projects[wae]="WhatStat: web app to display statistics for Whatsapp chats"
detail_projects[wae]="
"

techno_projects[app]="android"
desc_projects[app]="24Heure Android: application to provide live information to festival attendees"
detail_projects[app]="
"

techno_projects[toast]="android python"
desc_projects[toast]="ToastShooter: my first software, an Android replica of an iOs exclusive game"
detail_projects[toast]="
I didn't know Git at the time. No visioning, it makes it even more impressive.
A funny catchy looking smartphone game which consist of shooting toasts jumping out of toasters.

This project have a nice story. Have a tea with some biscuits and appreciate.
It's my first project ever, my discovery of the magnifique world of computing. I used to play a game on iOs, ToastShooter, for hours. I was scoring around 2000 then 4000, sometimes 600 but rarely. It was a challenge de beat my record. One day, my thumbs slipped on my screen until I reached 1600, which I couldn't ever beat. It was done, I lost the envy to play that game with one impossible challenge. I also lost my iPhone and the game doesn't exist on Android.

Then started the new challenge; play it once more, on Android.
I used a Python library that could export to apk for Android; PyGame. I developed it, my specs being my memories of the game I knew by heart. I implemented all the rules, the timings, the points systems and even shamely took the graphics from Google. It was running at 240 fps on my patato laptop, 30 fps on my Android, joy ! It was really joy but it didn't play it, instead I conceived some new type of game. My favorite was the one where you are a splashed egg and you have to jump from toast to toast using the accelerometer to lead your jump. Careful to the burnt toasts !

I still find it impressive, since I moved to the web and then to Android in Java which was so much move convenient but no game anymore.

I can try to offer it to you, via a Docker image and XTerm forwarding, hope it works.
"

# Define the dialog exit status codes
: ${DIALOG_OK=0}
: ${DIALOG_CANCEL=1}
: ${DIALOG_ESC=255}

function technoView(){
    # Duplicate file descriptor 1 on descriptor 3
    exec 3>&1

    result=$(dialog --clear --title " Bunch of technologies" \
    --ok-label "See related projects" --cancel-label "Exit" --checklist \
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
        all "Or go to the projects list" off \
        showtime "What about some Showtime !" off 2>&1 1>&3)

    # Get dialog's exit status
    return_value=$?

    # Close file descriptor 3
    exec 3>&-
}

function projectDetailView(){
    exec 3>&1
    result=$(dialog --clear --title "More about ${selected_project}" \
     --yes-label "Go back to technos list" --no-label "Go back to project list"  --yesno \
     "${detail_projects["${selected_project}"]}" 30 100 2>&1 1>&3)

    return_value=$?
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

    # compute projects to be displayed according to selected technos
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
            echo "Cancel pressed (projectView in)"
            return;;
        $DIALOG_ESC)
            echo "ESC pressed (projectView in)"
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
        if [[ "${selected_techos}" == *"showtime"* ]] ; then
            ./showtime_wae_on_voc.sh
        else
            projectController
        fi
        ;;
      $DIALOG_CANCEL)
        clear
        echo "Bye, don't forget to drop me a message on https://www.linkedin.com/in/remipichon/"
        exit 0;;
      $DIALOG_ESC)
          echo "ESC pressed (technoView)"
          exit 0;;
    esac
}

function projectController(){
    projectView

    case $return_value in
      $DIALOG_OK)
        echo "selected project $result"
        selected_project=$result
        projectDetailController
        ;;
      $DIALOG_CANCEL)
        technoController
        ;;
      $DIALOG_ESC)
          echo "ESC pressed (projectView)"
          exit 0;;
    esac
}

function projectDetailController(){
    projectDetailView

    case $return_value in
      $DIALOG_OK)
        technoController
        ;;
      $DIALOG_CANCEL)
        echo "Cancel pressed (projectDetailView)"
        projectController
        ;;
      $DIALOG_ESC)
          echo "ESC pressed (projectDetailView)"
          exit 0;;
    esac
}

technoController