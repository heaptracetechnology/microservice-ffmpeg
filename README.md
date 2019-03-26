# FFmpeg as a microservice
An OMG service for FFmpeg, it allows encode/decode, converts audio or video formats. It can also capture and encode in real-time from various hardware and software sources 

[![Open Microservice Guide](https://img.shields.io/badge/OMG-enabled-brightgreen.svg?style=for-the-badge)](https://microservice.guide)
<!-- [![Build Status](https://travis-ci.com/heaptracetechnology/microservice-urbanairship.svg?branch=master)](https://travis-ci.com/heaptracetechnology/microservice-urbanairship)
[![codecov](https://codecov.io/gh/heaptracetechnology/microservice-urbanairship/branch/master/graph/badge.svg)](https://codecov.io/gh/heaptracetechnology/microservice-urbanairship)
[![GolangCI](https://golangci.com/badges/github.com/golangci/golangci-web.svg)](https://golangci.com) -->


## [OMG](hhttps://microservice.guide) CLI

### OMG

* omg validate
```
omg validate
```
* omg build
```
omg build
```
### Test Service

* Test the service by following OMG commands

### CLI

##### Convert Video to Image
```sh
$ omg run convert_video_to_image -a video_base64=<BASE64_DATA> -a video_extension=<VIDEO_EXTENSION> -a image_extension=<IMAGE_EXTENSION>
```

## License
### [MIT](https://choosealicense.com/licenses/mit/)

## Docker
### Build
```
docker build -t microservice-ffmpeg .
```
### RUN
```
docker run -p 3000:3000 microservice-ffmpeg
```
