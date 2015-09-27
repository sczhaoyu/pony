package admin_server

import (
	"net/http"
)

func HtppRun() {
	http.HandleFunc("/logic", getLogicAddr)
	http.HandleFunc("/logic/list", getLogicLiost)
	http.ListenAndServe(":3869", nil)
}
