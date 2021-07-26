FROM  ubuntu

ADD . /goworkspace/src/gossh-honey
WORKDIR /goworkspace/src/gossh-honey
EXPOSE  2222
CMD ["/bin/bash"]