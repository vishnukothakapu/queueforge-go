package service

import (
	"jobQueue-go/internal/model"
	"jobQueue-go/pkg/db"
	"log"
)

func CreateJob(job model.Job) error {
	query := `INSERT INTO jobs (id,type,status,data) VALUES ($1,$2,$3,$4)`
	_, err := db.DB.Exec(query, job.ID, job.Type, job.Status, job.Data)
	if err != nil {
		log.Println("DB INSERT ERROR:", err)
	}
	return err
}

func UpdateJobStatus(id string, status string) error {
	query := `UPDATE jobs SET status=$1 WHERE id=$2`
	_, err := db.DB.Exec(query, status, id)
	return err
}

func GetJobByID(id string) (model.Job, error) {
	var job model.Job

	query := `SELECT id, type, status, data, retries, max_retries FROM jobs WHERE id=$1`

	err := db.DB.QueryRow(query, id).Scan(
		&job.ID,
		&job.Type,
		&job.Status,
		&job.Retries,
		&job.MaxRetries,
	)
	return job, err
}
