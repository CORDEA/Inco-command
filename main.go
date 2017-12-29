package main

import (
	"bufio"
	"crypto/rsa"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/jroimartin/gocui"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
)

const (
	KeyPath    = "private.pem"
	TableName  = "table"
	DialogName = "dialog"
)

func historyTitles(key *rsa.PrivateKey, histories []History) []string {
	var titles []string
	for i := range histories {
		decoded, err := base64.StdEncoding.DecodeString(histories[i].Url)
		if err != nil {
			glog.Warning(err)
		}
		decrypted, err := Decrypt(key, decoded)
		if err != nil {
			glog.Warning(err)
		}
		glog.Info(decrypted)
		titles = append(titles, decrypted)
	}
	return titles
}

func main() {
	flag.Parse()

	request := NewRequest()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Username: ")
	user, err := reader.ReadString('\n')
	if err != nil {
		glog.Error(err)
	}
	fmt.Print("Password: ")

	bytePass, err := terminal.ReadPassword(0)
	if err != nil {
		glog.Error(err)
	}

	user = strings.Trim(user, "\n")
	pass := strings.Trim(string(bytePass), "\n")

	if err := request.Login(user, pass); err != nil {
		glog.Error(err)
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		glog.Error(err)
	}
	defer g.Close()

	g.Cursor = true

	key, err := ReadPrivateKey(KeyPath)
	if err != nil {
		glog.Error(err)
	}

	histories, err := request.GetHistories()
	if err != nil {
		glog.Error(err)
	}
	titles := historyTitles(key, histories)
	table := NewTable(TableName, titles)
	g.SetManagerFunc(table.Layout)

	handler := NewHandler(
		request,
		TableName,
		DialogName,
		titles,
		histories,
	)

	if err := handler.KeyBindings(g); err != nil {
		glog.Error(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		glog.Error(err)
	}
}
