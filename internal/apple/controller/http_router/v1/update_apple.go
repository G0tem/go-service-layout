package v1

import "net/http"

// @Summary UpdateApple
// @Description Update Apple
// @Tags apple
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID"
// @Success 200 {object} map[string]string "status: Apple"
// @Failure 400 {object} map[string]string "invalid payload"
// @Failure 401 {object} map[string]string "unauthorized"
// @Failure 500 {object} map[string]string "internal error"
// @Router /apples/update_apple/{id} [put]
func (h *Handlers) UpdateApple(w http.ResponseWriter, r *http.Request) {

}
