package controller

import (
	"github.com/valyala/fasthttp"
	alg "app/model/algorithm"
	"strconv"
	"flag"
	"fmt"
	"log"
)

var (
	addr     = flag.String("addr", ":8080", "TCP address to listen to")
	compress = flag.Bool("compress", false, "Whether to enable transparent response compression")
)


func recommendHandler(ctx *fasthttp.RequestCtx) {
	domain := string(ctx.FormValue("domain"))
	boxidStr := string(ctx.FormValue("boxid"))
	guid := string(ctx.FormValue("guid"))
	itemid := string(ctx.FormValue("itemid"))

	ctx.SetContentType("application/json; charset=utf8")
	boxid, _ := strconv.Atoi(boxidStr)

	data := alg.GetRecommendNews(domain, boxid, guid, itemid)
	fmt.Fprintln(ctx, data)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func Load() {
	flag.Parse()
	h := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/recommend":
			recommendHandler(ctx)
		default:
			ctx.Error("not found", fasthttp.StatusNotFound)
		}
	}

	log.Println("Start server at 127.0.0.1:8080")
	if err := fasthttp.ListenAndServe(*addr, h); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}
