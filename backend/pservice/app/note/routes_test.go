package note

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"stonks/data"
	psApi "stonks/pservice/app/api"
	svcNote "stonks/service/note/lib"
)

type mockNoteSvc struct {
	createFn  func(context.Context, svcNote.CreateNoteIn) svcNote.CreateNoteOut
	listFn    func(context.Context, svcNote.ListNotesIn) svcNote.ListNotesOut
	archiveFn func(context.Context, svcNote.ArchiveNoteIn) svcNote.ArchiveNoteOut
}

func (m *mockNoteSvc) CreateNote(ctx context.Context, in svcNote.CreateNoteIn) svcNote.CreateNoteOut {
	if m.createFn == nil {
		return svcNote.CreateNoteOut{}
	}
	return m.createFn(ctx, in)
}

func (m *mockNoteSvc) ListNotes(ctx context.Context, in svcNote.ListNotesIn) svcNote.ListNotesOut {
	if m.listFn == nil {
		return svcNote.ListNotesOut{}
	}
	return m.listFn(ctx, in)
}

func (m *mockNoteSvc) ArchiveNote(ctx context.Context, in svcNote.ArchiveNoteIn) svcNote.ArchiveNoteOut {
	if m.archiveFn == nil {
		return svcNote.ArchiveNoteOut{}
	}
	return m.archiveFn(ctx, in)
}

func newTestRouter(service svcNote.AppClient) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	NewHandler(service).RegisterRoutes(r)
	return r
}

func postJSON(t *testing.T, router http.Handler, path string, body any) *httptest.ResponseRecorder {
	t.Helper()

	buf := &bytes.Buffer{}
	if body != nil {
		require.NoError(t, json.NewEncoder(buf).Encode(body))
	}
	req := httptest.NewRequest(http.MethodPost, path, buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func testNote(id int, title string, body string, status data.NoteStatus) *data.Note {
	now := time.Date(2026, 6, 14, 1, 2, 3, 0, time.UTC)
	return &data.Note{
		Id:            id,
		Title:         title,
		Body:          body,
		Status:        status,
		CreatedTs:     now,
		LastUpdatedTs: now,
	}
}

func TestHello(t *testing.T) {
	t.Run("returns database proof from first active note", func(t *testing.T) {
		router := newTestRouter(&mockNoteSvc{
			listFn: func(_ context.Context, in svcNote.ListNotesIn) svcNote.ListNotesOut {
				require.NotNil(t, in.Trace)
				assert.Equal(t, "/api/hello", in.Trace.Path)
				return svcNote.ListNotesOut{
					Success: true,
					Error:   svcNote.AppErrNone,
					Notes: []*data.Note{
						testNote(1, "Hello", "from db", data.NSActive),
						testNote(2, "Second", "ignored", data.NSActive),
					},
				}
			},
		})

		w := postJSON(t, router, "/api/hello", psApi.HelloReq{})

		assert.Equal(t, http.StatusOK, w.Code)
		var resp psApi.HelloResp
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, psApi.ApiErrNone, resp.Error)
		assert.Equal(t, "Hello, World!", resp.Message)
		assert.Equal(t, "from db", resp.DatabaseMessage)
		assert.Equal(t, 2, resp.NoteCount)
	})

	t.Run("maps service failure to API error", func(t *testing.T) {
		router := newTestRouter(&mockNoteSvc{
			listFn: func(context.Context, svcNote.ListNotesIn) svcNote.ListNotesOut {
				return svcNote.ListNotesOut{Error: svcNote.AppErrInternalError}
			},
		})

		w := postJSON(t, router, "/api/hello", psApi.HelloReq{})

		assert.Equal(t, http.StatusOK, w.Code)
		var resp psApi.HelloResp
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, psApi.ApiErrInternalError, resp.Error)
	})
}

func TestCreateNoteRoute(t *testing.T) {
	t.Run("binds request and returns created note", func(t *testing.T) {
		router := newTestRouter(&mockNoteSvc{
			createFn: func(_ context.Context, in svcNote.CreateNoteIn) svcNote.CreateNoteOut {
				require.NotNil(t, in.Trace)
				assert.Equal(t, "Title", in.Title)
				assert.Equal(t, "Body", in.Body)
				return svcNote.CreateNoteOut{
					Success: true,
					Error:   svcNote.AppErrNone,
					Note:    testNote(7, "Title", "Body", data.NSActive),
				}
			},
		})

		w := postJSON(t, router, "/api/note/create", psApi.CreateNoteReq{Title: "Title", Body: "Body"})

		assert.Equal(t, http.StatusOK, w.Code)
		var resp psApi.CreateNoteResp
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, psApi.ApiErrNone, resp.Error)
		require.NotNil(t, resp.Note)
		assert.Equal(t, 7, resp.Note.Id)
		assert.Equal(t, data.NSActive, resp.Note.Status)
	})

	t.Run("returns validation error for invalid JSON", func(t *testing.T) {
		router := newTestRouter(&mockNoteSvc{})
		req := httptest.NewRequest(http.MethodPost, "/api/note/create", bytes.NewBufferString("{"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp psApi.CreateNoteResp
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, psApi.ApiErrValidation, resp.Error)
		assert.Nil(t, resp.Note)
	})

	t.Run("maps service validation error", func(t *testing.T) {
		router := newTestRouter(&mockNoteSvc{
			createFn: func(context.Context, svcNote.CreateNoteIn) svcNote.CreateNoteOut {
				return svcNote.CreateNoteOut{Error: svcNote.AppErrValidation}
			},
		})

		w := postJSON(t, router, "/api/note/create", psApi.CreateNoteReq{Title: " ", Body: "Body"})

		assert.Equal(t, http.StatusOK, w.Code)
		var resp psApi.CreateNoteResp
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, psApi.ApiErrValidation, resp.Error)
		assert.Nil(t, resp.Note)
	})
}

func TestListNotesRoute(t *testing.T) {
	router := newTestRouter(&mockNoteSvc{
		listFn: func(context.Context, svcNote.ListNotesIn) svcNote.ListNotesOut {
			return svcNote.ListNotesOut{
				Success: true,
				Error:   svcNote.AppErrNone,
				Notes: []*data.Note{
					testNote(1, "A", "Body A", data.NSActive),
					testNote(2, "B", "Body B", data.NSActive),
				},
			}
		},
	})

	w := postJSON(t, router, "/api/note/list", psApi.ListNotesReq{})

	assert.Equal(t, http.StatusOK, w.Code)
	var resp psApi.ListNotesResp
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, psApi.ApiErrNone, resp.Error)
	require.Len(t, resp.Notes, 2)
	assert.Equal(t, "A", resp.Notes[0].Title)
	assert.Equal(t, "B", resp.Notes[1].Title)
}

func TestArchiveNoteRoute(t *testing.T) {
	t.Run("returns archived note", func(t *testing.T) {
		router := newTestRouter(&mockNoteSvc{
			archiveFn: func(_ context.Context, in svcNote.ArchiveNoteIn) svcNote.ArchiveNoteOut {
				assert.Equal(t, 8, in.NoteId)
				return svcNote.ArchiveNoteOut{
					Success: true,
					Error:   svcNote.AppErrNone,
					Note:    testNote(8, "Old", "Body", data.NSArchived),
				}
			},
		})

		w := postJSON(t, router, "/api/note/archive", psApi.ArchiveNoteReq{NoteId: 8})

		assert.Equal(t, http.StatusOK, w.Code)
		var resp psApi.ArchiveNoteResp
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, psApi.ApiErrNone, resp.Error)
		require.NotNil(t, resp.Note)
		assert.Equal(t, data.NSArchived, resp.Note.Status)
	})

	t.Run("maps not found", func(t *testing.T) {
		router := newTestRouter(&mockNoteSvc{
			archiveFn: func(context.Context, svcNote.ArchiveNoteIn) svcNote.ArchiveNoteOut {
				return svcNote.ArchiveNoteOut{Error: svcNote.AppErrNotFound}
			},
		})

		w := postJSON(t, router, "/api/note/archive", psApi.ArchiveNoteReq{NoteId: 99})

		assert.Equal(t, http.StatusOK, w.Code)
		var resp psApi.ArchiveNoteResp
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, psApi.ApiErrNotFound, resp.Error)
		assert.Nil(t, resp.Note)
	})
}
