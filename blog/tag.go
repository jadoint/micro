package blog

import (
	"database/sql"
	"errors"
	"regexp"
	"strings"

	"github.com/jadoint/micro/now"

	// Standard anonymous sql driver import
	_ "github.com/go-sql-driver/mysql"

	"github.com/jadoint/micro/conn"
	"github.com/jadoint/micro/logger"
	"github.com/jadoint/micro/validate"
)

// Tag contains data from tag table
type Tag struct {
	ID        int64  `db:"id_user,omitempty" json:"id,omitempty"`
	Tag       string `db:"tag,omitempty" json:"tag,omitempty"`
	Frequency string `db:"frequency,omitempty" json:"frequency,omitempty"`
}

// TagCSV contains data from blog_tags table
type TagCSV struct {
	Tags sql.NullString `db:"tags,omitempty" jason:"tags,omitempty"`
}

// Validate a tag
func (t *Tag) Validate() error {
	err := validate.Struct(t)
	if err != nil {
		return err
	}

	re := regexp.MustCompile("^[a-zA-Z0-9-]{2,25}$")
	if !re.MatchString(t.Tag) {
		return errors.New("Tag must only contain dashes or alphanumeric characters and be between 2 and 25 characters long")
	}
	return nil
}

// GetFrequentTags gets most frequently used tags
func GetFrequentTags(clients *conn.Clients) ([]*string, error) {
	db := clients.DB.Read
	rows, err := db.Query(`
		SELECT tag
		FROM tag
		ORDER BY frequency DESC
		LIMIT 5
		# GetFrequentTags`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*string
	for rows.Next() {
		var t string
		err := rows.Scan(&t)
		if err != nil {
			if err != sql.ErrNoRows {
				logger.Panic(err.Error())
			}
			return nil, err
		}
		tags = append(tags, &t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tags, nil
}

// GetTagsCSV Get CSV of a blog's tags
func GetTagsCSV(clients *conn.Clients, idBlog int64) (string, error) {
	db := clients.DB.Read
	var tagCsv string
	err := db.QueryRow(`
		SELECT tags
		FROM blog_tags
		WHERE id_blog = ?
		LIMIT 1
		# GetTagsCSV`, idBlog).
		Scan(&tagCsv)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Panic(err.Error())
		}
		return tagCsv, err
	}

	return tagCsv, nil
}

// SetTagsCSV Sets a blog's tag string in CSV.
// Only works as part of a transaction since all
// non-read tag operations are interdependent.
func SetTagsCSV(tx *sql.Tx, idBlog int64, tagCsv string) error {
	_, err := tx.Exec(`
		INSERT INTO blog_tags(id_blog, tags, modified)
		VALUES(?, ?, ?)
		ON DUPLICATE KEY UPDATE
		tags = ?, modified = ?`,
		idBlog, tagCsv, now.MySQLUTC(),
		tagCsv, now.MySQLUTC())
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// DeleteTagsCSV Deletes a record from blog_tags.
// Only works as part of a transaction since all
// non-read tag operations are interdependent.
func DeleteTagsCSV(tx *sql.Tx, idBlog int64) error {
	_, err := tx.Exec(`
		DELETE FROM blog_tags
		WHERE id_blog = ?
		LIMIT 1`,
		idBlog)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// AddTag inserts a tag into the tag table
func AddTag(clients *conn.Clients, idBlog int64, t *Tag) (int64, error) {
	t.Tag = strings.ToLower(t.Tag)

	// Get/Set tag csv for blog_tags
	tagCsv, _ := GetTagsCSV(clients, idBlog)
	var tags []string
	if tagCsv != "" {
		tags = strings.Split(tagCsv, ",")
	}
	numTags := len(tags)
	if numTags >= 20 {
		return 0, errors.New("Limit of 20 tags reached")
	}
	for _, v := range tags {
		if v == t.Tag {
			return 0, errors.New("Tag already exists")
		}
	}

	tx, err := clients.DB.Master.Begin()
	if err != nil {
		return 0, err
	}

	tags = append(tags, t.Tag)
	newTagCsv := strings.Join(tags, ",")
	err = SetTagsCSV(tx, idBlog, newTagCsv)
	if err != nil {
		return 0, err
	}

	// Insert tag into tag table to track usage
	// frequency and/or create a new tag ID.
	res, err := tx.Exec(`
		INSERT INTO tag(tag)
		VALUES(?)
		ON DUPLICATE KEY UPDATE
		frequency = frequency + 1`,
		t.Tag)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	idTag, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Insert tag ID into blog_tag to easily get
	// a list of blog IDs associated with a tag.
	_, err = tx.Exec(`
		INSERT IGNORE INTO blog_tag(id_tag, id_blog)
		VALUES(?, ?)`,
		idTag, idBlog)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return idTag, nil
}

// DeleteTag deletes a tag from the tag table
func DeleteTag(clients *conn.Clients, idBlog int64, tag string) error {
	// Get/Delete tag from csv in blog_tags
	tagCsv, _ := GetTagsCSV(clients, idBlog)
	var tags []string
	if tagCsv != "" {
		tags = strings.Split(tagCsv, ",")
	}
	numTags := len(tags)
	if numTags == 0 {
		return errors.New("Tag did not exist")
	}
	for i, v := range tags {
		if v == tag {
			tags[i] = tags[len(tags)-1]
			tags = tags[:len(tags)-1]
		}
	}
	newTagCsv := strings.Join(tags, ",")

	tx, err := clients.DB.Master.Begin()
	if err != nil {
		return err
	}

	if newTagCsv == "" {
		err = DeleteTagsCSV(tx, idBlog)
	} else {
		err = SetTagsCSV(tx, idBlog, newTagCsv)
	}
	if err != nil {
		return err
	}

	idTag, err := GetIDByTag(clients, tag)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM blog_tag
		WHERE id_tag = ? AND id_blog = ?
		LIMIT 1
		# DeleteTag`,
		idTag, idBlog)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Decrement frequency of tag
	_, err = tx.Exec(`
		UPDATE tag
		SET frequency = frequency - 1
		WHERE id_tag = ?
		LIMIT 1`,
		idTag)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// GetIDByTag gets a tag ID by its tag
func GetIDByTag(clients *conn.Clients, tag string) (int64, error) {
	db := clients.DB.Read
	var t Tag
	err := db.QueryRow(`
		SELECT id_tag
		FROM tag
		WHERE tag = ?
		LIMIT 1
		# GetTagByID`, tag).
		Scan(&t.ID)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Panic(err.Error())
		}
		return 0, err
	}
	return t.ID, nil
}

// GetTagByID gets a tag by its tag ID
func GetTagByID(clients *conn.Clients, idTag int64) (string, error) {
	db := clients.DB.Read
	var t Tag
	err := db.QueryRow(`
		SELECT tag
		FROM tag
		WHERE id_tag = ?
		LIMIT 1
		# GetTagByID`, idTag).
		Scan(&t.Tag)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Panic(err.Error())
		}
		return "", err
	}
	return t.Tag, nil
}
