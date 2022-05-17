package utils

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	db "github.com/astralservices/api/supabase"
	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/nedpals/supabase-go"
	log "github.com/sirupsen/logrus"
)

func CORSMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("ENV") == "production" {
			handlers.CORS(handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}), handlers.AllowedHeaders([]string{"Authorization", "Content-Type"}), handlers.AllowedOrigins([]string{"*", "http://localhost:3000", "http://localhost:8000", "https://*.astralapp.io", "https://astralapp.io"}), handlers.AllowCredentials())(h)
		}

		h.ServeHTTP(w, r)
	})
}

func JSONMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	})
}

func AuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		database := db.New()
		userCookie, err := r.Cookie("access_token")
		authHeader := r.Header.Get("Authorization")

		var res []byte
		var authorization string

		if authHeader != "" {
			authorization = strings.Split(authHeader, " ")[1]
		} else {
			if err != nil {
				return
			}
			authorization = userCookie.Value
		}

		if authorization == "" {
			res, err = json.Marshal(Response[struct {
				Message string `json:"message"`
			}]{
				Result: struct {
					Message string "json:\"message\""
				}{Message: "You must be logged in to access this page!"},
				Code:  http.StatusUnauthorized,
				Error: err.Error(),
			})

			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}

			w.Write(res)

			return
		}

		user, err := database.Auth.User(r.Context(), authorization)

		if err != nil {
			res, err = json.Marshal(Response[struct {
				Message string `json:"message"`
			}]{
				Result: struct {
					Message string "json:\"message\""
				}{Message: "You must be logged in to access this page!"},
				Code:  http.StatusUnauthorized,
				Error: err.Error(),
			})

			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}

			w.Write(res)

			return
		}

		context.Set(r, "user", user)

		h.ServeHTTP(w, r)
	})
}

func ProfileMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		database := db.New()

		user := context.Get(r, "user").(*supabase.User)

		var profile []IProfile

		err := database.DB.From("profiles").Select("*").Eq("id", user.ID).Execute(&profile)

		if err != nil {
			w.Header().Add("Content-Type", "application/json")
			res, err := json.Marshal(Response[any]{
				Result: nil,
				Code:   http.StatusNotFound,
				Error:  err.Error(),
			})

			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}

			w.Write(res)
		}

		context.Set(r, "profile", profile[0])

		h.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		uri := r.RequestURI
		method := r.Method
		h.ServeHTTP(w, r) // serve the original request

		duration := time.Since(start)

		// log request details
		log.WithFields(log.Fields{
			"uri":      uri,
			"method":   method,
			"duration": duration,
		}).Info("Request")
	})
}

type String string

func (s String) Format(data map[string]interface{}) (out string, err error) {
	t := template.Must(template.New("").Parse(string(s)))
	builder := &strings.Builder{}
	if err = t.Execute(builder, data); err != nil {
		return
	}
	out = builder.String()
	return
}

func GetCallbackURL(provider string) string {
	callbackUrl := os.Getenv("CALLBACK_URL")

	tmpl, err := template.New("callbackUrl").Delims("[[", "]]").Parse(callbackUrl)
	if err != nil {
		log.Fatal(err)
	}

	data := map[string]interface{}{
		"Provider": provider,
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)

	if err != nil {
		log.Fatal(err)
	}

	return buf.String()
}

func RandomWord() string {
	words := []string{
		"aardvark",
		"albatross",
		"alligator",
		"alpaca",
		"ant",
		"anteater",
		"antelope",
		"ape",
		"armadillo",
		"donkey",
		"badger",
		"barracuda",
		"bat",
		"bear",
		"beaver",
		"bee",
		"bison",
		"boar",
		"buffalo",
		"butterfly",
		"camel",
		"capybara",
		"caribou",
		"cassowary",
		"cat",
		"caterpillar",
		"chamois",
		"cheetah",
		"chicken",
		"chimpanzee",
		"chinchilla",
		"clam",
		"cobra",
		"cockroach",
		"coyote",
		"crab",
		"crane",
		"crocodile",
		"crow",
		"curlew",
		"deer",
		"dog",
		"dolphin",
		"dove",
		"dragonfly",
		"duck",
		"eagle",
		"eel",
		"elephant",
		"elk",
		"emu",
		"falcon",
		"ferret",
		"finch",
		"fish",
		"flamingo",
		"fly",
		"fox",
		"frog",
		"gazelle",
		"gerbil",
		"giraffe",
		"goat",
		"goldfish",
		"goose",
		"gorilla",
		"grasshopper",
		"grouse",
		"guanaco",
		"gull",
		"hamster",
		"hare",
		"hawk",
		"hedgehog",
		"heron",
		"herring",
		"hippopotamus",
		"hornet",
		"horse",
		"hummingbird",
		"jackal",
		"jaguar",
		"jay",
		"jellyfish",
		"kangaroo",
		"kingfisher",
		"koala",
		"lemur",
		"leopard",
		"lion",
		"llama",
		"lobster",
		"magpie",
		"mallard",
		"manatee",
		"mandrill",
		"mantis",
		"marten",
		"meerkat",
		"mink",
		"mole",
		"mongoose",
		"monkey",
		"moose",
		"mosquito",
		"mouse",
		"mule",
		"narwhal",
		"newt",
		"nightingale",
		"octopus",
		"okapi",
		"oryx",
		"ostrich",
		"otter",
		"owl",
		"oyster",
		"panther",
		"parrot",
		"partridge",
		"peafowl",
		"pelican",
		"penguin",
		"pheasant",
		"pig",
		"pigeon",
		"pony",
		"porcupine",
		"quail",
		"rabbit",
		"raccoon",
		"ram",
		"rat",
		"raven",
		"red deer",
		"red panda",
		"reindeer",
		"rhinoceros",
		"rook",
		"salamander",
		"salmon",
		"sandpiper",
		"scorpion",
		"seahorse",
		"seal",
		"shark",
		"sheep",
		"shrew",
		"snake",
		"sparrow",
		"spider",
		"squid",
		"squirrel",
		"starling",
		"stingray",
		"stork",
		"swan",
		"tiger",
		"toad",
		"trout",
		"turkey",
		"turtle",
		"viper",
		"vulture",
		"wallaby",
		"walrus",
		"wasp",
		"weasel",
		"whale",
		"wildcat",
		"wolf",
		"wolverine",
		"wombat",
		"woodpecker",
		"worm",
		"wren",
		"yak",
		"zebra",
	}

	return words[RandInt(0, len(words)-1)]
}

func RandInt(min int, max int) int {
    return min + rand.Intn(max-min)
}