package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/goodleby/golang-app/article"
)

type fakeArticleInserter struct {
	articles   []article.Article
	shouldFail bool
}

func (f *fakeArticleInserter) InsertArticle(ctx context.Context, article article.Article) error {
	if f.shouldFail {
		return errors.New("test error")
	}

	f.articles = append(f.articles, article)

	return nil
}

func TestAddArticle(t *testing.T) {
	type args struct {
		articleInserter *fakeArticleInserter
		req             *http.Request
	}
	tests := []struct {
		name         string
		args         args
		wantStatus   int
		wantErr      bool
		wantArticles []article.Article
	}{
		{
			name: "should add the passed article to the database",
			args: args{
				articleInserter: &fakeArticleInserter{
					articles: []article.Article{
						{
							ID:          "test_id",
							Title:       "Test title",
							Description: "Test description",
							Body:        "Test body",
						},
					},
					shouldFail: false,
				},
				req: addChiURLParams(httptest.NewRequest(http.MethodPatch, "/articles/test_id", makeJSONBody(t, &article.Article{
					ID:          "new_id",
					Title:       "New title",
					Description: "New description",
					Body:        "New body",
				})), map[string]string{
					"id": "test_id",
				}),
			},
			wantStatus: http.StatusNoContent,
			wantErr:    false,
			wantArticles: []article.Article{
				{
					ID:          "test_id",
					Title:       "Test title",
					Description: "Test description",
					Body:        "Test body",
				},
				{
					ID:          "new_id",
					Title:       "New title",
					Description: "New description",
					Body:        "New body",
				},
			},
		},
		{
			name: "should return an internal error if it fails to insert the article to the database",
			args: args{
				articleInserter: &fakeArticleInserter{
					articles: []article.Article{
						{
							ID:          "test_id",
							Title:       "Test title",
							Description: "Test description",
							Body:        "Test body",
						},
					},
					shouldFail: true,
				},
				req: addChiURLParams(httptest.NewRequest(http.MethodPatch, "/articles/test_id", makeJSONBody(t, &article.Article{
					ID:          "new_id",
					Title:       "New title",
					Description: "New description",
					Body:        "New body",
				})), map[string]string{
					"id": "test_id",
				}),
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    true,
			wantArticles: []article.Article{
				{
					ID:          "test_id",
					Title:       "Test title",
					Description: "Test description",
					Body:        "Test body",
				},
			},
		},
		{
			name: "should return a bad request error if wrong body is provided",
			args: args{
				articleInserter: &fakeArticleInserter{
					articles: []article.Article{
						{
							ID:          "test_id",
							Title:       "Test title",
							Description: "Test description",
							Body:        "Test body",
						},
					},
					shouldFail: true,
				},
				req: addChiURLParams(httptest.NewRequest(http.MethodPatch, "/articles/test_id", nil), map[string]string{
					"id": "test_id",
				}),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
			wantArticles: []article.Article{
				{
					ID:          "test_id",
					Title:       "Test title",
					Description: "Test description",
					Body:        "Test body",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			handler := AddArticle(tt.args.articleInserter)
			handler(w, tt.args.req)

			if w.Code != tt.wantStatus {
				t.Fatalf("GetArticle() status = %v, want %v", w.Code, tt.wantStatus)
			}

			if !reflect.DeepEqual(tt.args.articleInserter.articles, tt.wantArticles) {
				t.Fatalf("GetArticle() articles = %v, want %v", tt.args.articleInserter.articles, tt.wantArticles)
			}

			// If we expect an error, we just need to check the response body is not empty.
			if tt.wantErr && w.Body.Len() == 0 {
				t.Fatalf("GetArticle() response body is empty, want error")
			}
		})
	}
}
