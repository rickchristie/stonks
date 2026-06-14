package note

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"stonks/lib/tr"
	psApi "stonks/pservice/app/api"
	svcNoteLib "stonks/service/note/lib"
)

type Handler struct {
	service svcNoteLib.AppClient
}

func NewHandler(service svcNoteLib.AppClient) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	r.POST("/api/hello", h.hello)
	r.POST("/api/note/create", h.createNote)
	r.POST("/api/note/list", h.listNotes)
	r.POST("/api/note/archive", h.archiveNote)
}

func (h *Handler) hello(c *gin.Context) {
	out := h.service.ListNotes(c.Request.Context(), svcNoteLib.ListNotesIn{
		Trace: requestTrace(c),
	})
	if !out.Success {
		c.JSON(http.StatusOK, psApi.HelloResp{Error: apiErr(out.Error)})
		return
	}

	databaseMessage := "Database is reachable, but no active note is seeded."
	if len(out.Notes) > 0 {
		databaseMessage = out.Notes[0].Body
	}

	c.JSON(http.StatusOK, psApi.HelloResp{
		Error:           psApi.ApiErrNone,
		Message:         "Hello, World!",
		DatabaseMessage: databaseMessage,
		NoteCount:       len(out.Notes),
	})
}

func (h *Handler) createNote(c *gin.Context) {
	var req psApi.CreateNoteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, psApi.CreateNoteResp{Error: psApi.ApiErrValidation})
		return
	}

	out := h.service.CreateNote(c.Request.Context(), svcNoteLib.CreateNoteIn{
		Trace: requestTrace(c),
		Title: req.Title,
		Body:  req.Body,
	})

	c.JSON(http.StatusOK, psApi.CreateNoteResp{
		Error: apiErr(out.Error),
		Note:  psApi.NoteFromData(out.Note),
	})
}

func (h *Handler) listNotes(c *gin.Context) {
	out := h.service.ListNotes(c.Request.Context(), svcNoteLib.ListNotesIn{
		Trace: requestTrace(c),
	})

	c.JSON(http.StatusOK, psApi.ListNotesResp{
		Error: apiErr(out.Error),
		Notes: psApi.NotesFromData(out.Notes),
	})
}

func (h *Handler) archiveNote(c *gin.Context) {
	var req psApi.ArchiveNoteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, psApi.ArchiveNoteResp{Error: psApi.ApiErrValidation})
		return
	}

	out := h.service.ArchiveNote(c.Request.Context(), svcNoteLib.ArchiveNoteIn{
		Trace:  requestTrace(c),
		NoteId: req.NoteId,
	})

	c.JSON(http.StatusOK, psApi.ArchiveNoteResp{
		Error: apiErr(out.Error),
		Note:  psApi.NoteFromData(out.Note),
	})
}

func apiErr(err svcNoteLib.AppErr) psApi.ApiErr {
	switch err {
	case svcNoteLib.AppErrNone:
		return psApi.ApiErrNone
	case svcNoteLib.AppErrValidation:
		return psApi.ApiErrValidation
	case svcNoteLib.AppErrNotFound:
		return psApi.ApiErrNotFound
	default:
		return psApi.ApiErrInternalError
	}
}

func requestTrace(c *gin.Context) *tr.Trace {
	trace := tr.New(c.GetHeader("X-Request-ID"))
	trace.Path = c.FullPath()
	return trace
}
