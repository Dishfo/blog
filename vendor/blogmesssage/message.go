package blogmesssage

//在两个进程之间需要把文章的创建，修改与删除及时的传递给对方

const (
	CreateArticle = 1 << iota
	UpdateArticle
	DeleteAricle
)

type ArticleMessage struct {
	Op             int
	ArticleId      int64
	ArticleContent string
}
