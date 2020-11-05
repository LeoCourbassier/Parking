package api

import (
	"net/http"

	"br.com.mlabs/utils"
	"github.com/sirupsen/logrus"
)

// MethodNotAllowedHandler denies wrong access
func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content Type", "application/json")
	logrus.Warnf("%s %s %d %s", r.Method, r.URL, http.StatusMethodNotAllowed, utils.ErrMethodNotAllowed.Error())

	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte(utils.ErrMethodNotAllowed.Error()))
}

// NotFoundHandler denies wrong access
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content Type", "application/json")
	logrus.Warnf("%s %s %d %s", r.Method, r.URL, http.StatusNotFound, utils.ErrNotFound.Error())

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(utils.ErrNotFound.Error()))
}
