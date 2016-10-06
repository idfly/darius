Darius
======

Shell tasks runner: tangle webhook and shell task in your project in order to run task (build, test, deploy) by calling a link.


Installation
------------

```
go get github.com/idfly/darius
# or: sudo curl -o /bin/darius idfly.ru/darius/latest && sudo chmod oug+x /bin/darius
```


Usage
-----

Write simple configuration file:

```
# .darius.yml
tasks:
  say-hello: echo hello
```

Execute your task with command line:

```
darius say-hello
$ echo hello
  ! hello
```


Build
-----

  * clone repo
  * install latest `docker`, `docker-compose` and `darius`
  * execute `darius up` in order to install dependecies and run build
    containers
  * execute `darius build` in order to run tests


Authors
-------

  * [Leonid Shagabutdinov](http://github.com/shagabutdinov)
