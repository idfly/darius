FROM sickp/centos-sshd

COPY . /root/.ssh
RUN chmod -R og-wrx /root/.ssh
RUN chown -R root: /root/.ssh

RUN useradd user
COPY ./user.pub /home/user/.ssh/authorized_keys
RUN chmod -R og-wrx /home/user/.ssh
RUN chown -R user: /home/user/.ssh
