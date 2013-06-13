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

Verbs under development are:
- [ ] Clever uploading to avoid sending already known files
- [ ] "/render" : Starts rendering on a file
- [ ] "/status" : Provides a status for a rendering job
- [ ] "/cancel" : Stops a rendering job

## File locations
Each user has its own folder for his projects. Same files, for example models, are done using hardlinks. Most likely the structure is the same as on the user's local machine, but it is TBD for now.
