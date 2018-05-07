# pages


doc
ftp://ftp.traduc.org/pub/lgazette/html/2004/101/lg101-P.html
https://linux.die.net/man/1/dialog


welcome page
    txt: Here is my portfolio using state of the art UI technologies
    


technologies page
    txt: Following is a list of the most needed technologies in modern software companies
         Please select at least one to see the related projects I did   
    


projects page
    txt: All the project I did OR All the project I did using <given technologies>
        Please select one to see the details        


project page
    txt: name
        Github or TOP_SECRET
        one sentence summary
        description:
            what ?
            user ?
            how ?
            with who ?
    txt: run terraform/ansible to provision it and see it live         
    
    
showtime page
    txt: let's go crazy and deploy it on Amazon ! 
    input: amazon access key and access token for terraform
    
    
# workflows  

* select some technologies to see related project
* select project to display details
* showtime: WAE on VOC
* showtime: VOC on Amazon   

## info only

## show time VOC on Amazon
* using terraform, provision on Amazon an EC2 and the correct access group
   * create terraform config 
* using Ansible, provision VOC (with one runner)
   * get runner token from user/password
* seed a test project (like a simple web server with a nice html) to deploy on local mode
   * see what's done for WAE

## show time WAE on VOC

| VOC is already fully setup with DNS records (needed for the mail)

* ~~create normal gitlab user~~
* ~~create empty project~~
* push repo with
  * simple-stack-instance on remote repo mode
  * .gitlab -ci.yml
* explain what does it do with WAE 
* explain how to use WAE  







