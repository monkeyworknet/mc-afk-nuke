MC AFK NUKER
------------

Minecraft AFK detection isn't that great, if a player is at an AFK fisher or any other such farms they can sit there forever.   This can be a problem on servers where you do not want to enable single player sleep command blocks.   I wrote this script in GOLANG (partly to help myself learn the langauge) and it checks each players stats file and looks for movement .. if a player doesn't move they are marked as AFK no matter what they are swinging with the arms / etc.

The code isn't the cleanest but it is functional.   Feel free to fork and make improvements or make suggestions on how I can improve.



Prerequisites
-------------

1.  Minecraft running on a linux server
2.  An init script for starting/stoping minecraft (example init included or https://minecraft.gamepedia.com/Tutorials/Server_startup_script)
          -  Note if you are using a different start script the program is written to expect to give the following command to kick:
			/usr/sbin/service/minecraft command kick PLAYER REASON
	  -  If you need to change this edit the main.go file and modify "kickcmd" and "kickargs" variables.
3.  An ability to use cron or some other scheduler to run the script on a regular basis.
4.  Your minecraft server needs to be able to be quieried:
          -  In your server.properties file for minecraft make sure "enable-query=true" and "query.port=25565" is set (port can be whatever you wish)


Installing Binary
-----------------

1.  Place mc-afk-nuke in any directoy you like (recommend /home/minecraft)
2.  Move default-config.json to config.json and edit it as required
3.  Run  ./mc-afk-nuke  with no command options and insure it works (make sure someone is on the server when you execute this) 
4.  If it executes properly then setup a crontab to run it on a regular schedule (see schedule example below)
5.  A log file is created and appended to everytime the script runs.   You may also want to setup a job to rotate this log as well.


Schedule Example
----------------

Let's say you want to boot anyone who is AFK for more than 25 mins,  but to keep your players from knowing the exact moment they will be kicked you want it be set to a variable between 24 and 30 mins.

Set afkkickvaluemin to 8 and akfkickvaluemax to 10 and then set the cron job to run every 3 mins (3x8 = 24 mins,   3x10 = 30 mins).

If you don't want it to be variable set the min and max values to the exact same (lets say 10) and run the cron job every 3 mins and it will kick anyone afk for more than 30 mins.


Building from Source
-------------------

Built with go version go1.10 linux/amd64
Uses all standard libs + "github.com/nanobox-io/golang-scribble"

Entire program is held in main.go with exception of the configuration section - should be split out more in future.


Future To-Do
------------

Refactor and simplify code
Write go tests 
Expand the configuration options  
Allow the ability to mark a person AFK if they don't move a certian distance (right now any movement resets you)



