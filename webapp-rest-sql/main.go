package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world!")
}

func handlerShowFile(w http.ResponseWriter, r *http.Request) {

	var resp string

	defer func() {
		fmt.Fprintf(w, "%s", resp)
	}()

	ex, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// change url to folder structure
	mp := strings.Replace(r.URL.Path, "/show/", "/markdown/", 1)

	fullpath := path.Join(ex, mp)
	log.Println("using path ", fullpath)

	fileinfo, err := os.Stat(fullpath)
	if err != nil {
		log.Printf("stat error %s path %s", err, fullpath)
		http.NotFound(w, r)
		return
	}

	file, err := os.Open(fullpath)
	if err != nil {
		log.Printf("open error %s path %s", err, fullpath)
		return
	}

	data := make([]byte, fileinfo.Size())
	count, err := io.ReadFull(file, data)
	if int64(count) != fileinfo.Size() {
		log.Printf("read error %s path %s", err, fullpath)
		return
	}

	resp = resp + fmt.Sprintf("%s", data)
}

func handlerMarkdown(w http.ResponseWriter, r *http.Request) {

	var resp string

	defer func() {
		fmt.Fprintf(w, "%s", resp)
	}()

	ex, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fullpath := path.Join(ex, r.URL.Path)
	log.Println("using path ", fullpath)

	fileinfo, err := os.Stat(fullpath)
	if err != nil || fileinfo.IsDir() {
		log.Printf("stat error %s path %s", err, fullpath)
		http.NotFound(w, r)
		return
	}

	data, err := exec.Command("markdown", fullpath).Output()
	if err != nil {
		log.Printf("exec error %s path %s", err, fullpath)
		return
	}

	resp = resp + fmt.Sprintf("%s", data)
}

type TLDentry struct {
	Id    string
	Name  string
	Tld   string
	Count int
}

type Status struct {
	Stat     bool
	RowCount int64
}

//
// instead of using a router or additional libs
// this demo does all the boilerplate code
//

func handlerRestTld(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("sqlite3", "./sample.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("url", r.Method, r.URL.Path)

	custid := strings.Trim(strings.Replace(r.URL.Path, "/resttlds/", "", 1), " ")
	fmt.Println("custid", custid)

	where := ""
	if len(custid) > 0 {
		where = fmt.Sprintf(" where id ='%v'", custid)
	}

	switch r.Method {
	case http.MethodGet:
		rows, err := db.Query("select * from tlds" + where)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		resp := make([]TLDentry, 0)

		for rows.Next() {
			var id string
			var name string
			var tld string
			var count int
			err = rows.Scan(&id, &name, &tld, &count)
			if err != nil {
				log.Fatal(err)
			}

			entry := TLDentry{id, name, tld, count}
			resp = append(resp, entry)
			//fmt.Println(entry)
		}

		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(resp)

		b, err := json.Marshal(resp)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", b)

	case http.MethodPost:

		entry := TLDentry{}
		err = json.NewDecoder(r.Body).Decode(&entry)
		if err != nil {
			log.Println(err)
			http.Error(w, "bad request", http.StatusBadRequest)
			break
		}

		if len(custid) > 0 && len(entry.Id) > 0 && strings.Compare(custid, entry.Id) == 0 {

			// update

			stmt, err := db.Prepare("update tlds set name=?, tld=?, count=? where id=?")
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()

			st := Status{false, 0}

			dbres, err := stmt.Exec(entry.Name, entry.Tld, entry.Count, entry.Id)

			if err == nil {
				st.Stat = true
				st.RowCount, _ = dbres.RowsAffected()
			}

			// create the response
			resp, err := json.Marshal(st)
			if err != nil {
				log.Fatal(err)
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, "%s", resp)

			return
		}

		if len(custid) == 0 && len(entry.Id) == 0 {

			// insert new

			uid := uuid.New().String()

			tx, err := db.Begin()
			if err != nil {
				log.Fatal(err)
			}
			stmt, err := tx.Prepare("insert into tlds(id, name, tld, count) values(?, ?, ?, ?)")
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()

			_, err = stmt.Exec(uid, entry.Name, entry.Tld, entry.Count)
			if err != nil {
				log.Fatal(err)
			}

			tx.Commit()

			st := Status{true, 1}

			// create the response
			resp, err := json.Marshal(st)
			if err != nil {
				log.Fatal(err)
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, "%s", resp)

			return
		}

		http.Error(w, "id mismatch", http.StatusBadRequest)

	case http.MethodDelete:

		if len(custid) == 0 {
			http.Error(w, "id missing", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("delete from tlds" + where)
		if err != nil {
			log.Fatal(err)
		}

		// create the response
		resp, err := json.Marshal(Status{true, 0})
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", resp)

	default:
		http.Error(w, http.ErrNotSupported.ErrorString, http.StatusNotImplemented)

	}

}

func main() {
	http.Handle("/stat/", http.StripPrefix("/stat/", http.FileServer(http.Dir("static/"))))
	http.HandleFunc("/", handler)
	http.HandleFunc("/show/", handlerShowFile)
	http.HandleFunc("/markdown/", handlerMarkdown)
	http.HandleFunc("/resttlds/", handlerRestTld)
	log.Fatal(http.ListenAndServe(":8082", nil))
}
