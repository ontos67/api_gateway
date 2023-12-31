package api2

import (
	"github.com/gorilla/mux"
)

type API struct {
	r    *mux.Router
	Conf Config
}

type pPage interface{}

type Paging struct {
	PageTotal   int
	ItemPerPage int
	PageN       int
	Page        pPage
}

type Config struct {
	Port        string `json:"port"`
	Commentator string `json:"commentator"`
	Agrigator   string `json:"agrigator"`
}

type Article struct {
	ID        int    // номер записи
	Title     string // заголовок публикации
	Content   string // содержание публикации
	PubTime   int64  // время публикации
	Url       string // ссылка на источник
	Publisher string // название источника
	Autor     string // Имя автора
	Comments  []Comment
}

type Comment struct {
	ID         int    // номер комментария
	User_id    int    // ID автора комментария
	Text       string // содержание комментария
	PubTime    int64  // время публикации, Unixtime
	ParentType string // тип родителя (A - статья (комментарий на саму статью), С - комментарий (отчеточка на комментарий) )
	ParentID   int    // ID родителя (или статьи или комментария)
}

func New() *API {
	a := API{r: mux.NewRouter()}
	a.endpoints()
	return &a
}

func (api *API) Router() *mux.Router {
	return api.r
}

//ок Comments.saveCom http://localhost:999/comment/save?userid=64&text=заманали%20комары&pubtime=12344134&ptype=A&pid=2345
//ок Comments.deleteCom http://localhost:999/comment/del?id=64
// Comments.comListP (Parent) http://localhost:999/comment/comListP?pT=C&pId=47
// Comments.comListP по паганацией (Parent) http://localhost:999/comment/comListP?pT=C&pId=47&page=2

//ok news.lastArticles http://localhost:998/news/last?n=100&itemperpage=5&page=2
//ok news.lastArticlesList http://localhost:998/news/lastlist?n=5
//ok news.newsFilteredDetailed http://localhost:998/news/filter?time1=1699016144&time2=1700293140&lim=0&field=title&contains=putin&sortfield=id&dir=s
//ok news.newsFullDetailed http://localhost:998/news/news?id=5

func (api *API) endpoints() {
	api.r.HandleFunc("/comment/save", api.commentSaveHandler)
	api.r.HandleFunc("/comment/del", api.commentDelHandler)
	api.r.HandleFunc("/comment/comListP", api.commenListPHandler)
	api.r.HandleFunc("/comment/comListPPage", api.commenListPHandler)

	api.r.HandleFunc("/news/last", api.lastHandler)
	api.r.HandleFunc("/news/lastlist", api.lastListHandler)
	api.r.HandleFunc("/news/filter", api.filterHandler)
	api.r.HandleFunc("/news/news", api.newsHandler)
	api.r.HandleFunc("/news/newss", api.newsHandlerSynh)
}
