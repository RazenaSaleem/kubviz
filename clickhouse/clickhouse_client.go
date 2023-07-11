package clickhouse

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go"
	"github.com/intelops/kubviz/model"
)

func GetClickHouseConnection(url string) (*sql.DB, error) {
	connect, err := sql.Open("clickhouse", url)
	//connect, err := sql.Open("clickhouse", "tcp://kubviz-client-clickhouse:9000?debug=true")
	if err != nil {
		log.Fatal(err)
	}
	if err := connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println(err)
		}
		return nil, err
	}

	return connect, nil
}

func CreateSchema(connect *sql.DB) {
	_, err := connect.Exec(`
		CREATE TABLE IF NOT EXISTS events (
			ClusterName String,
			Id           UUID,
			EventTime   DateTime,
			OpType      String,
			Name         String,
			Namespace    String,
			Kind         String,
			Message      String,
			Reason       String,
			Host         String,
			Event        String,
			FirstTime   DateTime,
			LastTime    DateTime
		) engine=File(TabSeparated)
	`)

	if err != nil {
		log.Fatal(err)
	}
}
func CreateRakeesMetricsSchema(connect *sql.DB) {
	_, err := connect.Exec(`
		CREATE TABLE IF NOT EXISTS rakkess (
			ClusterName String,
			Name String,
			Create String,
			Delete String,
			List String,
			Update String
        ) engine=File(TabSeparated)
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateKubePugSchema(connect *sql.DB) {
	_, err := connect.Exec(`
        CREATE TABLE IF NOT EXISTS DeprecatedAPIs (
			ClusterName String,
			ObjectName String,
            Description String,
            Kind String,
            Deprecated UInt8,
            Scope String
        ) engine=File(TabSeparated)
    `)
	if err != nil {
		log.Fatal(err)
	}

	_, err = connect.Exec(`
        CREATE TABLE IF NOT EXISTS DeletedAPIs (
			ClusterName String,
			ObjectName String,
            Group String,
            Kind String,
            Version String,
            Name String,
            Deleted UInt8,
            Scope String
        ) engine=File(TabSeparated)
    `)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateKetallSchema(connect *sql.DB) {
	_, err := connect.Exec(`
		CREATE TABLE IF NOT EXISTS getall_resources (
			ClusterName String,
			Namespace String,
			Kind String,
			Resource String,
			Age String
        ) engine=File(TabSeparated)
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateOutdatedSchema(connect *sql.DB) {
	_, err := connect.Exec(`
	    CREATE TABLE IF NOT EXISTS outdated_images (
			ClusterName String,
			Namespace String,
			Pod String,
		    CurrentImage String,
			CurrentTag String,
			LatestVersion String,
			VersionsBehind Int64
	    ) engine=File(TabSeparated)
	`)
	if err != nil {
		log.Fatal(err)
	}
}
func InsertRakeesMetrics(connect *sql.DB, metrics model.RakeesMetrics) {
	var (
		tx, _   = connect.Begin()
		stmt, _ = tx.Prepare("INSERT INTO rakkess (ClusterName, Name, Create, Delete, List, Update) VALUES (?, ?, ?, ?, ?, ?)")
	)
	defer stmt.Close()
	if _, err := stmt.Exec(
		metrics.ClusterName,
		metrics.Name,
		metrics.Create,
		metrics.Delete,
		metrics.List,
		metrics.Update,
	); err != nil {
		log.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}

func CreateKubeScoreSchema(connect *sql.DB) {
	_, err := connect.Exec(`
	    CREATE TABLE IF NOT EXISTS kubescore (
		    id UUID,
			namespace String,
			cluster_name String,
			recommendations String
	    ) engine=File(TabSeparated)
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func InsertKetallEvent(connect *sql.DB, metrics model.Resource) {
	var (
		tx, _   = connect.Begin()
		stmt, _ = tx.Prepare("INSERT INTO getall_resources (ClusterName, Namespace, Kind, Resource, Age) VALUES (?, ?, ?, ?, ?)")
	)
	defer stmt.Close()
	if _, err := stmt.Exec(
		metrics.ClusterName,
		metrics.Namespace,
		metrics.Kind,
		metrics.Resource,
		metrics.Age,
	); err != nil {
		log.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}

func InsertOutdatedEvent(connect *sql.DB, metrics model.CheckResultfinal) {
	var (
		tx, _   = connect.Begin()
		stmt, _ = tx.Prepare("INSERT INTO outdated_images (ClusterName, Namespace, Pod, CurrentImage, CurrentTag, LatestVersion, VersionsBehind) VALUES (?, ?, ?, ?, ?, ?, ?)")
	)
	defer stmt.Close()
	if _, err := stmt.Exec(
		metrics.ClusterName,
		metrics.Namespace,
		metrics.Pod,
		metrics.Image,
		metrics.Current,
		metrics.LatestVersion,
		metrics.VersionsBehind,
	); err != nil {
		log.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}

func InsertDeprecatedAPI(connect *sql.DB, deprecatedAPI model.DeprecatedAPI) {
	var (
		tx, _   = connect.Begin()
		stmt, _ = tx.Prepare("INSERT INTO DeprecatedAPIs (ClusterName, ObjectName, Description, Kind, Deprecated, Scope) VALUES (?, ?, ?, ?, ?, ?)")
	)
	defer stmt.Close()

	deprecated := uint8(0)
	if deprecatedAPI.Deprecated {
		deprecated = 1
	}

	for _, item := range deprecatedAPI.Items {
		if _, err := stmt.Exec(
			deprecatedAPI.ClusterName,
			item.ObjectName,
			deprecatedAPI.Description,
			deprecatedAPI.Kind,
			deprecated,
			item.Scope,
		); err != nil {
			log.Fatal(err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}

func InsertDeletedAPI(connect *sql.DB, deletedAPI model.DeletedAPI) {
	var (
		tx, _   = connect.Begin()
		stmt, _ = tx.Prepare("INSERT INTO DeletedAPIs (ClusterName, ObjectName, Group, Kind, Version, Name, Deleted, Scope) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	)
	defer stmt.Close()

	deleted := uint8(0)
	if deletedAPI.Deleted {
		deleted = 1
	}

	for _, item := range deletedAPI.Items {
		if _, err := stmt.Exec(
			deletedAPI.ClusterName,
			item.ObjectName,
			deletedAPI.Group,
			deletedAPI.Kind,
			deletedAPI.Version,
			deletedAPI.Name,
			deleted,
			item.Scope,
		); err != nil {
			log.Fatal(err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}

func InsertEvent(connect *sql.DB, metrics model.Metrics) {
	var (
		tx, _   = connect.Begin()
		stmt, _ = tx.Prepare("INSERT INTO events (ClusterName, Id, EventTime, OpType, Name, Namespace, Kind, Message, Reason, Host, Event, FirstTime, LastTime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	)
	defer stmt.Close()
	eventJson, _ := json.Marshal(metrics.Event)
	if _, err := stmt.Exec(
		metrics.ClusterName,
		metrics.Event.UID,
		time.Now(),
		metrics.Type,
		metrics.Event.Name,
		metrics.Event.Namespace,
		metrics.Event.InvolvedObject.Kind,
		metrics.Event.Message,
		metrics.Event.Reason,
		metrics.Event.Source.Host,
		string(eventJson),
		metrics.Event.FirstTimestamp.Time,
		metrics.Event.LastTimestamp.Time,
	); err != nil {
		log.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}

func InsertKubeScoreMetrics(connect *sql.DB, metrics model.KubeScoreRecommendations) {
	var (
		tx, _   = connect.Begin()
		stmt, _ = tx.Prepare("INSERT INTO kubescore (id, namespace, cluster_name, recommendations) VALUES (?, ?, ?, ?)")
	)
	defer stmt.Close()
	if _, err := stmt.Exec(
		metrics.ID,
		metrics.Namespace,
		metrics.ClusterName,
		metrics.Recommendations,
	); err != nil {
		log.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}

func RetriveKetallEvent(connect *sql.DB) ([]model.Resource, error) {
	rows, err := connect.Query("SELECT ClusterName, Namespace, Kind, Resource, Age FROM getall_resources")
	if err != nil {
		log.Printf("Error: %s", err)
		return nil, err
	}
	defer rows.Close()
	var events []model.Resource
	for rows.Next() {
		var result model.Resource
		if err := rows.Scan(&result.ClusterName, &result.Namespace, &result.Kind, &result.Resource, &result.Age); err != nil {
			log.Printf("Error: %s", err)
			return nil, err
		}
		events = append(events, result)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error: %s", err)
		return nil, err
	}
	return events, nil
}

func RetriveOutdatedEvent(connect *sql.DB) ([]model.CheckResultfinal, error) {
	rows, err := connect.Query("SELECT ClusterName, Namespace, Pod, CurrentImage, CurrentTag, LatestVersion, VersionsBehind FROM outdated_images")
	if err != nil {
		log.Printf("Error: %s", err)
		return nil, err
	}
	defer rows.Close()
	var events []model.CheckResultfinal
	for rows.Next() {
		var result model.CheckResultfinal
		if err := rows.Scan(&result.ClusterName, &result.Namespace, &result.Pod, &result.Image, &result.Current, &result.LatestVersion, &result.VersionsBehind); err != nil {
			log.Printf("Error: %s", err)
			return nil, err
		}
		events = append(events, result)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error: %s", err)
		return nil, err
	}
	return events, nil
}

func RetriveKubepugEvent(connect *sql.DB) ([]model.Result, error) {
	rows, err := connect.Query("SELECT result, cluster_name FROM deprecatedAPIs_and_deletedAPIs")
	if err != nil {
		log.Printf("Error: %s", err)
		return nil, err
	}
	defer rows.Close()
	var events []model.Result
	for rows.Next() {
		var result model.Result
		if err := rows.Scan(&result, &result.ClusterName); err != nil {
			log.Printf("Error: %s", err)
			return nil, err
		}
		events = append(events, result)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error: %s", err)
		return nil, err
	}
	return events, nil
}

func RetrieveEvent(connect *sql.DB) ([]model.DbEvent, error) {
	rows, err := connect.Query("SELECT ClusterName, Id, EventTime, OpType, Name, Namespace, Kind, Message, Reason, Host, Event, FirstTime, LastTime FROM events")
	if err != nil {
		log.Printf("Error: %s", err)
		return nil, err
	}
	defer rows.Close()
	var events []model.DbEvent
	for rows.Next() {
		var dbEvent model.DbEvent
		if err := rows.Scan(&dbEvent.Cluster_name, &dbEvent.Id, &dbEvent.Event_time, &dbEvent.Op_type, &dbEvent.Name, &dbEvent.Namespace, &dbEvent.Kind, &dbEvent.Message, &dbEvent.Host, &dbEvent.Event, &dbEvent.First_time, &dbEvent.Last_time); err != nil {
			log.Printf("Error: %s", err)
			return nil, err
		}
		eventJson, _ := json.Marshal(dbEvent)
		log.Printf("DB Event: %s", string(eventJson))
		events = append(events, dbEvent)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error: %s", err)
		return nil, err
	}
	return events, nil
}
