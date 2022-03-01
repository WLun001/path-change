package main

import (
	"encoding/json"
	"fmt"
	"github.com/bmatcuk/doublestar/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/go-homedir"
	"github.com/tektoncd/triggers/pkg/apis/triggers/v1beta1"
	"google.golang.org/grpc/codes"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	clonePath         = "./repo"
	defaultConfigFile = "config.yaml"
	refParam          = "ref"
	repoParam         = "repo"
)

var (
	EmptyPattern = "EMPTY_PATTERN"
	Match        = "MATCH"
	NotMatch     = "NOT_MATCH"
)

type repos struct {
	Repos map[string]config `yaml:"repos"`
}

type config struct {
	URL   string   `yaml:"url"`
	Paths []string `yaml:"paths,omitempty"`
}

// todo: git credentials
// update examples and readme

func fileChange(ctx *fiber.Ctx, repos *repos, repo, branch string) error {
	var url string
	var paths []string
	val, ok := repos.Repos[repo]
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "repo not exist in config file")
	}
	url = val.URL
	paths = val.Paths

	paths = removeEmptyStrings(paths)
	if len(paths) <= 0 {
		return ctx.JSON(v1beta1.InterceptorResponse{
			Extensions: fiber.Map{"paths": EmptyPattern},
			Continue:   true,
		})
	}

	// mkdir
	if _, err := os.Stat(clonePath); os.IsNotExist(err) {
		err := os.Mkdir(clonePath, 0700)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprint(err))

		}
	}

	homepath, err := homedir.Dir()
	if err != nil {
		log.Printf("Unexpected error getting the user home directory: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprint(err))
	}
	if os.Geteuid() == 0 {
		homepath = "/root"
	}

	ensureHomeEnv(homepath)

	// git clone
	log.Printf("cloning %s at branch %s", repo, branch)
	clone := exec.Command("git", "clone", "-b", branch, "--depth", "2", url, clonePath)
	//err = clone.Run()
	out, err := clone.CombinedOutput()
	log.Println(string(out))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprint(err))
	}

	// git diff
	err = os.Chdir(clonePath)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprint(err))
	}
	cmd := exec.Command("git", "--no-pager", "diff", "--name-only", "HEAD^")
	stdout, _ := cmd.CombinedOutput()

	// remove dir
	err = os.Chdir("../")
	err = os.RemoveAll(clonePath)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprint(err))

	}

	// process output
	output := strings.Split(strings.TrimSpace(string(stdout)), "\n")
	log.Printf("file change %s:", output)
	for _, e := range output {
		for _, path := range paths {
			m, err := doublestar.PathMatch(path, e)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprint(err))
			}
			if m {
				return ctx.JSON(v1beta1.InterceptorResponse{
					Extensions: fiber.Map{"paths": Match},
					Continue:   true,
				})
			}
		}
	}
	return ctx.JSON(v1beta1.InterceptorResponse{
		Extensions: fiber.Map{"paths": NotMatch},
		Continue:   false,
	})
}

func main() {

	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = defaultConfigFile
	}

	// read config.yaml
	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}
	var repos *repos
	err = yaml.Unmarshal(yamlFile, &repos)
	if err != nil {
		panic(err)
	}

	app := fiber.New(fiber.Config{ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		// Retrieve the custom status code if it's a fiber.*Error
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}
		return ctx.Status(code).JSON(v1beta1.InterceptorResponse{
			Continue: false,
			Status: v1beta1.Status{
				Code:    codes.Internal,
				Message: err.Error(),
			},
		})
	}})

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("hello world")
	})

	app.Post("/local", func(ctx *fiber.Ctx) error {
		repo := ctx.Query("repo")
		ref := ctx.Query("ref")
		branch := getBranch(ref)
		return fileChange(ctx, repos, repo, branch)
	})

	app.Post("/", func(ctx *fiber.Ctx) error {
		byteBody := ctx.Body()
		req := new(v1beta1.InterceptorRequest)
		err = json.Unmarshal(byteBody, &req)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		repo := fmt.Sprint(req.InterceptorParams[repoParam])
		ref := fmt.Sprint(req.Extensions[refParam])
		branch := getBranch(ref)
		return fileChange(ctx, repos, repo, branch)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	log.Fatal(app.Listen(port))
}
