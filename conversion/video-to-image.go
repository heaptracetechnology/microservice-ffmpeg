package conversion

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/3d0c/gmf"
	"github.com/heaptracetechnology/microservice-ffmpeg/result"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	extention   string
	format      string
	fileCount   int
	srcFileName string
	swsctx      *gmf.SwsCtx
)

type ArgumentData struct {
	VideoBase64     string `json:"video_base64"`
	InputExtension  string `json:"input_extension"`
	OutputExtension string `json:"output_extension"`
}

type Message struct {
	Success    string `json:"success"`
	Message    string `json:"message"`
	StatusCode int    `json:"statuscode"`
}

//Convert video to images
func VideoToImage(responseWriter http.ResponseWriter, request *http.Request) {

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		result.WriteErrorResponse(responseWriter, err)
		return
	}
	defer request.Body.Close()

	var argumentData ArgumentData
	unmarshalErr := json.Unmarshal(body, &argumentData)
	if unmarshalErr != nil {
		result.WriteErrorResponse(responseWriter, unmarshalErr)
		return
	}

	data, decodeErr := base64.StdEncoding.DecodeString(argumentData.VideoBase64)
	if decodeErr != nil {
		result.WriteErrorResponse(responseWriter, decodeErr)
		return
	}
	filename := time.Now().String() + ".3gp"
	filepath := os.TempDir() + "/" + filename
	f, createFileErr := os.Create(filepath)
	if createFileErr != nil {
		result.WriteErrorResponse(responseWriter, createFileErr)
		return
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		panic(err)
	}
	if err := f.Sync(); err != nil {
		panic(err)
	}

	flag.StringVar(&srcFileName, "src", filename, "source video")
	flag.StringVar(&extention, "ext", "png", "destination type, e.g.: png, jpg, tiff, whatever encoder you have")
	flag.Parse()

	os.MkdirAll("./tmp", 0755)

	inputCtx, err := gmf.NewInputCtx(srcFileName)
	if err != nil {
		log.Fatalf("Error creating context - %s\n", err)
	}
	defer inputCtx.Free()

	srcVideoStream, err := inputCtx.GetBestStream(gmf.AVMEDIA_TYPE_VIDEO)
	if err != nil {
		log.Printf("No video stream found in '%s'\n", srcFileName)
		return
	}

	codec, err := gmf.FindEncoder(extention)
	if err != nil {
		log.Fatalf("%s\n", err)
	}

	cc := gmf.NewCodecCtx(codec)
	defer gmf.Release(cc)

	cc.SetTimeBase(gmf.AVR{Num: 1, Den: 1})

	cc.SetPixFmt(gmf.AV_PIX_FMT_RGBA).SetWidth(srcVideoStream.CodecCtx().Width()).SetHeight(srcVideoStream.CodecCtx().Height())
	if codec.IsExperimental() {
		cc.SetStrictCompliance(gmf.FF_COMPLIANCE_EXPERIMENTAL)
	}

	if err := cc.Open(nil); err != nil {
		log.Fatal(err)
	}
	defer cc.Free()

	ist, err := inputCtx.GetStream(srcVideoStream.Index())
	if err != nil {
		log.Fatalf("Error getting stream - %s\n", err)
	}
	defer ist.Free()

	icc := srcVideoStream.CodecCtx()
	if swsctx, err = gmf.NewSwsCtx(icc.Width(), icc.Height(), icc.PixFmt(), cc.Width(), cc.Height(), cc.PixFmt(), gmf.SWS_BICUBIC); err != nil {
		panic(err)
	}
	defer swsctx.Free()

	ln := int(math.Log10(float64(ist.NbFrames()))) + 1
	format = "./tmp/" + "%0" + strconv.Itoa(ln) + "d." + extention

	start := time.Now()

	var (
		pkt        *gmf.Packet
		frames     []*gmf.Frame
		drain      int = -1
		frameCount int = 0
	)

	for {
		if drain >= 0 {
			break
		}

		pkt, err = inputCtx.GetNextPacket()
		if err != nil && err != io.EOF {
			if pkt != nil {
				pkt.Free()
			}
			log.Printf("error getting next packet - %s", err)
			break
		} else if err != nil && pkt == nil {
			drain = 0
		}

		if pkt != nil && pkt.StreamIndex() != srcVideoStream.Index() {
			continue
		}

		frames, err = ist.CodecCtx().Decode(pkt)
		if err != nil {
			log.Printf("Fatal error during decoding - %s\n", err)
			break
		}

		if len(frames) == 0 && drain < 0 {
			continue
		}

		if frames, err = gmf.DefaultRescaler(swsctx, frames); err != nil {
			panic(err)
		}

		encode(cc, frames, drain)

		for i, _ := range frames {
			frames[i].Free()
			frameCount++
		}

		if pkt != nil {
			pkt.Free()
			pkt = nil
		}
	}

	for i := 0; i < inputCtx.StreamsCnt(); i++ {
		st, _ := inputCtx.GetStream(i)
		st.CodecCtx().Free()
		st.Free()
	}

	since := time.Since(start)
	log.Printf("Finished in %v, avg %.2f fps", since, float64(frameCount)/since.Seconds())

	message := Message{"true", "Converted video to images", http.StatusOK}
	bytes, _ := json.Marshal(message)
	result.WriteJsonResponse(responseWriter, bytes, http.StatusOK)

}

func encode(cc *gmf.CodecCtx, frames []*gmf.Frame, drain int) {
	packets, err := cc.Encode(frames, drain)
	if err != nil {
		log.Fatalf("Error encoding - %s\n", err)
	}
	if len(packets) == 0 {
		return
	}

	for _, p := range packets {
		writeFile(p.Data())
		p.Free()
	}

	return
}

func writeFile(b []byte) {
	name := fmt.Sprintf(format, fileCount)

	fp, err := os.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("%s\n", err)
	}

	if n, err := fp.Write(b); err != nil {
		log.Fatalf("%s\n", err)
	} else {
		log.Printf("%d bytes written to '%s'", n, name)
	}

	fp.Close()
	fileCount++
}
