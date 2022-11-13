package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Ladence/go-url-shortener/internal/bll"
	"github.com/Ladence/go-url-shortener/internal/model"
	"github.com/Ladence/go-url-shortener/internal/storage"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const API_QUOTA = 10 // calls quota for clients

type Server struct {
	shortener  *bll.Shortener
	resolver   *bll.Resolver
	ipStorage  storage.KvStorage
	urlStorage storage.KvStorage
	httpServer http.Server
}

func NewServer(address string, urlStorage, ipStorage storage.KvStorage) *Server {
	server := &Server{
		shortener:  bll.NewShortener(urlStorage),
		resolver:   bll.NewResolver(urlStorage),
		ipStorage:  ipStorage,
		urlStorage: urlStorage,
	}

	server.initHttpServer(address)
	return server
}

func (s *Server) Run() {
	_ = s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() error {
	err := s.urlStorage.Close()
	if err != nil {
		return fmt.Errorf("error on closing urlStorage: %v", err)
	}
	if err = s.ipStorage.Close(); err != nil {
		return fmt.Errorf("error on closing ipStorage: %v", err)
	}
	return nil
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
	req := &model.PostShortenRequest{}
	err = json.Unmarshal(bytes, req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	expiryStorage := time.Hour * 24
	if req.Expiry != nil {
		expiryStorage = *req.Expiry
	}

	// rate limiting
	clientIp := getIpFromRequest(r)
	clientRate, err := s.ipStorage.Get(context.Background(), clientIp)
	if clientRate == nil {
		_ = s.ipStorage.Push(context.Background(), clientIp, API_QUOTA, time.Minute*30)
	} else if err == nil {
		valInt, _ := strconv.Atoi(clientRate.(string))
		if valInt <= 0 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
	}

	shortenUrl, err := s.shortener.ShortenUrl(req.Url, req.CustomShort, expiryStorage)
	if err != nil {
		w.WriteHeader(http.StatusAlreadyReported)
		return
	}

	remainingQuota, _ := s.ipStorage.Decr(context.Background(), clientIp)
	response := model.PostShortenResponse{
		Url:             req.Url,
		CustomShort:     shortenUrl,
		Expiry:          expiryStorage,
		XRateLimitReset: 30, // todo: care about quota, need to receive a TTL
		XRateRemaining:  remainingQuota.(int64),
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

func getIpFromRequest(r *http.Request) (ipAddr string) {
	ipAddr = r.Header.Get("X-Real-Ip")
	if len(ipAddr) != 0 {
		return
	}
	ipAddr = r.Header.Get("X-Forwarded-For")
	if len(ipAddr) != 0 {
		return
	}
	ipAddr = r.RemoteAddr
	return
}
