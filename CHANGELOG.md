# Changelog

## 1.12.2
### New Features
- Add GET, POST, PUT, DELETE requests between microservices

## 1.12.1

### Bugs Fixes
- Issue with wrong path Request
- Issue with sessions
- other small fixes

### Improvemetns
- Fixed issue with file upload
- Remove Request struct

## 1.12.0

### Core Improvemetns
- Remove plugins support
- Change Error response structure
- Remove API v2 suport
- Proto modifications. Preparations for streaming connection

## 1.11.6

### Bug Fixing
- Fixed issue with wrong token initialisation

## 1.11.5

### Gufo functions

- Add ErrorReturn Function as general Error Handler
- Add CheckForSign and Gufo Sign. Gufo Sign is necessary for safety connection between GRPC microservices. By this sign, Microservice can understand that request was made from right GUFO instance

### proto file (GRPC Request)

- Add Sign data with GUFO Sign
- Add Requestor's IP address and UserAgent Data
