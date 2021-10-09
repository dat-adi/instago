# Instago

A simple REST API server, connected with a MongoDB backend.\
This project has been built to submit in the Task Round for the internship opportunity at Appointy.

## Tasks
Here's a few of the requirements listed and the status of their completion.

- [X] `POST` : Create a User.
- [-] `GET`  : Get a User using the ObjectID.
- [X] `POST` : Create a Post.
- [-] `GET`  : Get a Post using the ObjectID.
- [-] `GET`  : List all the Posts using the User ObjectID.

Additional extension for the task.
- [X] Documented.
- [X] Pagination.

## Notes

While the creation of the code itself was relatively easy, the configuration setup for Go, was quite the hurdle.
So, there will be a few issues with the modules and the compatibility.

The routing for the GET requests continue to fail for some reason, unsure of what the issue is,
and while the general routing for POST and GET requests with parameters are feasible and easy to implement.
The parsing of the `:id` attribute proves to be a challenge with the native `http` library.

The passwords, before storage, are encrypted using the SHA512 hash algorithm.
The generated salt is not stored as of right now, considering that there is no login interface present, 
as a form of validation for the password field.

DB References are still unsupported in MongoDB for Golang, and as such, I've opted for a manual DB references for the final task.

However, none of the sections in the application go beyond the provided list of standard Golang libraries and MongoDB drivers.
