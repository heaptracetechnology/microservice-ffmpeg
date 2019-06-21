# _FFmpeg_ as a microservice

An OMG service for FFmpeg, it allows encode/decode, converts audio or video formats. It can also capture and encode in real-time from various hardware and software sources 

[![Open Microservice Guide](https://img.shields.io/badge/OMG%20Enabled-👍-green.svg?)](https://microservice.guide)

## Direct usage in [Storyscript](https://storyscript.io/):

##### Convert Video to Image
```coffee
>>> ffmpeg convertVideoToImage videoBase64:'Base64 data'

```

Curious to [learn more](https://docs.storyscript.io/)?

✨🍰✨

## Usage with [OMG CLI](https://www.npmjs.com/package/omg)

##### Convert Video to Image
```shell
$ omg run convertVideoToImage -a videoBase64=<BASE64_DATA>
```

**Note**: The OMG CLI requires [Docker](https://docs.docker.com/install/) to be installed.

## License
[MIT License](https://github.com/omg-services/ffmpeg/blob/master/LICENSE).
