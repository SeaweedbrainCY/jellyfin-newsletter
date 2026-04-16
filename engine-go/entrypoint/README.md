# Jellyfin-Newsletter entrypoint

Jellyfin-Newsletter needs persistent data to 
- Read the configuration file written by the user (real-only)
- Save some data related to previous newsletter (e.g. the last newsletter data, read/write)

In another hand, Jellyfin-Newsletter should **never** run as root for countless security reasons. 

That being said, the mounted folders need to be configured by the script to have the correct permissions (the is the only use case for now). 

For user convenience, it has been decided to **not** rely on the user to configure correctly the permissions of mounted directories:
- The user is not necessarily aware of the exact UID/GUID of the user set up at build 
- The script should have to run several checks at startup to immediatly warn the user if anything is incorrect.
- In some setup, the user barely interact with the filesystem/docker compose file and therefore, has limited capacity/action on the directories.
- It requires handy, not portable, and manual configuration, which defeat one of the Jellyfin-Newsletter purpose: Be easy and plug-and-play. 


## The solution
As many docker containers are currently doing it, there are several ways to do this. The main idea is always the same: **start the container as root. Let the entrypoint configure permission/filesystem and then DROP the privileges to execute (no fork) the program with a non-root user.**

Since Jellyfin-Newsletter Go engine is using distroless as a base image, no shell is available to do this. 

This is the whole purpose of this entrypoint Go script. It configures, as Jellyfin-Newsletter script will require, the filesystem and then **drops** the privileges and **executes** Jellyfin-Newsletter script with low privileges.
