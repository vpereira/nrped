
GO NRPE - Nagios Remote Plugin Executor 

[![Build Status](https://travis-ci.org/vpereira/nrped.svg?branch=master)](https://travis-ci.org/vpereira/nrped)

Contents
--------

There are two pieces to this addon:

  1) NRPE       - This program runs as a background process on the 
                  remote host and processes command execution requests
	              from the check_nrpe plugin on the Nagios host.
		          Upon receiving a plugin request from an authorized
                  host, it will execute the command line associated
                  with the command name it received and send the
                  program output and return code back to the 
                  check_nrpe plugin

  2) check_nrpe - This is a plugin that is run on the Nagios host
                  and is used to contact the NRPE process on remote
	              hosts.  The plugin requests that a plugin be
                  executed on the remote host and wait for the NRPE
                  process to execute the plugin and return the result.
                  The plugin then uses the output and return code
                  from the plugin execution on the remote host for
                  its own output and return code.


NOTE:

It's a ongoing project, however it works. There are some features missing, like SSL support, systemd integration, that will be done with the time. git pull are, as usual, welcome :o)


Compiling it:

you will need go and make

    make
