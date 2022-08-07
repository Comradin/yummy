[![Go Report Card](https://goreportcard.com/badge/github.com/Comradin/yummy)](https://goreportcard.com/report/github.com/Comradin/yummy) [![Build Status](https://travis-ci.org/Comradin/yummy.svg?branch=master)](https://travis-ci.org/Comradin/yummy)

# yummy
A simple YUM repository server to upload to and serve rpm packages from.

It provides a simple API (/api/upload) to use tools like curl to upload rpm
packages and automatically rebuild the rpm repository metadata.

# API endpoints
## /api/upload
The upload endpoint emulates a filled-in form in which a user has pressed
the submit button. The formular field that contains the file is **fileupload**.

To use this with curl use the -F switch with the following format:
_-F fileupload=@file.rpm_

```example
    curl -F fileupload=@python-elasticsearch-curator-5.1.1-1.noarch.rpm \
         http://yum.example.com/api/upload
```

**Caution**: This is the **only** supported way of uploading files to the
service at the moment. Using a different form value or an http PUT will
result in an error.

## /api/delete
The delete endpoint deletes a given rpm.
The file is declared by adding it to the url:

```example
    curl -X DELETE http://localhost:8080/api/delete/python-elasticsearch-curator-5.1.1-1.noarch.rpm
```

## / BaseURL
Will serve the repository content as a Filebrowser, so it can be used as
the baseurl in a yum repository config file.

```
[yum-repository]
baseurl = https://yum.example.com/
enabled = 1
name = YUM Repository Example
```
**Caution** Will serve a blank page in case no files were uploaded yet.

## /help
Will render the README.md file as HTML.

# configuration
yummy can be configured by editing the .yummy.yml file.

```
yum:
  repopath: /opt/
  workers: 2
  createrepoBinary: /bin/createrepo
  rpmBinary: /bin/rpm
  helpFile: /usr/share/doc/yummy/README.md
  auth:
    enabled: false
    user: yummy
    password: yummy
```

The options should be mostly self explanatory. Here a brief overview:

1. repopath, where to store the rpms and create the repository
2. workers, the amount of workers for the createrepo call to index the
   packages. If the update process takes too long, more worker processes
   might work
3. createrepoBinary, path to the executable to create the repository metadata
4. rpmBinary, not in use yet, for further verification of uploaded rpm files
5. helpFile, path to load the help file from
6. auth enabled, enable/disable authentication
7. auth user, the user for authentificate for the delete endpoint
8. auth password, the password for authentificate for the delete endpoint

# hints for contributors
This project uses the Go dependency management tool `dep`

If you added a new dependency, please add it with `dep ensure -add`
