
# GO NRPE - Nagios Remote Plugin Executor 

![test.yml](https://github.com/vpereira/nrped/actions/workflows/test.yml/badge.svg)
![codeql-analysis.yml](https://github.com/vpereira/nrped/actions/workflows/codeql-analysis.yml/badge.svg)
[![Build Status](https://travis-ci.org/vpereira/nrped.svg?branch=master)](https://travis-ci.org/vpereira/nrped)

## Status: Maintained

Contents
--------

There are two pieces to this addon:

  * `nrped`: This program runs as a background process on the remote host and processes command execution requests
	     from the check_nrpe plugin on the Nagios host.  Upon receiving a plugin request from an authorized
             host, it will execute the command line associated with the command name it received and send the
             program output and return code back to the `check_nrpe` plugin

  * `check_nrpe`: This is a plugin that is run on the Nagios host and is used to contact the `nrped` process on remote
	          hosts.  The plugin requests that a plugin be executed on the remote host and wait for the `nrped`
                  process to execute the plugin and return the result.
                  The plugin then uses the output and return code from the plugin execution on the remote host for
                  its own output and return code.


NOTE:

It's a ongoing project, however it works. There are some features missing, like SSL support that will be done with the time. git pull are, as usual, welcome :o)


Compiling it:

you will need go and make

    make
