# EZProxy

EZProxy (read as "easy proxy") is an open source solution to store your own Go modules in a repository at your __own VDS__.

## When it to use?
In a case when you want to use your own Go modules in your different projects you can:
- publish a module at Github
- copy&paste a module as a package from one project to next
- use Git submodules
- use Go proxy such as Athens or __EZProxy__

## How EZProxy works?
![](//6wings-tech/go-ezproxy/blob/master/docs/pic-publish.png?raw=true)

![](//6wings-tech/go-ezproxy/blob/master/docs/pic-download.png?raw=true)

## Prerequisites
1. VDS
2. Domain name (for example, *git.example.com*) binded to your VDS
3. Git & Go lang at your PC
4. Git at VDS
