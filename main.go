package main

import (
	"bytes"
	"encoding/json"
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
	flags "github.com/jessevdk/go-flags"
	goCache "github.com/patrickmn/go-cache"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	hercules "gopkg.in/src-d/hercules.v5"
	"gopkg.in/src-d/hercules.v5/leaves"
)

type options struct {
	Storage       string `short:"s" long:"storage" env:"STORAGE" default:"memory" choice:"memory" choice:"disk" description:"store backend for analysis results"`
	DiskDir       string `long:"disk-storage-dir" env:"DISK_STORAGE_DIR" default:"/tmp" description:"directory for disk storage"`
	RepoSizeLimit int    `long:"repository-size-limit" env:"REPOSITORY_SIZE_LIMIT" default:"102400" description:"reject repositories bigger than (kb)"`
}

var opts options
var parser = flags.NewParser(&opts, flags.Default)

func main() {
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	var storage storage
	switch opts.Storage {
	case "memory":
		storage = newCachedStorage()
	case "disk":
		var err error
		storage, err = newDiskStorage(opts.DiskDir)
		if err != nil {
			log.Fatal(err)
		}
	}

	static := newStaticServer()
	api := newAPIServer(storage)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware)
	r.Use(middleware.DefaultCompress)

	r.Mount("/", static.Router())
	r.Mount("/api/", api.Router())

	log.Println("running...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func memClone(uri string) (*git.Repository, error) {
	backend := memory.NewStorage()
	cloneOptions := &git.CloneOptions{URL: uri}
	return git.Clone(backend, nil, cloneOptions)
}

type herculesResponse struct {
	Begin      int64                `json:"begin"`
	End        int64                `json:"end"`
	Project    [][]int64            `json:"project"`
	Files      map[string][][]int64 `json:"filesData"`
	PeopleData [][][]int64          `json:"peopleData"`
	PeopleList []string             `json:"peopleList"`
}

func herculesRun(repository *git.Repository) (*herculesResponse, error) {
	pipeline := hercules.NewPipeline(repository)
	commits, err := pipeline.Commits(false)
	if err != nil {
		return nil, err
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
		return nil, err
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

	return &herculesResponse{
		Begin:      commonResult.BeginTime,
		End:        commonResult.EndTime,
		Project:    r.GlobalHistory,
		Files:      r.FileHistories,
		PeopleData: r.PeopleHistories,
		PeopleList: facts[hercules.FactIdentityDetectorReversedPeopleDict].([]string),
	}, nil
}

type validationError struct {
	msg string
}

func (e *validationError) Error() string {
	return e.msg
}

func newValidationError(format string, args ...interface{}) *validationError {
	return &validationError{msg: fmt.Sprintf(format, args...)}
}

func validateRepo(uri string) error {
	if !strings.HasPrefix(uri, "https://github.com/") {
		return newValidationError("unsupported provider: only github is supported for now")
	}
	apiURI := strings.Replace(uri, "https://github.com/", "https://api.github.com/repos/", 1)
	resp, err := http.Get(apiURI)
	if err != nil {
		return newValidationError("can't access github api: %s", err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return newValidationError("repository not found")
	}
	if resp.StatusCode != http.StatusOK {
		return newValidationError("can't access github api: %s", resp.Status)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return newValidationError("can't read github api response: %s", err)
	}
	var r struct{ Size int }
	if err := json.Unmarshal(b, &r); err != nil {
		return newValidationError("can't parse github api response: %s", err)
	}
	if r.Size == 0 {
		return newValidationError("incorrect repository")
	}
	if r.Size > opts.RepoSizeLimit {
		return newValidationError("repository is too big")
	}
	return nil
}

//

type storage interface {
	BurndownProject(uri string) (*burndownProjectResp, error)
	BurndownPeople(uri string) (*burndownPeopleResp, error)
	BurndownFiles(uri string) (*burndownFilesResp, error)
}

type burndownResp struct {
	Begin int64 `json:"begin"`
	End   int64 `json:"end"`
}

func toBurndownResp(res *herculesResponse) burndownResp {
	return burndownResp{
		Begin: res.Begin,
		End:   res.End,
	}
}

type burndownProjectResp struct {
	burndownResp
	Project [][]int64 `json:"project"`
}

type burndownPeopleResp struct {
	burndownResp
	PeopleData [][][]int64 `json:"peopleData"`
	PeopleList []string    `json:"peopleList"`
}

type burndownFilesResp struct {
	burndownResp
	Files map[string][][]int64 `json:"filesData"`
}

type cachedStorage struct {
	cache *goCache.Cache
}

var _ storage = &cachedStorage{}

func newCachedStorage() *cachedStorage {
	return &cachedStorage{cache: goCache.New(6*time.Hour, time.Hour)}
}

func (s *cachedStorage) BurndownProject(uri string) (*burndownProjectResp, error) {
	data, err := s.cached(uri)
	if err != nil {
		return nil, err
	}

	return &burndownProjectResp{
		burndownResp: toBurndownResp(data),
		Project:      data.Project,
	}, nil
}

func (s *cachedStorage) BurndownPeople(uri string) (*burndownPeopleResp, error) {
	data, err := s.cached(uri)
	if err != nil {
		return nil, err
	}

	return &burndownPeopleResp{
		burndownResp: toBurndownResp(data),
		PeopleData:   data.PeopleData,
		PeopleList:   data.PeopleList,
	}, nil
}

func (s *cachedStorage) BurndownFiles(uri string) (*burndownFilesResp, error) {
	data, err := s.cached(uri)
	if err != nil {
		return nil, err
	}

	return &burndownFilesResp{
		burndownResp: toBurndownResp(data),
		Files:        data.Files,
	}, nil
}

func (s *cachedStorage) cached(uri string) (*herculesResponse, error) {
	v, ok := s.cache.Get(uri)
	if ok {
		return v.(*herculesResponse), nil
	}

	if err := validateRepo(uri); err != nil {
		return nil, err
	}

	repo, err := memClone(uri)
	if err != nil {
		return nil, err
	}

	res, err := herculesRun(repo)
	if err != nil {
		return nil, err
	}

	s.cache.Set(uri, res, goCache.DefaultExpiration)

	return res, nil
}

type diskStorage struct {
	root string
}

func newDiskStorage(root string) (*diskStorage, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("'%s' must be a directory", root)
	}

	cacheRoot := path.Join(root, "hercules-cache")
	if err := os.MkdirAll(cacheRoot, os.ModePerm); err != nil {
		return nil, err
	}

	return &diskStorage{root: cacheRoot}, nil
}

func (s *diskStorage) BurndownProject(uri string) (*burndownProjectResp, error) {
	data, err := s.cached(uri)
	if err != nil {
		return nil, err
	}

	return &burndownProjectResp{
		burndownResp: toBurndownResp(data),
		Project:      data.Project,
	}, nil
}

func (s *diskStorage) BurndownPeople(uri string) (*burndownPeopleResp, error) {
	data, err := s.cached(uri)
	if err != nil {
		return nil, err
	}

	return &burndownPeopleResp{
		burndownResp: toBurndownResp(data),
		PeopleData:   data.PeopleData,
		PeopleList:   data.PeopleList,
	}, nil
}

func (s *diskStorage) BurndownFiles(uri string) (*burndownFilesResp, error) {
	data, err := s.cached(uri)
	if err != nil {
		return nil, err
	}

	return &burndownFilesResp{
		burndownResp: toBurndownResp(data),
		Files:        data.Files,
	}, nil
}

func (s *diskStorage) cached(uri string) (*herculesResponse, error) {
	resultFile := path.Join(s.root, uriToFilename(uri))

	info, err := os.Stat(resultFile)
	if err == nil {
		if info.IsDir() {
			return nil, fmt.Errorf("'%s' is a directory but should be file", resultFile)
		}

		b, err := ioutil.ReadFile(resultFile)
		if err != nil {
			return nil, err
		}

		var res herculesResponse
		if err = json.Unmarshal(b, &res); err != nil {
			return nil, err
		}

		return &res, err
	}
	if os.IsNotExist(err) {
		if err := validateRepo(uri); err != nil {
			return nil, err
		}

		repo, err := memClone(uri)
		if err != nil {
			return nil, err
		}

		res, err := herculesRun(repo)
		if err != nil {
			return nil, err
		}

		f, err := os.Create(resultFile)
		if err != nil {
			return nil, err
		}

		b, err := json.Marshal(res)
		if err != nil {
			return nil, err
		}

		_, err = f.Write(b)
		if err != nil {
			return nil, err
		}

		return res, nil
	}

	return nil, err
}

func uriToFilename(uri string) string {
	charsToReplace := []string{"/", ":", "@", "#"}

	for _, c := range charsToReplace {
		uri = strings.Replace(uri, c, "_", -1)
	}

	return uri
}

// http staff

type staticServer struct {
	exDir     string
	indexHTML []byte
}

func newStaticServer() *staticServer {
	exDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

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

	return &staticServer{
		exDir:     exDir,
		indexHTML: indexHTML,
	}
}

func (s *staticServer) Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", s.Index)
	r.Mount("/static", http.FileServer(http.Dir(path.Join(s.exDir, "dist"))))

	return r
}

func (s *staticServer) Index(w http.ResponseWriter, r *http.Request) {
	w.Write(s.indexHTML)
}

type apiServer struct {
	st storage
}

func newAPIServer(s storage) *apiServer {
	return &apiServer{st: s}
}

func (s *apiServer) Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/analysis/project/*", s.BurndownProject)
	r.Get("/analysis/people/*", s.BurndownPeople)
	r.Get("/analysis/files/*", s.BurndownFiles)

	return r
}

func (s *apiServer) BurndownProject(w http.ResponseWriter, r *http.Request) {
	data, err := s.st.BurndownProject(s.uri(r))
	if err != nil {
		s.handleError(w, err)
		return
	}

	renderJSON(w, data)
}

func (s *apiServer) BurndownPeople(w http.ResponseWriter, r *http.Request) {
	data, err := s.st.BurndownPeople(s.uri(r))
	if err != nil {
		s.handleError(w, err)
		return
	}

	renderJSON(w, data)
}

func (s *apiServer) BurndownFiles(w http.ResponseWriter, r *http.Request) {
	data, err := s.st.BurndownFiles(s.uri(r))
	if err != nil {
		s.handleError(w, err)
		return
	}

	renderJSON(w, data)
}

func (s *apiServer) uri(r *http.Request) string {
	repo := chi.URLParam(r, "*")
	return "https://" + repo
}

func (s *apiServer) handleError(w http.ResponseWriter, err error) {
	var res *errResponse
	if verr, ok := err.(*validationError); ok {
		res = newErrResponse(verr, http.StatusBadRequest)
	} else if err == git.ErrRepositoryNotExists {
		res = newErrResponse(err, http.StatusBadRequest)
	} else {
		res = newErrResponse(err, http.StatusInternalServerError)
	}

	renderJSON(w, res)
}

type errResponse struct {
	Status int    `json:"-"`
	Error  string `json:"error"`
}

func newErrResponse(err error, status int) *errResponse {
	return &errResponse{
		Status: status,
		Error:  err.Error(),
	}
}

func renderJSON(w http.ResponseWriter, d interface{}) {
	status := http.StatusOK

	if resp, ok := d.(*errResponse); ok {
		status = resp.Status
	}

	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	if err := enc.Encode(d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(buf.Bytes())
}

func corsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := w.Header()
		headers.Set("Access-Control-Allow-Origin", "*")
		headers.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}
