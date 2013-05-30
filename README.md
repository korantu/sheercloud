# Sheer Industries Cloud System

This repository contains a server and a Qt/C++ client library which can be used with it.
Server can accept files, allow user to delete/download/render them afterwards.
Communication happens over HTTP/HTTPS protocol.

## Commands

The server understands usual http verbs formed as:

   <server address>/<verb>?param1=value1&param2=value2...

Each connection contains user/password, so no special login is required beforehand at the moment.

Existing verbs include:

- "/authorize" : Verify that user/pass are ok, not necessary for other tasks.
- "/upload"    : Post contents of a file to server.
- "/download"  : Retreive contents of a file from server.
- "/delete"    : Remove file from server.
