package v1

import "github.com/gin-gonic/gin"

// @Summary DeleteApple
// @Description Delete Apple
// @Tags apple
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string "status: Apple delete"
// @Failure 400 {object} map[string]string "invalid payload"
// @Failure 401 {object} map[string]string "unauthorized"
// @Failure 500 {object} map[string]string "internal error"
// @Router /apples/delete_apple/{id} [delete]
func (h *Handlers) DeleteApple(c *gin.Context) {

}
