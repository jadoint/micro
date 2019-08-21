package model

import (
	"database/sql"
	"fmt"

	"github.com/go-redis/redis"
	// Standard anonymous sql driver import
	_ "github.com/go-sql-driver/mysql"

	"github.com/jadoint/micro/conn"
	"github.com/jadoint/micro/logger"
)

// Blog contains data from the blog tables
type Blog struct {
	ID         int64  `db:"id_blog,omitempty" json:"idPost,omitempty" validate:"min=0"`
	IDAuthor   int64  `db:"id_author,omitempty" json:"idAuthor,omitempty" validate:"required,min=0"`
	Title      string `db:"title,omitempty" json:"title,omitempty" validate:"required,max=255"`
	Post       string `db:"post,omitempty" json:"post,omitempty" validate:"required,max=16777215"`
	WordCount  int    `db:"word_count,omitempty" json:"wordCount,omitempty" validate:"required"`
	Created    string `db:"created,omitempty" json:"created,omitempty"`
	Modified   string `db:"modified,omitempty" json:"modified,omitempty"`
	IsUnlisted bool   `db:"is_unlisted,omitempty" json:"isUnlisted,omitempty"`
	IsDraft    bool   `db:"is_draft,omitempty" json:"isDraft,omitempty"`
}

// GetRecentAuthorBlogs gets all recent blogs by an author
func GetRecentAuthorBlogs(clients *conn.Clients, idAuthor int64) ([]*Blog, error) {
	db := clients.DB.Read
	rows, err := db.Query(`
		SELECT b.id_blog, b.title
		FROM blog AS b
		INNER JOIN blog_settings AS bs ON b.id_blog = bs.id_blog
		WHERE
			b.id_author = ? AND
			bs.is_unlisted = 0
		ORDER BY b.id_blog DESC
		LIMIT 5
		# GetRecentAuthorBlogs`, idAuthor)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []*Blog
	for rows.Next() {
		b := &Blog{}
		err := rows.Scan(&b.ID, &b.Title)
		if err != nil {
			if err != sql.ErrNoRows {
				logger.Panic(err.Error())
			}
			return nil, err
		}
		blogs = append(blogs, b)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return blogs, nil
}

// GetBlog gets single blog post and its settings
func GetBlog(clients *conn.Clients, idBlog int64) (*Blog, error) {
	db := clients.DB.Read
	b := &Blog{}
	err := db.QueryRow(`
		SELECT b.id_blog, b.id_author, b.title, b.post,
			b.word_count, b.created, b.modified,
			bs.is_unlisted, bs.is_draft
		FROM blog AS b
		INNER JOIN blog_settings AS bs ON b.id_blog = bs.id_blog
		WHERE b.id_blog = ?
		LIMIT 1
		# GetBlog`, idBlog).
		Scan(&b.ID, &b.IDAuthor, &b.Title, &b.Post, &b.WordCount,
			&b.Created, &b.Modified, &b.IsUnlisted, &b.IsDraft)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Panic(err.Error())
		}
		return b, err
	}

	return b, nil
}

// GetBlogPost gets single blog post
func GetBlogPost(clients *conn.Clients, idBlog int64) (*Blog, error) {
	db := clients.DB.Read
	b := &Blog{}
	err := db.QueryRow(`
		SELECT id_blog, title, post, word_count, created, modified
		FROM blog
		WHERE id_blog = ?
		LIMIT 1
		# GetBlogPost`, idBlog).
		Scan(&b.ID, &b.Title, &b.Post, &b.WordCount,
			&b.Created, &b.Modified)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Panic(err.Error())
		}
		return b, err
	}
	return b, nil
}

// GetIDAuthor gets author ID of blog
func GetIDAuthor(clients *conn.Clients, idBlog int64) (int64, error) {
	db := clients.DB.Read
	var idAuthor int64
	err := db.QueryRow(`
		SELECT id_author
		FROM blog
		WHERE id_blog = ?
		LIMIT 1
		# GetIdAuthor`, idBlog).
		Scan(&idAuthor)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Panic(err.Error())
		}
		return idAuthor, err
	}
	return idAuthor, nil
}

// AddBlog add blog post
func AddBlog(clients *conn.Clients, b *Blog) (int64, error) {
	// Blog
	res, err := clients.DB.Exec(`
		INSERT INTO blog(id_author, title, post, word_count)
		VALUES(?, ?, ?, ?)`,
		b.IDAuthor, b.Title, b.Post, b.WordCount)
	if err != nil {
		return 0, err
	}
	idBlog, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	b.ID = idBlog

	// Settings
	go UpdateBlogSettings(clients, b)

	return idBlog, nil
}

// UpdateBlog update blog post
func UpdateBlog(clients *conn.Clients, b *Blog) error {
	_, err := clients.DB.Exec(`
		UPDATE blog
		SET title = ?, post = ?, word_count = ?, modified = ?
		WHERE id_blog = ?
		LIMIT 1`,
		b.Title, b.Post, b.WordCount, b.Modified, b.ID)
	if err != nil {
		return err
	}

	// Settings
	go UpdateBlogSettings(clients, b)

	return nil
}

// UpdateBlogSettings saves blog settings
func UpdateBlogSettings(clients *conn.Clients, b *Blog) {
	_, err := clients.DB.Exec(`
		UPDATE blog_settings
		SET is_unlisted = ?, is_draft = ?
		WHERE id_blog = ?
		LIMIT 1`,
		b.IsUnlisted, b.IsDraft, b.ID)
	if err != nil {
		logger.Panic("UpdateBlogSettings() failed to save")
	}
}

// DeleteBlog delete blog post
func DeleteBlog(clients *conn.Clients, idBlog int64) error {
	_, err := clients.DB.Exec(`
		DELETE FROM blog
		WHERE id_blog = ?
		LIMIT 1`,
		idBlog)
	if err != nil {
		return err
	}
	return nil
}

// GetViews gets blog view count
func GetViews(clients *conn.Clients, idBlog int64) (int64, error) {
	var views int64
	dbRead := clients.DB.Read
	err := dbRead.QueryRow(`
			SELECT views
			FROM blog_views
			WHERE id_blog = ?
			LIMIT 1
			# IncrViews`, idBlog).
		Scan(&views)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Panic(err.Error())
		}
		return 0, err
	}
	return views, nil
}

// IncrViews increments blog views
func IncrViews(clients *conn.Clients, idBlog int64) (int64, error) {
	bvKey := fmt.Sprintf("blog:views:%d", idBlog)
	res := clients.Cache.Get(bvKey)
	if res.Err() == redis.Nil {
		// Not found in cache so check the database
		views, err := GetViews(clients, idBlog)
		if err != nil {
			return 0, err
		}

		// Reset cache to latest view count
		// found in the database.
		if views == 0 {
			views = 1
		}
		clients.Cache.IncrBy(bvKey, views)

		return views, nil
	}

	views, err := res.Int64()
	if err != nil {
		return 0, err
	}

	// Reduce database writes by
	// limiting view count updates
	if views%10 == 0 {
		_, err := clients.DB.Exec(`
			UPDATE blog_views
			SET views = views + 10
			WHERE id_blog = ?
			LIMIT 1`, idBlog)
		if err != nil {
			logger.Panic("IncrViews() failed to save")
		}
	}

	clients.Cache.Incr(bvKey)

	return views + 1, nil
}
