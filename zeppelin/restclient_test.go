package zeppelin

import (
	"fmt"
	"testing"
	"time"
)

var client *Client

func init() {
	config := ClientConfig{
		ZeppelinRestUrl: "http://127.0.0.1:8080",
		Timeout:         time.Millisecond * 2000,
		RetryCount:      2,
		QueryInterval:   2000,
	}
	client = NewZeppelinClient(config)
}

func Test_CreateNote(t *testing.T) {
	noteId, err := client.createNote("/a/flink")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(noteId)
}

func Test_DeleteNote(t *testing.T) {
	var noteId = "2GTBE85ZZ"
	err := client.deleteNote(noteId)
	if err != nil {
		t.Error(err)
	}
}

func Test_AddFlinkConfigParagraph(t *testing.T) {
	var noteId = "2GRK6JUF5"
	var title = "Flink Config"
	var text = `%flink.conf 
FLINK_HOME /Users/apple/develop/bigdata/flink-1.12.5
flink.execution.mode remote
flink.execution.remote.host	127.0.0.2
flink.execution.remote.port 8082`
	paragraph, err := client.addParagraph(noteId, title, text)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(paragraph)
}

func Test_AddFlinkSqlParagraph(t *testing.T) {
	var noteId = "2GRK6JUF5"
	var title = "Flink Sql"
	var text = `%flink.ssql
create table if not exists datagen(id int,name string) with ('connector' = 'datagen','rows-per-second' = '2');
create table if not exists print(id int,name string) with ('connector' = 'print');
insert into print select * from datagen;`
	paragraph, err := client.addParagraph(noteId, title, text)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(paragraph)
}

func Test_WaitUtilFinished(t *testing.T) {
	var noteId = "2GRK6JUF5"
	var paragraphId = "paragraph_1640772348317_1851116219"
	paragraph, err := client.submitParagraph(noteId, paragraphId)
	if err != nil {
		t.Error(err)
	}
	for _, v := range paragraph.Results {
		fmt.Println(v)
	}
	paragraphFinish, err := client.waitUtilParagraphFinish(noteId, paragraphId)
	if err != nil {
		t.Error(err)
	}
	for _, v := range paragraphFinish.Results {
		fmt.Println(v)
	}
}

func Test_WaitUtilRunning(t *testing.T) {
	var noteId = "2GRK6JUF5"
	var paragraphId = "paragraph_1640772348317_1851116219"
	paragraph, err := client.submitParagraph(noteId, paragraphId)
	if err != nil {
		t.Error(err)
	}
	for _, v := range paragraph.Results {
		fmt.Println(v)
	}
	paragraphRunning, err := client.waitUtilParagraphRunning(noteId, paragraphId)
	if err != nil {
		t.Error(err)
	}
	for _, url := range paragraphRunning.JobUrls {
		fmt.Println(url)
	}
	for _, v := range paragraphRunning.Results {
		fmt.Println(v)
	}
}

func Test_NextParagraph(t *testing.T) {
	var noteId = "2GRNBSUQD"
	paragraph, err := client.nextSessionParagraph(noteId, 100)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(paragraph)
}
