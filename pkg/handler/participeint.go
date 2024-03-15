package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"testApi"
)

func (h *Handler) createParticipant(c *gin.Context) {
	groupId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid group param")
		return
	}
	var input testApi.Participant
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid input")
	}

	id, err := h.services.Participant.Create(groupId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) tossParticipant(c *gin.Context) {
	groupId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid group param")
		return
	}

	participants, err := h.services.Participant.Toss(groupId)
	if err != nil {
		newErrorResponse(c, http.StatusConflict, err.Error())
		return
	}

	c.JSON(http.StatusOK, participants)
}

func (h *Handler) getInfoAboutRecipient(c *gin.Context) {
	groupId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid Group Id param")
		return
	}

	participantId, err := strconv.Atoi(c.Param("participantId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid Participant Id param")
		return
	}

	recipient, err := h.services.Participant.GetRecipientInfo(groupId, participantId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Recipient not found or not exist")
		return
	}

	c.JSON(http.StatusOK, recipient)
}

func (h *Handler) deleteParticipant(c *gin.Context) {
	groupId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid group ID param")
		return
	}

	participantId, err := strconv.Atoi(c.Param("participantId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid participant ID param")
		return
	}

	err = h.services.Participant.Delete(groupId, participantId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "OK",
	})
}
