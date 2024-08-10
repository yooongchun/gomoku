package ai

import (
	"github.com/duke-git/lancet/v2/fileutil"
	"testing"
)

func TestSave(t *testing.T) {
	board := NewBoard(12, ROLE_HUMAN)
	board.Move(NewPoint(1, 1))
	board.Move(NewPoint(2, 2))
	board.Move(NewPoint(3, 3))
	board.Move(NewPoint(4, 4))
	board.Save("test.txt")
	textLines, err := fileutil.ReadFileByLine("../history/test.txt")
	if err != nil {
		t.Errorf("error reading file:%v", err)
	}
	if len(textLines) != 7 {
		t.Errorf("Expected line 7 but got:%d", len(textLines))
	}
	if textLines[0] != "#meta info:size,firstRole" {
		t.Errorf("Expected #meta info:size,firstRole but got:%s", textLines[0])
	}
	if textLines[1] != "12,1" {
		t.Errorf("Expected 12,1 but got:%s", textLines[1])
	}
	if textLines[2] != "#history info:x,y,chess" {
		t.Errorf("Expected #history info:x,y,chess but got:%s", textLines[2])
	}
	if textLines[3] != "1,1,1" {
		t.Errorf("Expected 1,1,1 but got:%s", textLines[3])
	}
}

func TestLoad(t *testing.T) {
	board := Load("../history/test.txt")
	if board == nil {
		t.Error("Expected board but got nil")
		return
	}
	if board.size != 12 {
		t.Errorf("Expected size 12 but got:%d", board.size)
	}
	if board.firstRole != ROLE_HUMAN {
		t.Errorf("Expected firstRole 1 but got:%d", board.firstRole)
	}
	if board.current != CHESS_BLACK {
		t.Errorf("Expected current 2 but got:%d", board.current)
	}
	if board.history[0].point != NewPoint(1, 1) {
		t.Errorf("Expected history 1,1 but got:%v", board.history[0])
	}
}
