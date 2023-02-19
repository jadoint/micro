package blog

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	// Standard anonymous sql driver import
	_ "github.com/go-sql-driver/mysql"

	"github.com/jadoint/micro/pkg/clean"
	"github.com/jadoint/micro/pkg/conn"
	"github.com/jadoint/micro/pkg/logger"
	"github.com/jadoint/micro/pkg/paginate"
)

// Blog contains data from the blog tables
type Blog struct {
	ID               int64    `db:"id_blog,omitempty" json:"idPost,omitempty" validate:"min=0"`
	IDAuthor         int64    `db:"id_author,omitempty" json:"idAuthor,omitempty" validate:"required,min=0"`
	Title            string   `db:"title,omitempty" json:"title,omitempty" validate:"required,max=255"`
	Post             string   `db:"post,omitempty" json:"post,omitempty" validate:"required,max=16777215"`
	WordCount        int      `db:"word_count,omitempty" json:"wordCount,omitempty" validate:"required"`
	Created          string   `db:"created,omitempty" json:"created,omitempty"`
	Modified         string   `db:"modified,omitempty" json:"modified,omitempty"`
	ModifiedDatetime string   `db:"modified,omitempty" json:"modifiedDatetime,omitempty"`
	IsUnlisted       bool     `db:"is_unlisted,omitempty" json:"isUnlisted,omitempty"`
	IsDraft          bool     `db:"is_draft,omitempty" json:"isDraft,omitempty"`
	Tags             []string `db:"tags,omitempty" json:"tags,omitempty"`
}

// GetLatest gets latest blogs posted
func GetLatest(clients *conn.Clients, pageNum int, pageSize int) ([]*Blog, error) {
	offset := paginate.GetOffset(pageNum, pageSize)
	db := clients.DB.Read
	rows, err := db.Query(`
		SELECT b.id_blog, b.id_author, b.title, LEFT(b.post, 300) AS post,
			b.word_count, b.created, b.modified, bt.tags
		FROM blog AS b
		INNER JOIN blog_settings AS bs ON b.id_blog = bs.id_blog
		LEFT JOIN blog_tags AS bt ON b.id_blog = bt.id_blog
		WHERE bs.is_unlisted = 0
		ORDER BY b.id_blog DESC
		LIMIT ?, ?
		# GetLatest`, offset, pageSize)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Log(err)
		}
		return nil, err
	}
	defer rows.Close()

	// Strip Post of all tags
	strict := clean.Strict()

	var blogs []*Blog
	for rows.Next() {
		var b Blog
		var t TagCSV
		err := rows.Scan(&b.ID, &b.IDAuthor, &b.Title, &b.Post,
			&b.WordCount, &b.Created, &b.Modified, &t.Tags)
		if err != nil {
			if err != sql.ErrNoRows {
				logger.Log(err)
			}
			return nil, err
		}
		if t.Tags.Valid {
			tagCsv := t.Tags.String
			b.Tags = strings.Split(tagCsv, ",")
		}
		b.Post = strict.Sanitize(b.Post)
		ti, _ := time.Parse("2006-01-02 15:04:05", b.Modified)
		b.Modified = ti.Format("January 02, 2006")
		blogs = append(blogs, &b)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return blogs, nil
}

// GetLatestByTag gets latest blogs posted by tag
func GetLatestByTag(clients *conn.Clients, tag string, pageNum int, pageSize int) ([]*Blog, error) {
	offset := paginate.GetOffset(pageNum, pageSize)
	db := clients.DB.Read
	rows, err := db.Query(`
		SELECT b.id_blog, b.id_author, b.title, LEFT(b.post, 300) AS post,
			b.word_count, b.created, b.modified, bts.tags
		FROM blog AS b
		INNER JOIN blog_settings AS bs ON b.id_blog = bs.id_blog
		INNER JOIN blog_tags AS bts ON b.id_blog = bts.id_blog
		INNER JOIN blog_tag AS bt ON b.id_blog = bt.id_blog
		INNER JOIN tag AS t ON bt.id_tag = t.id_tag
		WHERE t.tag = ? AND bs.is_unlisted = 0
		ORDER BY b.id_blog DESC
		LIMIT ?, ?
		# GetLatest`, tag, offset, pageSize)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Log(err)
		}
		return nil, err
	}
	defer rows.Close()

	// Strip Post of all tags
	strict := clean.Strict()

	var blogs []*Blog
	for rows.Next() {
		var b Blog
		var t TagCSV
		err := rows.Scan(&b.ID, &b.IDAuthor, &b.Title, &b.Post,
			&b.WordCount, &b.Created, &b.Modified, &t.Tags)
		if err != nil {
			if err != sql.ErrNoRows {
				logger.Log(err)
			}
			return nil, err
		}
		if t.Tags.Valid {
			tagCsv := t.Tags.String
			b.Tags = strings.Split(tagCsv, ",")
		}
		b.Post = strict.Sanitize(b.Post)
		ti, _ := time.Parse("2006-01-02 15:04:05", b.Modified)
		b.Modified = ti.Format("January 02, 2006")
		blogs = append(blogs, &b)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return blogs, nil
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
		if err != sql.ErrNoRows {
			logger.Log(err)
		}
		return nil, err
	}
	defer rows.Close()

	var blogs []*Blog
	for rows.Next() {
		var b Blog
		err := rows.Scan(&b.ID, &b.Title)
		if err != nil {
			if err != sql.ErrNoRows {
				logger.Log(err)
			}
			return nil, err
		}
		blogs = append(blogs, &b)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return blogs, nil
}

// Get gets single blog post and its settings
func Get(clients *conn.Clients, idBlog int64) (*Blog, error) {
	db := clients.DB.Read
	var b Blog
	err := db.QueryRow(`
		SELECT b.id_blog, b.id_author, b.title, b.post,
			b.word_count, b.created, b.modified,
			bs.is_unlisted, bs.is_draft
		FROM blog AS b
		INNER JOIN blog_settings AS bs ON b.id_blog = bs.id_blog
		WHERE b.id_blog = ?
		LIMIT 1
		# Get`, idBlog).
		Scan(&b.ID, &b.IDAuthor, &b.Title, &b.Post, &b.WordCount,
			&b.Created, &b.Modified, &b.IsUnlisted, &b.IsDraft)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Log(err)
		}
		return &b, err
	}
	return &b, nil
}

// GetPostInit gets blog settings and credentials used
// for retrieving details from GetPost.
func GetPostInit(clients *conn.Clients, idBlog int64) (*Blog, error) {
	db := clients.DB.Read
	var b Blog
	err := db.QueryRow(`
		SELECT b.id_blog, b.id_author, b.modified,
			bs.is_unlisted, bs.is_draft
		FROM blog AS b
		INNER JOIN blog_settings AS bs ON b.id_blog = bs.id_blog
		WHERE b.id_blog = ?
		LIMIT 1
		# GetPostInit`, idBlog).
		Scan(&b.ID, &b.IDAuthor, &b.Modified,
			&b.IsUnlisted, &b.IsDraft)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Log(err)
		}
		return &b, err
	}
	return &b, nil
}

// GetPost gets single blog post
func GetPost(clients *conn.Clients, idBlog int64) (*Blog, error) {
	db := clients.DB.Read
	var b Blog
	err := db.QueryRow(`
		SELECT id_blog, title, post, word_count, created, modified
		FROM blog
		WHERE id_blog = ?
		LIMIT 1
		# GetPost`, idBlog).
		Scan(&b.ID, &b.Title, &b.Post, &b.WordCount,
			&b.Created, &b.Modified)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Log(err)
		}
		return &b, err
	}
	return &b, nil
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
			logger.Log(err)
		}
		return idAuthor, err
	}
	return idAuthor, nil
}

// Add creates a new blog post
func Add(clients *conn.Clients, b *Blog) (int64, error) {
	// Blog
	res, err := clients.DB.Exec(`
		INSERT INTO blog(id_author, title, post, word_count)
		VALUES(?, ?, ?, ?)`,
		b.IDAuthor, b.Title, b.Post, b.WordCount)
	if err != nil {
		logger.Log(err)
		return 0, err
	}
	idBlog, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	b.ID = idBlog

	// Settings
	err = UpdateSettings(clients, b)
	if err != nil {
		return 0, err
	}

	return idBlog, nil
}

// Update update blog post
func Update(clients *conn.Clients, b *Blog) error {
	_, err := clients.DB.Exec(`
		UPDATE blog
		SET title = ?, post = ?, word_count = ?, modified = ?
		WHERE id_blog = ?
		LIMIT 1`,
		b.Title, b.Post, b.WordCount, b.Modified, b.ID)
	if err != nil {
		logger.Log(err)
		return err
	}

	// Settings
	err = UpdateSettings(clients, b)
	if err != nil {
		return err
	}

	return nil
}

// UpdateSettings saves blog settings
func UpdateSettings(clients *conn.Clients, b *Blog) error {
	_, err := clients.DB.Exec(`
		UPDATE blog_settings
		SET is_unlisted = ?, is_draft = ?
		WHERE id_blog = ?
		LIMIT 1`,
		b.IsUnlisted, b.IsDraft, b.ID)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Log(err)
		}
		return err
	}
	return nil
}

// Delete delete blog post
func Delete(clients *conn.Clients, idBlog int64) error {
	_, err := clients.DB.Exec(`
		DELETE FROM blog
		WHERE id_blog = ?
		LIMIT 1`,
		idBlog)
	if err != nil {
		logger.Log(err)
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
			logger.Log(err)
		}
		return 0, err
	}
	return views, nil
}

// IncrViews increments blog views
func IncrViews(clients *conn.Clients, idBlog int64) (int64, error) {
	ctx := context.Background()
	bvKey := fmt.Sprintf("blog:views:%d", idBlog)
	res := clients.Cache.Get(ctx, bvKey)
	if res.Err() != nil && res.Err().Error() == redis.Nil.Error() {
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
		clients.Cache.IncrBy(ctx, bvKey, views)

		return views, nil
	}

	views, err := res.Int64()
	if err != nil {
		log.Println(err.Error())
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
			if err != sql.ErrNoRows {
				logger.Log(err)
			}
			return 0, err
		}
	}

	clients.Cache.Incr(ctx, bvKey)

	return views + 1, nil
}
