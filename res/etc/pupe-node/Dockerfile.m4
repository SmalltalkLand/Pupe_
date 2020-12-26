changeword(`@\([_a-zA-Z0-9]*\)')
FROM ubuntu:20.04
RUN apt-get update
RUN apt-get install -y pcmanfm xterm xserver-xephyr openbox

RUN apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg-agent \
    software-properties-common && curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add - && add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable" && apt-get update && apt-get install -y docker-ce docker-ce-cli containerd.io && curl -fsSL https://raw.githubusercontent.com/mviereck/x11docker/master/x11docker | bash -s -- --update
   
RUN snap install anbox && apt-get install -y adb

RUN snap install chromium && apt-get install -y unzip xvfb libxi6 libgconf-2-4 default-jdk python-pip && wget https://chromedriver.storage.googleapis.com/2.41/chromedriver_linux64.zip && unzip chromedriver_linux64.zip && mv chromedriver /usr/bin/chromedriver && chmod +x /usr/bin/chromedriver && pip install selenium
