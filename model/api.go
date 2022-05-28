package model

type ApiResponse struct {
	Message *string `json:"message,omitempty"`
}

type PostLoginRequest struct {
	Email string `json:"email"`
	Pw    string `json:"pw"`
}

type PostLoginResponse struct {
	ApiResponse
	AccessToken  *string `json:"accessToken"`
	RefreshToken *string `json:"refreshToken"`
}

type GetArticleListResponse struct {
	ApiResponse
	ArticleList []GetArticleList `json:"articleList"`
	TotalPage   int64            `json:"totalPage"`
	Current     int64            `json:"current"`
}

type GetArticleList struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Date   string `json:"date"`
	Writer string `json:"writer"`
}

type GetArticleResponse struct {
	ApiResponse
	Article GetArticle `json:"article"`
}

type GetArticle struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Date    string `json:"date"`
	Writer  string `json:"writer"`
}
