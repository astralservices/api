package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strings"
	"text/template"
	"time"

	db "github.com/astralservices/api/supabase"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/handlers"
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

func AuthMiddleware(ctx *fiber.Ctx) error {
	auth_header := ctx.GetReqHeaders()["Authorization"]
	auth_cookie := ctx.Cookies("token")
	if (auth_header != "" && !strings.HasPrefix(auth_header, "Bearer")) || auth_cookie == "" {
		return ctx.Status(http.StatusUnauthorized).JSON(Response[struct {
			Message string `json:"message"`
		}]{
			Result: struct {
				Message string "json:\"message\""
			}{Message: "You must be logged in to access this page!"},
			Code:  http.StatusUnauthorized,
			Error: "",
		})
	}

	var tokenString string

	if auth_header != "" {
		tokenString = strings.TrimPrefix(auth_header, "Bearer ")
	} else if auth_cookie != "" {
		tokenString = auth_cookie
	}

	claims, err := GetClaimsFromToken(tokenString)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(Response[struct {
			Message string `json:"message"`
		}]{
			Result: struct {
				Message string "json:\"message\""
			}{Message: "There was an error while trying to authenticate you. Please try again."},
			Code:  http.StatusUnauthorized,
			Error: "",
		})
	}

	ctx.Locals("user", claims.UserInfo)

	return ctx.Next()
}

// Injects user if the user exists
func AuthInjectorMiddleware(ctx *fiber.Ctx) error {
	auth_header := ctx.GetReqHeaders()["Authorization"]
	auth_cookie := ctx.Cookies("token")
	if (auth_header != "" && !strings.HasPrefix(auth_header, "Bearer")) || auth_cookie == "" {
		return ctx.Next()
	}

	var tokenString string

	if auth_header != "" {
		tokenString = strings.TrimPrefix(auth_header, "Bearer ")
	} else if auth_cookie != "" {
		tokenString = auth_cookie
	}

	claims, err := GetClaimsFromToken(tokenString)
	if err != nil {
		return ctx.Next()
	}

	ctx.Locals("user", claims.UserInfo)

	return ctx.Next()
}

func WorkspaceMiddleware(ctx *fiber.Ctx) error {
	workspace_id := ctx.Params("workspace_id")

	if workspace_id == "" {
		return ctx.Status(http.StatusBadRequest).JSON(Response[struct {
			Message string `json:"message"`
		}]{
			Result: struct {
				Message string "json:\"message\""
			}{Message: "You must provide a workspace ID!"},
			Code:  http.StatusBadRequest,
			Error: "",
		})
	}

	database := db.New()

	var workspace IWorkspace

	err := database.DB.From("workspaces").Select("*").Single().Eq("id", workspace_id).Execute(&workspace)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(Response[struct {
			Message string `json:"message"`
		}]{
			Result: struct {
				Message string "json:\"message\""
			}{Message: "There was an error while trying to find the workspace!"},
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
	}

	ctx.Locals("workspace", workspace)

	return ctx.Next()
}

func WorkspaceMemberMiddleware(ctx *fiber.Ctx) error {
	workspace := ctx.Locals("workspace").(IWorkspace)

	user := ctx.Locals("user").(IProvider)

	database := db.New()

	var workspaceMember IWorkspaceMember

	err := database.DB.From("workspace_members").Select("id, created_at, role, pending").Single().Eq("workspace", *workspace.ID).Eq("profile", *user.ID).Execute(&workspaceMember)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(Response[struct {
			Message string `json:"message"`
		}]{
			Result: struct {
				Message string "json:\"message\""
			}{Message: "There was an error while trying to find the authenticated workspace member!"},
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
	}

	ctx.Locals("workspace_member", workspaceMember)

	return ctx.Next()
}

func ProfileMiddleware(ctx *fiber.Ctx) error {
	database := db.New()

	user := ctx.Locals("user").(IProvider)

	var profile []IProfile

	err := database.DB.From("profiles").Select("*").Eq("id", *user.ID).Execute(&profile)

	if err != nil {
		return ctx.JSON(Response[any]{
			Result: nil,
			Code:   http.StatusNotFound,
			Error:  err.Error(),
		})
	}

	if len(profile) == 0 {
		return ctx.JSON(Response[any]{
			Result: nil,
			Code:   http.StatusNotFound,
			Error:  "Profile not found",
		})
	}

	ctx.Locals("profile", profile[0])

	return ctx.Next()
}

func BotMiddleware(ctx *fiber.Ctx) error {
	workspace := ctx.Locals("workspace").(IWorkspace)

	var bots []IBot

	database := db.New()

	err := database.DB.From("bots").Select("id, created_at, owner, region, settings, token, commands").Eq("workspace", *workspace.ID).Execute(&bots)

	if err != nil {
		return ctx.Status(500).JSON(Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	if len(bots) == 0 {
		return ctx.Status(404).JSON(Response[any]{
			Result: nil,
			Code:   http.StatusNotFound,
			Error:  "Bot not found",
		})
	}

	bot := bots[0]

	ctx.Locals("bot", bot)

	return ctx.Next()
}

func WorkspaceIntegrationMiddleware(ctx *fiber.Ctx) error {
	integrationId := ctx.Params("integrationId")
	workspace := ctx.Locals("workspace").(IWorkspace)

	var integrations []IWorkspaceIntegration

	database := db.New()

	fmt.Println(integrationId, *workspace.ID)

	err := database.DB.From("workspace_integrations").Select("*").Eq("workspace", *workspace.ID).Eq("integration", integrationId).Execute(&integrations)

	if err != nil {
		return ErrorResponse(ctx, 500, err.Error())
	}

	fmt.Printf("%+v\n", integrations)

	if len(integrations) == 0 {
		return ErrorResponse(ctx, 404, "Integration not found")
	}

	ctx.Locals("integration", integrations[0])

	return ctx.Next()
}

func ErrorResponse(ctx *fiber.Ctx, code int, message string) error {
	redirect := ctx.FormValue("redirect")

	if redirect != "" {
		return ctx.Redirect(redirect + "?error=" + message)
	}

	return ctx.Status(500).JSON(Response[any]{
		Result: nil,
		Code:   code,
		Error:  message,
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

type UserClaims struct {
	UserInfo IProvider
	*jwt.RegisteredClaims
}

var secret = []byte(os.Getenv("SECRET"))

func CreateToken(sub string, userInfo IProvider) (string, error) {
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	exp := time.Now().Add(time.Hour * 24)
	token.Claims = &UserClaims{
		userInfo,
		&jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			Subject:   sub,
		},
	}
	val, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}
	return val, nil
}

func GetClaimsFromToken(tokenString string) (UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("unexpected signing method ", token.Header["alg"])
			return nil, jwt.ErrSignatureInvalid
		}
		return secret, nil
	})
	if err != nil {
		return UserClaims{}, err
	}

	claims := token.Claims.(*UserClaims)
	ok := token.Valid

	if ok {
		return *claims, nil
	}
	return UserClaims{}, err
}

type OrderedMap struct {
	Order []string
	Map   map[string]string
}

func (om *OrderedMap) UnmarshalJSON(b []byte) error {
	json.Unmarshal(b, &om.Map)

	index := make(map[string]int)
	for key := range om.Map {
		om.Order = append(om.Order, key)
		esc, _ := json.Marshal(key) //Escape the key
		index[key] = bytes.Index(b, esc)
	}

	sort.Slice(om.Order, func(i, j int) bool { return index[om.Order[i]] < index[om.Order[j]] })
	return nil
}

func (om OrderedMap) MarshalJSON() ([]byte, error) {
	var b []byte
	buf := bytes.NewBuffer(b)
	buf.WriteRune('{')
	l := len(om.Order)
	for i, key := range om.Order {
		km, err := json.Marshal(key)
		if err != nil {
			return nil, err
		}
		buf.Write(km)
		buf.WriteRune(':')
		vm, err := json.Marshal(om.Map[key])
		if err != nil {
			return nil, err
		}
		buf.Write(vm)
		if i != l-1 {
			buf.WriteRune(',')
		}
		fmt.Println(buf.String())
	}
	buf.WriteRune('}')
	fmt.Println(buf.String())
	return buf.Bytes(), nil
}
