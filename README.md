[![Build Status](https://travis-ci.org/motomux/glcp.svg?branch=master)](https://travis-ci.org/motomux/glcp)
[![Go Report Card](https://goreportcard.com/badge/github.com/motomux/glcp)](https://goreportcard.com/report/github.com/motomux/glcp)
# glcp
glcp is a helper tool for golint to put comment placeholder

## Installation

```
go get -u github.com/motomux/glcp
```

## Usage
Dry run
```
glcp $PACKAGE_NAME
```
  
Write change to source files
```
glcp -w $PACKAGE_NAME
```
  
Add comment placeholder recursively
```
glcp -w $(go list ./...)
```
