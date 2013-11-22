# Sheer Industries Cloud System

This repository contains a server and a Qt/C++ client library which can be used with it.
Server can accept files, allow user to delete/download/render them afterwards.
Communication happens over HTTP/HTTPS protocol.

## Commands

The server understands usual http verbs formed as:

   *server address*/*verb*?param1=value1&param2=value2...

Each connection contains user/password, so no special login is required beforehand at the moment.

Existing verbs include:

- [x] "/authorize" : Verify that user/pass are ok, not necessary for other tasks.
- [x] "/upload"    : Post contents of a file to server.
- [x] "/list"      : Retrieve list of files starting with the provided prefix with their checksums.
- [x] "/download"  : Retrieve contents of a file from server.
- [x] "/delete"    : Remove file from server.
- [x] "/job"       : Starts rendering on a file.

## File locations
Each user has its own folder for his projects. Same files, for example models, are done using hardlinks. The structure is the same as on the user's local machine.

## Jobs 
To start a rendering job, user uploads the .xml file with meta-info about the job, and calls /job with the xml file.
Rendering result is written as follws: example.xml -> example.xml.png 
