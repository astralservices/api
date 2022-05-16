package utils

import (
	"encoding/json"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	db "github.com/astralservices/api/supabase"
	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/nedpals/supabase-go"
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

	s, err := String(callbackUrl).Format(map[string]interface{}{
		"Provider": provider,
	})

	if err != nil {
		return ""
	}

	return s
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

	rand.Seed(time.Now().UTC().UnixNano())

	return words[rand.Intn(len(words))]
}
