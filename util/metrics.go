package util

import (
	"database/sql"
	_ "github.com/lib/pq"
	"hitachienergy.com/cr-operator/generated"
	"log"
	"os"
)

var db *sql.DB

func init() {
	dbConnString := os.Getenv("DATABASE")
	if dbConnString == "" {
		return
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("postgres", dbConnString)
	if err != nil {
		db = nil
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		db = nil
		log.Fatal(pingErr)
	}
	log.Println("Connected to DB!")
	initSchema()
}

func initSchema() {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS ckp_metrics(
		id serial primary key,
		start_time int,
		kubelet_duration int,
		checkpoint_image_creation_start_time int,
		checkpoint_image_creation_duration int,
		transfer_start_time int,
		transfer_duration int,
		freezing_time int,
		frozen_time int,
		memdump_time int,
		memwrite_time int,
		pages_scanned int,
		pages_skipped_parent int,
		pages_written int,
		irmap_resolve int,
		pages_lazy int,
		page_pipes int,
		page_pipe_bufs int,
		shpages_scanned int,
		shpages_skipped_parent int,
		shpages_written int,
		pod_name varchar,
		application varchar,
		host varchar
	);`)
	if err != nil {
		return
	}
}

func SaveMetrics(
	statsEntry *generated.StatsEntry,
	startTime int64,
	kubeletDuration int,
	checkpointImageCreationStartTime int64,
	checkpointImageCreationDuration int,
	transferStartTime int64,
	transferDuration int,
	metricName string,
	application string,
	host string,
) {
	if db == nil {
		return
	}
	_, err := db.Exec(`INSERT INTO ckp_metrics(
                        start_time,
                        kubelet_duration,
                        checkpoint_image_creation_start_time,
                        checkpoint_image_creation_duration,
                        transfer_start_time,
                        transfer_duration,
                        freezing_time,
					    frozen_time,
						memdump_time,
						memwrite_time,
						pages_scanned,
						pages_skipped_parent,
						pages_written,
						irmap_resolve,
						pages_lazy,
					  	page_pipes,
					  	page_pipe_bufs,
						shpages_scanned,
					  	shpages_skipped_parent,
					  	shpages_written,
                        pod_name,
                        application,
                        host)
                   VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23);`,
		startTime,
		kubeletDuration,
		checkpointImageCreationStartTime,
		checkpointImageCreationDuration,
		transferStartTime,
		transferDuration,
		statsEntry.Dump.FreezingTime,
		statsEntry.Dump.FrozenTime,
		statsEntry.Dump.MemdumpTime,
		statsEntry.Dump.MemwriteTime,
		statsEntry.Dump.PagesScanned,
		statsEntry.Dump.PagesSkippedParent,
		statsEntry.Dump.PagesWritten,
		statsEntry.Dump.IrmapResolve,
		statsEntry.Dump.PagesLazy,
		statsEntry.Dump.PagePipes,
		statsEntry.Dump.PagePipeBufs,
		statsEntry.Dump.ShpagesScanned,
		statsEntry.Dump.ShpagesSkippedParent,
		statsEntry.Dump.ShpagesWritten,
		metricName,
		application,
		host,
	)
	if err != nil {
		return
	}
}
