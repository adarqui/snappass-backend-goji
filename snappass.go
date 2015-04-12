package main

import (
	"flag"
	"fmt"
	"github.com/adarqui/snappass-core-go"
	"github.com/hypebeast/gojistatic"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"log"
	"net/http"
	"os"
)

var prog = "snappass-backend-goji"

func usage() {
	log.Println("usage: ./snappass-backend-goji <path/to/config.json>")
	os.Exit(1)
}

func getPassword(snap *snappass_core.SnapPass, c web.C, w http.ResponseWriter, r *http.Request) {
	key := c.URLParams["key"]
	password, err := snap.GetPassword([]byte(key))
	if err != nil {
		http.Error(w, "Key not found.", 400)
	} else {
		fmt.Fprintf(w, "%s", password)
	}
}

func postPassword(snap *snappass_core.SnapPass, c web.C, w http.ResponseWriter, r *http.Request) {
	password := c.URLParams["password"]
	ttl := c.URLParams["ttl"]
	key, err := snap.SetPasswordStrTTL([]byte(password), ttl)
	if err != nil {
		http.Error(w, "Unexpected Error", 400)
	} else {
		fmt.Fprintf(w, "%s", key)
	}
}

func main() {

	args := os.Args
	if len(args) < 2 {
		usage()
	}

	parsed_config, err_config := snappass_core.ParseConfig(args[1])
	if err_config != nil {
		log.Fatal(prog+": Failure while attempting to parse config file: ", err_config)
	}

	flag.Set("bind", fmt.Sprintf("%s:%d", parsed_config.Web.Host, parsed_config.Web.Port))

	redisdb, err_redis := snappass_core.NewRedisDatabase(fmt.Sprintf("%s:%d", parsed_config.Redis.Host, parsed_config.Redis.Port), parsed_config.Redis.Auth, parsed_config.Redis.Db)
	if err_redis != nil {
		log.Fatal(prog+": Failure while attempting to connect to redis: ", err_redis)
	}
	uuidkg, err_kg := snappass_core.NewUUIDKeyGenerator([]byte("snap:"))
	if err_kg != nil {
		log.Fatal(prog+": Unable to create uuid key generator: ", err_kg)
	}

	snap, err := snappass_core.New(redisdb, uuidkg)
	if err != nil {
		log.Fatal(prog+": Failure creating snappass service: ", err)
	}

	goji.Use(gojistatic.Static(parsed_config.Web.Static, gojistatic.StaticOptions{SkipLogging: true}))
	goji.Get("/key/:key", func(c web.C, w http.ResponseWriter, r *http.Request) { getPassword(snap, c, w, r) })
	goji.Post("/pass/:password/:ttl", func(c web.C, w http.ResponseWriter, r *http.Request) { postPassword(snap, c, w, r) })
	goji.Serve()
}
