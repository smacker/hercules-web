package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	goCache "github.com/patrickmn/go-cache"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	hercules "gopkg.in/src-d/hercules.v5"
	"gopkg.in/src-d/hercules.v5/leaves"
)

var cache = goCache.New(time.Hour, 10*time.Minute)

func main() {
	exDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headers := w.Header()
			headers.Set("Access-Control-Allow-Origin", "*")
			headers.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			if r.Method == "OPTIONS" {
				return
			}
			h.ServeHTTP(w, r)
		})
	})

	externalHost := os.Getenv("EXTERNAL_HOST")
	if externalHost == "" {
		externalHost = "http://127.0.0.1:8080"
	}
	script := `<script>window.hercules = {apiHost: '` + externalHost + `'}</script>`
	indexHTML, err := ioutil.ReadFile(path.Join(exDir, "dist", "index.html"))
	if err != nil {
		panic(err)
	}
	indexHTML = bytes.Replace(indexHTML, []byte("</head>"), []byte(script+"</head>"), 1)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(indexHTML)
	})
	r.Mount("/static", http.FileServer(http.Dir(path.Join(exDir, "dist"))))

	r.Get("/api/burndown/*", jsonResponse(func(r *http.Request) (response, error) {
		repo := chi.URLParam(r, "*")
		return burndownCached("https://" + repo)
	}))

	fmt.Println("running...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

type respFunc func(*http.Request) (response, error)

func jsonResponse(f respFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := f(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(res.Status)
	}
}

type response struct {
	Status int         `json:"-"`
	Data   interface{} `json:"data"`
	Error  string      `json:"error"`
}

type burndownResponse struct {
	Begin      int64                `json:"begin"`
	End        int64                `json:"end"`
	Project    [][]int64            `json:"project"`
	Files      map[string][][]int64 `json:"filesData"`
	PeopleData [][][]int64          `json:"peopleData"`
	PeopleList []string             `json:"peopleList"`
}

func burndownCached(uri string) (response, error) {
	v, ok := cache.Get(uri)
	if ok {
		return v.(response), nil
	}

	res, err := burndown(uri)
	if err != nil {
		return response{}, err
	}
	cache.Set(uri, res, goCache.DefaultExpiration)

	return res, nil
}

func burndown(uri string) (response, error) {
	if err := validateRepo(uri); err != nil {
		return response{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}

	backend := memory.NewStorage()
	cloneOptions := &git.CloneOptions{URL: uri}
	repository, err := git.Clone(backend, nil, cloneOptions)
	if err != nil {
		// FIXME it can be internal error too
		return response{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}

	pipeline := hercules.NewPipeline(repository)
	commits, err := pipeline.Commits(false)
	if err != nil {
		return response{}, err
	}

	facts := map[string]interface{}{
		hercules.ConfigPipelineCommits: commits,
		// maybe move to another endpoint? but actually it's cheap enough compare to downloading repo
		leaves.ConfigBurndownGranularity: 30,
		leaves.ConfigBurndownSampling:    30,
		leaves.ConfigBurndownTrackPeople: true,
		leaves.ConfigBurndownTrackFiles:  true,
		// this constants are internal in hercules
		"RenameAnalysis.SimilarityThreshold": 80,
		"TreeDiff.Languages":                 []string{"all"},
		"TreeDiff.EnableBlacklist":           true,
		"TreeDiff.BlacklistedPrefixes":       []string{"vendor/", "vendors/", "node_modules/"},
	}

	burndownItem := hercules.Registry.Summon("Burndown")[0]
	pipeline.DeployItem(burndownItem)

	pipeline.Initialize(facts)

	results, err := pipeline.Run(commits)
	if err != nil {
		return response{}, err
	}
	// it's super ugly, but hercules api isn't very friendly or I just didn't get it
	var r leaves.BurndownResult
	for li, v := range results {
		if li == nil {
			continue
		}
		if li.Name() == "Burndown" {
			r = v.(leaves.BurndownResult)
		}
	}

	commonResult := results[nil].(*hercules.CommonAnalysisResult)

	return response{
		Status: http.StatusOK,
		Data: burndownResponse{
			Begin:      commonResult.BeginTime,
			End:        commonResult.EndTime,
			Project:    r.GlobalHistory,
			Files:      r.FileHistories,
			PeopleData: r.PeopleHistories,
			PeopleList: facts[hercules.FactIdentityDetectorReversedPeopleDict].([]string),
		},
	}, nil
}

const repoSizeLimit = 102400 // kb

func validateRepo(uri string) error {
	if !strings.HasPrefix(uri, "https://github.com/") {
		return errors.New("unsupported provider: only github is supported for now")
	}
	apiURI := strings.Replace(uri, "https://github.com/", "https://api.github.com/repos/", 1)
	resp, err := http.Get(apiURI)
	if err != nil {
		return fmt.Errorf("can't access github api: %s", err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return errors.New("repository not found")
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("can't access github api: %s", resp.Status)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("can't read github api response: %s", err)
	}
	var r struct{ Size int }
	if err := json.Unmarshal(b, &r); err != nil {
		return fmt.Errorf("can't parse github api response: %s", err)
	}
	if r.Size == 0 {
		return errors.New("incorrect repository")
	}
	if r.Size > repoSizeLimit {
		return errors.New("repository is too big")
	}
	return nil
}
