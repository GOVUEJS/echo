package model

type ApiResponse struct {
	ErrCode *int    `json:"errCode,omitempty"`
	Message *string `json:"message,omitempty"`
}

type GetArticleList struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Date  string `json:"date"`
}

type GetArticleListResponse struct {
	ArticleList []GetArticleList `json:"articleList"`
	TotalPage   int64            `json:"totalPage"`
	Current     int64            `json:"current"`
}

type GetArticle struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Date    string `json:"date"`
}

type GetArticleResponse struct {
	Article GetArticle `json:"article"`
}
