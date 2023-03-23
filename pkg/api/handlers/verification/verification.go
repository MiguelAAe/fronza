package verification

import (
	"fmt"
	"net/http"
	"portal-backend/pkg/addresssvc"
)

// Postcode
// @Summary      Checks if a postcode is valid
// @Description  Checks if a postcode is valid returns 200 if successful
// @Accept       json
// @Produce      json
// @Param        Authorization  header  string  true  "token: test"
// @Param        postcode       query   string  true  "CR0 2HS"
// @Success      200
// @Failure      404
// @Router       /verify/postcode [post]
func Postcode(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("postcode")
	fmt.Println("poscode is:", code)
	ok := addresssvc.CheckPostCode(code)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
