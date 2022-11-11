package server

import (
	"encoding/json"
	"github.com/Ladence/go-url-shortener/internal/bll"
	"github.com/Ladence/go-url-shortener/internal/model"
	"github.com/Ladence/go-url-shortener/internal/storage"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"time"
)

type Server struct {
	shortener  *bll.Shortener
	resolver   *bll.Resolver
	ipStorage  storage.KvStorage
	httpServer http.Server
}

func NewServer(address string, urlStorage, ipStorage storage.KvStorage) *Server {
	server := &Server{
		shortener: bll.NewShortener(urlStorage),
		resolver:  bll.NewResolver(urlStorage),
		ipStorage: ipStorage,
	}

	server.initHttpServer(address)
	return server
}

func (s *Server) Run() {
	_ = s.httpServer.ListenAndServe()
}

func (s *Server) initHttpServer(address string) {
	r := mux.NewRouter()
	r.HandleFunc("/shorten", s.handleShorten).Methods(http.MethodPost)
	r.HandleFunc("/resolve", s.handleResolve).Methods(http.MethodGet)

	s.httpServer = http.Server{
		Addr:    address,
		Handler: r,
	}
}

// handleShorten handles POST /shorten requests
func (s *Server) handleShorten(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req := &model.GetShortenRequest{}
	err = json.Unmarshal(bytes, req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expiryStorage := time.Hour * 24
	if req.Expiry != nil {
		expiryStorage = *req.Expiry
	}

	shortenUrl, err := s.shortener.ShortenUrl(req.Url, req.CustomShort, expiryStorage)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := model.GetShortenResponse{
		Url:             req.Url,
		CustomShort:     shortenUrl,
		Expiry:          expiryStorage,
		XRateLimitReset: 30, // todo: care about quota
		XRateRemaining:  20,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// handleResolve handles GET /resolve?url=... requests
func (s *Server) handleResolve(w http.ResponseWriter, r *http.Request) {
	shortUrl := r.URL.Query().Get("url")
	url, err := s.resolver.Resolve(shortUrl)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, url, http.StatusSeeOther)
}
