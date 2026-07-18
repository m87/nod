package main

import (
	"log/slog"

	"github.com/m87/nod"
	sqlite_nod "github.com/m87/nod/sqlite"
)

// Article implements nod.NodeCodec, so it does not need a registered adapter.
type Article struct {
	Title     string
	Slug      string
	Body      string
	Published bool
	Labels    []string
}

func (article *Article) ToNode() (*nod.Node, error) {
	status := "draft"
	if article.Published {
		status = "published"
	}

	tags := make([]*nod.Tag, 0, len(article.Labels))
	for _, label := range article.Labels {
		tags = append(tags, &nod.Tag{Name: label})
	}

	return &nod.Node{
		Core: nod.NodeCore{
			Name:   article.Title,
			Kind:   "article",
			Status: status,
		},
		KV: map[string]*nod.NodeKV{
			"slug": {Key: "slug", ValueText: nod.Ptr(article.Slug)},
		},
		Content: map[string]*nod.NodeContent{
			"body": {Key: "body", Value: nod.Ptr(article.Body)},
		},
		Tags: tags,
	}, nil
}

func (article *Article) FromNode(node *nod.Node) error {
	article.Title = node.Core.Name
	article.Published = node.Core.Status == "published"

	if slug := node.KV["slug"]; slug != nil && slug.ValueText != nil {
		article.Slug = *slug.ValueText
	}
	if body := node.Content["body"]; body != nil && body.Value != nil {
		article.Body = *body.Value
	}

	article.Labels = make([]string, 0, len(node.Tags))
	for _, tag := range node.Tags {
		article.Labels = append(article.Labels, tag.Name)
	}

	return nil
}

func (article *Article) IsApplicable(node *nod.Node) bool {
	return node != nil && node.Core.Kind == "article"
}

func main() {
	repo, err := sqlite_nod.NewRepositoryInMemory(slog.Default(), nil)
	if err != nil {
		panic(err)
	}
	defer repo.Close()

	article := &Article{
		Title:     "Typed queries with codecs",
		Slug:      "typed-codec-query",
		Body:      "A model can encode and decode itself.",
		Published: true,
		Labels:    []string{"go", "nod"},
	}
	articles := nod.Nodes[Article](repo)
	id, err := articles.SaveNode(article)
	if err != nil {
		panic(err)
	}

	found, err := articles.Query().
		WithKV().
		WithContent().
		WithTags().
		Where(nod.NodeFields.Id.Equals(id)).
		FindAll()
	if err != nil {
		panic(err)
	}
	if len(found) == 0 {
		panic("article not found")
	}

	slog.Info("Found article", "title", found[0].Title, "slug", found[0].Slug)
}
