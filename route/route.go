package route

import (
    "github.com/gorilla/mux"
    conversion "github.com/heaptracetechnology/microservice-ffmpeg/conversion"
    "log"
    "net/http"
)

type Route struct {
    Name        string
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
    Route{
        "VideoToImage",
        "POST",
        "/convertvideotoimage",
        conversion.VideoToImage,
    },
    Route{
        "Watermark",
        "POST",
        "/watermark",
        conversion.Watermark,
    },
}

func NewRouter() *mux.Router {
    router := mux.NewRouter().StrictSlash(true)
    for _, route := range routes {
        var handler http.Handler
        log.Println(route.Name)
        handler = route.HandlerFunc

        router.
            Methods(route.Method).
            Path(route.Pattern).
            Name(route.Name).
            Handler(handler)
    }
    return router
}
