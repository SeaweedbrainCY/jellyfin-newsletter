#! /bin/bash 

####
# Pre-execution set up 
# User: root
####

set -e 


if [ "$USER_UID" = "" ]; then
    echo "USER_UID environment variable is not set. Defaulting to 1001"
    USER_UID="1001"
fi

if [ "$USER_GID" = "" ]; then
    echo "USER_GID environment variable is not set. Defaulting to 1001"
    USER_GID="1001"
fi

# Check if config file exists

if [ -f /app/config/config.yaml ] && [ ! -f /app/config/config.yml ]; then
    echo "config.yaml found, renaming to config.yml"
    mv /app/config/config.yaml /app/config/config.yml
fi


if  [ ! -f /app/config/config.yml ]; then
    echo "WARNING. Config file not found. Generating default config."
    cp /app/default/config-example.yml /app/config/config.yml 
    cp /app/default/config-example.yml /app/config/config-example.yml
    echo "Config file generated at /app/config/config.yml"
    echo "Please edit the config file before running the application."
    echo "Refer to the documentation for more information."
    exit 1
fi

chown -R $USER_UID:$USER_GID /app/config/

mkdir -p /app/previews
chown -R $USER_UID:$USER_GID /app/previews/

###
# Switch to user 1001 and execute the main script
###

exec gosu $USER_UID:$USER_GID python main.py