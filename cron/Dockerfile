FROM ubuntu:20.04

# Updating packages and installing cron
RUN apt-get update && apt-get install cron -y

# Copying script file into the container
COPY clean-script.sh .

# Giving executable permission to the script file
RUN chmod +x clean-script.sh

# Adding crontab to the appropriate location
ADD crontab /etc/cron.d/my-cron-file

# Giving permission to crontab file
RUN chmod 0644 /etc/cron.d/my-cron-file

# Running crontab
RUN crontab /etc/cron.d/my-cron-file

# Creating entry point for cron
ENTRYPOINT ["cron", "-f"]