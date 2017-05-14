FROM golang:1.8.1

RUN apt-get update
RUN apt-get install -y sudo git markdown

ADD *.sh /
RUN chmod +x /*.sh


EXPOSE 8080
EXPOSE 8081


RUN git clone https://github.com/kr-g/go-sample.git /repo/go-sample


ENTRYPOINT /startup.sh /



