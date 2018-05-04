#!/bin/bash
DIALOG=${DIALOG=dialog}


#   welcome page
$DIALOG --clear --title " Hello you" \
--ok-label "Next" --msgbox " \
            YOU !
You are not the world, but let's pretend you are

And me ?

" 40 80

#   technologies list
dialog --clear --title " Bunch of technologies" \
--ok-label "see project" --checklist \
"Following is a list of the most needed technologies in modern software companies
Please select at least one to see the related projects I did
" 40 80 10 \
    docker "Docker" off \
    gitlab "Gitlab 8 with CI" off \
    swarm "Docker Swarm" off \
    k8s "Kubernetes" off \
    all "Or go to the projects list" off

#   projects page
dialog --clear --title "All the cool projects I did" \
--ok-label "Project detail" --radiolist \
"Please select one to see the details
" 40 80 10 \
  voc "VOC: automatic deployment for Swarm based on Gitlab CI (ready to deploy on Amazon)" off \
  mm "ManifMaker: stuff to assign" off \
  android "24Heures Android app" off

#   project page
dialog --clear --title "More about VOC" \
--no-label "back" --yes-label "deploy on Amazon ?" \
--yesno "Github or TOP_SECRET
    one sentence summary
    description:
        what ?
        user ?
        how ?
        with who ?" 40 80

# Amazon configuration page
dialog --title "Deploy on Amazon !" \
--ok-label "deploy" --form \
"I need a amazon access key and access token (careful to not copy paste a new line as it would press enter)" 25 60 16 \
"Amazon Access Key:" 1 1 "" 1 25 25 30 \
"Amazon Access Token" 2 1 "" 2 25 25 30


dialog --title "Terraform will create the AWS ressources" \
--timeout 2 --msgbox "Terraform will create the AWS ressources" 40 80 2> /dev/null

clear
# Fall back to regular script to get Terraform then Ansible
cat /var/loogs

echo "******************"
echo "done, press any key continue ?"
read

dialog --title "Ansible will provision" \
--timeout 2 --msgbox "Ansible will provision" 40 80 2> /dev/null


clear
# Fall back to regular script to get Terraform then Ansible
cat /var/loogs

echo "******************"
echo "done, press any key continue ?"
read

dialog --title "All good" \
--msgbox "Go to there to enjoy" 40 80


clear

#
#$DIALOG --title " Mon premier dialog" --clear \
#	--yesno "Bonjour, ceci est mon premier programme dialog" 10 30
#
#case $? in
#	0)	echo "Oui choisi. ";;
#	1)	echo "Non choisi. ";;
#	255)	echo "Appuyé sur Echap. ";;
#esac
#
#
# fichtemp=`tempfile 2>/dev/null` || fichtemp=/tmp/test$$
#trap "rm -f $fichtemp" 0 1 2 5 15
#$DIALOG --clear --title "Mon chanteur français favori" \
#	--menu "Bonjour, choisissez votre chanteur français favori :" 20 51 4 \
#	 "Brel" "Jacques Brel" \
#	 "Aznavour" "Charles Aznavour" \
# 	 "Brassens" "Georges Brassens" \
#	 "Nougaro" "Claude Nougaro" \
#	 "Souchon" "Alain Souchon" \
#	 "Balavoine" "Daniel Balavoine" 2> $fichtemp
#valret=$?
#choix=`cat $fichtemp`
#case $valret in
# 0)	echo "'$choix' est votre chanteur français préféré";;
# 1) 	echo "Appuyé sur Annuler.";;
#255) 	echo "Appuyé sur Echap.";;
#esac
#clear

