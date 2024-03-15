package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"testApi"
)

func (h *Handler) createGroup(c *gin.Context) {
	var input testApi.Group
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid input")
		return
	}

	id, err := h.services.Group.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Error of create")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getAllGroups(c *gin.Context) {
	groups, err := h.services.Group.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Error of get all groups")
		return
	}

	c.JSON(http.StatusOK, getAllGroupsResponse{
		Data: groups,
	})

}

func (h *Handler) getGroupById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid id param")
		return
	}

	group, err := h.services.Group.GetById(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Error of get group by ID")
		return
	}

	c.JSON(http.StatusOK, group)
}

func (h *Handler) updateGroup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid ID param")
		return
	}

	var input testApi.UpdateGroupInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid input")
		return
	}

	if err := h.services.Group.Update(id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Error of update group")
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "OK",
	})

}

func (h *Handler) deleteGroup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid id param")
		return
	}

	err = h.services.Group.Delete(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Error of delete group")
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "OK",
	})
}
