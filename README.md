# jorum [![GoDoc](https://godoc.org/github.com/zignd/jorum?status.svg)](https://godoc.org/github.com/zignd/jorum) [![Report card](https://goreportcard.com/badge/github.com/zignd/jorum)](https://goreportcard.com/report/github.com/zignd/jorum)     

A large drinking vessel into which you can group services. But here's the original definition of [jorum](https://en.wiktionary.org/wiki/jorum).

Before we continue, we need to agree on a definition of service, from now on, by service, it is meant something with the following characteristics: performs IO, initializes and stops, may emit events throughout its lifecycle which could indicate errors, warnings or general information.

It often occurs that an application is required to manage the lifecycle of multiple services; those services may be HTTP servers, Apache Kafka consumers, database connections, etc. Handling of those services usually requires repetitive code, which is typically boilerplate that will be repeated the next time an application with the same services ends up being created.

Jorum aims to help you handle that by literally providing a large vessel into which you can group services and manage their similarities from a single control panel.
