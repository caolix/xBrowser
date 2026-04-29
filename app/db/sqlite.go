package db

import (
	"database/sql"
	"errors"
	"os"
	"time"

	_ "github.com/glebarez/go-sqlite"
	"oBrowser/app/models"
)

type AppSqlite struct {
	DB     *sql.DB
	Logger interface{}
}

func (s *AppSqlite) Init(dir string) (err error) {
	dbAddress := dir + string(os.PathSeparator) + DB_NAME
	s.DB, err = sql.Open("sqlite", dbAddress)
	if err != nil {
		return err
	}

	_, err = s.DB.Exec(`
		CREATE TABLE IF NOT EXISTS login_info (
			endpoint   TEXT NOT NULL,
			access_key TEXT NOT NULL,
			secret_key TEXT NOT NULL,
			account_id TEXT NOT NULL DEFAULT '',
			remark     TEXT NOT NULL DEFAULT '',
			prepath    TEXT NOT NULL DEFAULT '',
			login_time TEXT NOT NULL DEFAULT '',
			use_ssl    INTEGER NOT NULL DEFAULT 0,
			PRIMARY KEY (endpoint, access_key, secret_key)
		);
		CREATE TABLE IF NOT EXISTS upload_task (
			account_id     TEXT    NOT NULL,
			task_id        TEXT    NOT NULL,
			bucket         TEXT    NOT NULL DEFAULT '',
			key            TEXT    NOT NULL DEFAULT '',
			name           TEXT    NOT NULL DEFAULT '',
			source         TEXT    NOT NULL DEFAULT '',
			size           INTEGER NOT NULL DEFAULT 0,
			human_size     TEXT    NOT NULL DEFAULT '',
			uploaded_size  INTEGER NOT NULL DEFAULT 0,
			upload_id      TEXT    NOT NULL DEFAULT '',
			is_multipart   INTEGER NOT NULL DEFAULT 0,
			part_size      INTEGER NOT NULL DEFAULT 0,
			status         INTEGER NOT NULL DEFAULT 0,
			modified_time  TEXT    NOT NULL DEFAULT '',
			PRIMARY KEY (account_id, task_id)
		);
		CREATE TABLE IF NOT EXISTS download_task (
			account_id     TEXT    NOT NULL,
			task_id        TEXT    NOT NULL,
			bucket         TEXT    NOT NULL DEFAULT '',
			key            TEXT    NOT NULL DEFAULT '',
			version_id     TEXT    NOT NULL DEFAULT '',
			name           TEXT    NOT NULL DEFAULT '',
			destination    TEXT    NOT NULL DEFAULT '',
			size           INTEGER NOT NULL DEFAULT 0,
			human_size     TEXT    NOT NULL DEFAULT '',
			status         INTEGER NOT NULL DEFAULT 0,
			PRIMARY KEY (account_id, task_id)
		);
		CREATE TABLE IF NOT EXISTS completed_download_part (
			task_id   TEXT    NOT NULL,
			offset    INTEGER NOT NULL,
			part_size INTEGER NOT NULL DEFAULT 0,
			PRIMARY KEY (task_id, offset)
		);
		CREATE TABLE IF NOT EXISTS settings (
			account_id               TEXT    NOT NULL PRIMARY KEY,
			part_size_mb             INTEGER NOT NULL DEFAULT 5,
			upload_parts_concurrency INTEGER NOT NULL DEFAULT 2,
			upload_concurrency       INTEGER NOT NULL DEFAULT 2,
			download_concurrency     INTEGER NOT NULL DEFAULT 1
		);
	`)
	if err != nil {
		s.DB = nil
		return err
	}
	return nil
}

func (s *AppSqlite) Close() {
	if s.DB != nil {
		s.DB.Close()
	}
}

func (s *AppSqlite) UpsertLoginInfo(l *models.LoginInfo) (err error) {
	var count int
	err = s.DB.QueryRow(`
		SELECT COUNT(*) FROM login_info
		WHERE endpoint = ? AND access_key = ? AND secret_key = ?
	`, l.Endpoint, l.AccessKey, l.SecretKey).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		_, err = s.DB.Exec(`
			INSERT INTO login_info (endpoint, access_key, secret_key, account_id, remark, prepath, login_time, use_ssl)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, l.Endpoint, l.AccessKey, l.SecretKey, l.AccountId, l.Remark, l.Prepath,
			l.LoginTime.Format(time.RFC3339), boolToInt(l.UseSSL))
		return err
	}
	_, err = s.DB.Exec(`
		UPDATE login_info SET login_time = ?, remark = ?, prepath = ?, use_ssl = ?
		WHERE endpoint = ? AND access_key = ? AND secret_key = ?
	`, time.Now().Local().Format(time.RFC3339), l.Remark, l.Prepath, boolToInt(l.UseSSL),
		l.Endpoint, l.AccessKey, l.SecretKey)
	return err
}

func (s *AppSqlite) GetLatestLoginInfo() (*models.LoginInfo, error) {
	info := &models.LoginInfo{}
	var loginTimeStr string
	var useSSL int
	err := s.DB.QueryRow(`
		SELECT endpoint, access_key, secret_key, account_id, remark, prepath, login_time, use_ssl
		FROM login_info ORDER BY login_time DESC LIMIT 1
	`).Scan(&info.Endpoint, &info.AccessKey, &info.SecretKey, &info.AccountId,
		&info.Remark, &info.Prepath, &loginTimeStr, &useSSL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	info.LoginTime, _ = time.Parse(time.RFC3339, loginTimeStr)
	info.UseSSL = intToBool(useSSL)
	return info, nil
}

func (s *AppSqlite) ListAllLoginInfo() ([]models.LoginInfo, error) {
	rows, err := s.DB.Query(`
		SELECT endpoint, access_key, secret_key, account_id, remark, prepath, login_time, use_ssl
		FROM login_info
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.LoginInfo
	for rows.Next() {
		var info models.LoginInfo
		var loginTimeStr string
		var useSSL int
		if err := rows.Scan(&info.Endpoint, &info.AccessKey, &info.SecretKey, &info.AccountId,
			&info.Remark, &info.Prepath, &loginTimeStr, &useSSL); err != nil {
			return nil, err
		}
		info.LoginTime, _ = time.Parse(time.RFC3339, loginTimeStr)
		info.UseSSL = intToBool(useSSL)
		results = append(results, info)
	}
	return results, rows.Err()
}

func (s *AppSqlite) ListAllUploadTasks(accountId string) ([]models.UploadTask, error) {
	rows, err := s.DB.Query(`
		SELECT account_id, task_id, bucket, key, name, source, size, human_size,
		       uploaded_size, upload_id, is_multipart, part_size, status, modified_time
		FROM upload_task WHERE account_id = ?
	`, accountId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.UploadTask
	for rows.Next() {
		var t models.UploadTask
		var modifiedTimeStr string
		if err := rows.Scan(&t.AccountId, &t.TaskId, &t.Bucket, &t.Key, &t.Name, &t.Source,
			&t.Size, &t.HumanSize, &t.UploadedSize, &t.UploadId, &t.IsMultipart, &t.PartSize,
			&t.Status, &modifiedTimeStr); err != nil {
			return nil, err
		}
		t.ModifiedTime, _ = time.Parse(time.RFC3339, modifiedTimeStr)
		results = append(results, t)
	}
	return results, rows.Err()
}

func (s *AppSqlite) CreateUploadTask(u *models.UploadTask) error {
	_, err := s.DB.Exec(`
		INSERT INTO upload_task (account_id, task_id, bucket, key, name, source, size, human_size,
		                         uploaded_size, upload_id, is_multipart, part_size, status, modified_time)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, u.AccountId, u.TaskId, u.Bucket, u.Key, u.Name, u.Source, u.Size, u.HumanSize,
		u.UploadedSize, u.UploadId, u.IsMultipart, u.PartSize, u.Status,
		u.ModifiedTime.Format(time.RFC3339))
	return err
}

func (s *AppSqlite) DeleteUploadTask(accountId string, taskId string) {
	s.DB.Exec(`DELETE FROM upload_task WHERE account_id = ? AND task_id = ?`, accountId, taskId)
}

func (s *AppSqlite) ListAllDownloadTasks(accountId string) ([]models.DownloadTask, error) {
	rows, err := s.DB.Query(`
		SELECT account_id, task_id, bucket, key, version_id, name, destination,
		       size, human_size, status
		FROM download_task WHERE account_id = ?
	`, accountId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.DownloadTask
	for rows.Next() {
		var t models.DownloadTask
		if err := rows.Scan(&t.AccountId, &t.TaskId, &t.Bucket, &t.Key, &t.VersionId,
			&t.Name, &t.Destination, &t.Size, &t.HumanSize, &t.Status); err != nil {
			return nil, err
		}
		results = append(results, t)
	}
	return results, rows.Err()
}

func (s *AppSqlite) CreateDownloadTask(u *models.DownloadTask) error {
	var existingTaskId string
	err := s.DB.QueryRow(`
		SELECT task_id FROM download_task WHERE account_id = ? AND task_id = ?
	`, u.AccountId, u.TaskId).Scan(&existingTaskId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err = s.DB.Exec(`
				INSERT INTO download_task (account_id, task_id, bucket, key, version_id, name, destination,
				                           size, human_size, status)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`, u.AccountId, u.TaskId, u.Bucket, u.Key, u.VersionId, u.Name, u.Destination,
				u.Size, u.HumanSize, u.Status)
			return err
		}
		return err
	}
	_, err = s.DB.Exec(`
		UPDATE download_task SET status = ?
		WHERE account_id = ? AND task_id = ?
	`, u.Status, u.AccountId, u.TaskId)
	return err
}

func (s *AppSqlite) CreateDownloadPart(taskId string, offset, partSize int64) error {
	_, err := s.DB.Exec(`
		INSERT INTO completed_download_part (task_id, offset, part_size)
		VALUES (?, ?, ?)
	`, taskId, offset, partSize)
	return err
}

func (s *AppSqlite) GetAllDownloadParts(taskId string) ([]*models.CompletedDownloadPart, error) {
	rows, err := s.DB.Query(`
		SELECT task_id, offset, part_size
		FROM completed_download_part WHERE task_id = ? ORDER BY offset
	`, taskId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parts []*models.CompletedDownloadPart
	var partOffsetMarker int64 = 0
	for rows.Next() {
		var p models.CompletedDownloadPart
		if err := rows.Scan(&p.TaskId, &p.Offset, &p.PartSize); err != nil {
			return nil, err
		}
		if p.Offset == partOffsetMarker {
			partOffsetMarker = p.Offset + p.PartSize
			parts = append(parts, &p)
		} else {
			break
		}
	}
	return parts, rows.Err()
}

func (s *AppSqlite) DeleteDownloadTask(accountId string, taskId string) {
	tx, err := s.DB.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()
	tx.Exec(`DELETE FROM download_task WHERE account_id = ? AND task_id = ?`, accountId, taskId)
	tx.Exec(`DELETE FROM completed_download_part WHERE task_id = ?`, taskId)
	tx.Commit()
}

func (s *AppSqlite) UpdateSettings(settings *models.Settings) error {
	_, err := s.DB.Exec(`
		INSERT INTO settings (account_id, part_size_mb, upload_parts_concurrency, upload_concurrency, download_concurrency)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(account_id) DO UPDATE SET
			part_size_mb = excluded.part_size_mb,
			upload_parts_concurrency = excluded.upload_parts_concurrency,
			upload_concurrency = excluded.upload_concurrency,
			download_concurrency = excluded.download_concurrency
	`, settings.AccountId, settings.PartSizeMB, settings.UploadPartsConcurrency,
		settings.UploadConcurrency, settings.DownloadConcurrency)
	return err
}

func (s *AppSqlite) LoadSettings(accountId string) (*models.Settings, error) {
	settings := &models.Settings{}
	err := s.DB.QueryRow(`
		SELECT account_id, part_size_mb, upload_parts_concurrency, upload_concurrency, download_concurrency
		FROM settings WHERE account_id = ?
	`, accountId).Scan(&settings.AccountId, &settings.PartSizeMB, &settings.UploadPartsConcurrency,
		&settings.UploadConcurrency, &settings.DownloadConcurrency)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			settings.SetDefault()
			settings.AccountId = accountId
			_, err = s.DB.Exec(`
				INSERT INTO settings (account_id, part_size_mb, upload_parts_concurrency, upload_concurrency, download_concurrency)
				VALUES (?, ?, ?, ?, ?)
			`, settings.AccountId, settings.PartSizeMB, settings.UploadPartsConcurrency, settings.UploadConcurrency, settings.DownloadConcurrency)
			if err != nil {
				return nil, err
			}
			return settings, nil
		}
		return nil, err
	}
	return settings, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func intToBool(i int) bool {
	return i != 0
}
