# gittagger
A rest API interface for tagging or changing files in git repositories. 

Useful for triggering CI that uses a git tag as input or performs a special build based on an input file contents.

An example usage is when you want to use CI for triggering a container run that performs some long running business job. In this case you can parametrize the build with tag name or an input file contents and leverage all CI tooling for controlling the execution of this task (even thou it is not a "regular" development task).

## Example

* Create file 'docker-compose.yml'
```yml
version: '3.7'

services:

  gittagger:
    image: flaviostutz/gittagger
    ports:
      - 50000:50000
    environment:
      - GIT_REPO_URL=https://myuser:mypass@github.com/flaviostutz/gittagger-test.git
      - GIT_USERNAME=myuser
      - GIT_EMAIL=myuser@mail.com
```
* Run 'docker-compose up'

* Push a new tag to repository
```shell
curl -X POST http://localhost:50000/tag/1.0.2
```
```json
{"message":"Tag 1.0.3 pushed successfully to git repository"}
```

* Push new file contents to repository
```shell
curl -X POST http://localhost:50000/files/test2 -d test2
```
```json
{"message":"File 'test2' updated and pushed to git repo successfully"}
```

