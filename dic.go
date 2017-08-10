package main

import (
	"io/ioutil"
	"log"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"errors"
	"net/http"
	"os"
	"io"
	"strconv"
	"reflect"
	"github.com/jeek120/dic/player"
	"fmt"
	"os/exec"
	"path/filepath"
)

var CMD = []string{"play", "text", "sh"}
var OP= []string{"attr", "sub", "text", "index"}

type _config struct {
	Url 		string		`json:"url"`
	Filter 		[]string	`json:"filter"`
	Dir 		[]string 		`json:"dir"`
	Enable		bool		`json:"enable"`
	Cmd 		string		`json:"cmd"`

	dir 		string
	scheme		string
}

var configs map[string]*_config = make(map[string]*_config)
func getCurrentDir() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(filepath.Dir(file))
	if err != nil {
		return "", err
	}
	return path, nil
}

func init(){

	// json.Unmarshal is nil if configs is map[string]*_config
	var cs map[string]_config = make(map[string]_config)
	bytes,err := ioutil.ReadFile("config.json")
	if err != nil {
		var dir string
		dir, err = getCurrentDir()
		if err != nil {
			log.Fatal("get current dir:",err)
		}
		bytes,err = ioutil.ReadFile(dir+ "/config.json")
		//bytes,err = ioutil.ReadFile("/Users/jeekyuan/go/src/github.com/jeek120/dic/config.json")
	}
	if err != nil {
		log.Fatal("config.json not found:",err)
	}

	err = json.Unmarshal(bytes, &cs)
	if err != nil {
		log.Fatal("config file:" , err)
		os.Exit(-1)
	}

	for k,c := range cs{
		_c := _config{}
		_c.scheme = k
		_c.Cmd = c.Cmd
		_c.Url = c.Url
		_c.Dir = c.Dir
		_c.Enable = c.Enable
		_c.Filter = c.Filter
		for _,d := range c.Dir {
			dir := d
			idx := strings.Index(d, "/{")
			if idx > 0 {
				dir = d[0:idx]
			}
			_, err := os.Stat(dir)
			if err == nil {
				_c.dir = d
				break
			}
		}
		if c.Cmd == "" {
			_c.Cmd = CMD[1]
		}
		if _c.dir == "" {
			_c.Enable = false
		}
		configs[k] = &_c
	}
}

func get(obj interface{},op, param string) (dest interface{}, err error) {
	dest = nil
	err = nil
	if o, ok := obj.(*goquery.Document); ok {
		switch op {
		case OP[0]:
			//attr
			dest,_ = o.Attr(param)
		case OP[1]:
			dest = o.Find(param)
		case OP[2]:
			dest = o.Text()
		case OP[3]:
			var idx int
			if idx,err = strconv.Atoi(param); err ==nil {
				dest = o.Eq(idx)
			}else{
				err = errors.New(op + ":" + param + " error:" + err.Error())
				return
			}
		default:
			err = errors.New(op + ":" + param + "cannot support op for " + reflect.TypeOf(obj).String())
			return
		}
	}else if o,ok := obj.(*goquery.Selection); ok {
		switch op {
		case OP[0]:
			//attr
			dest,_ = o.Attr(param)
		case OP[1]:
			dest = o.Find(param)
		case OP[2]:
			dest = o.Text()
		case OP[3]:
			var idx int
			if idx,err = strconv.Atoi(param); err == nil {
				dest = o.Eq(idx)
			}else{
				err = errors.New(op + ":" + param + " error" + err.Error())
				return
			}
		default:
			err = errors.New(op + ":" + param + "cannot support op for " + reflect.TypeOf(obj).String())
			return
		}
	}else{
		err = errors.New(op + ":" + param + "cannot support type for " + reflect.TypeOf(obj).String())
		return
	}
	return
}

func (c *_config)exec(word string) error{
	var url = strings.Replace(c.Url, "{word}", word, -1)
	var doc interface{}
	var err error
	doc, err = goquery.NewDocument(url)
	if err != nil {
		return err
	}
	for _,f := range c.Filter {
		v := strings.Split(f, ":")
		var op,param string
		if len(v) == 2 {
			op = v[0]
			param = v[1]
		} else if len(v) == 1 {
			op = "sub"
			param = v[0]
		} else {
			return errors.New("error filter for " + f)
		}

		if doc, err = get(doc, op, param); err != nil {
			return err
		}
	}
	var result string

	if o,ok1 := doc.(*goquery.Selection); ok1 {
		result = o.Text()
	} else if result,ok1 = doc.(string); !ok1{
		return errors.New("The result is not text")
	}

	//cache
	var path = strings.Replace(c.dir, "{first}", word[0:0], -1)
	path = strings.Replace(path, "{word}", word, -1)
	path = strings.Replace(path, "{scheme}", c.scheme, -1)
	if strings.Contains(path, "{suff}") {
		if strings.LastIndex(result, ".") > 0 {
			suffix := result[strings.LastIndex(result, "."):]
			path = strings.Replace(path, "{suff}", suffix, -1)
		}else{
			path = strings.Replace(path, "{suff}", "", -1)
		}

	}

	c.Cmd = strings.Replace(c.Cmd, "{path}", path, -1)
	cmds := strings.Split(c.Cmd, " ")

	var bytes []byte
	if bytes, err = ioutil.ReadFile(path); err != nil {
		//not exist
		switch cmds[0] {
		case CMD[0],CMD[2]:
			res, err := http.Get(result)
			if err != nil {
				log.Println(err)
			}
			_, err = os.Create(path)
			if err != nil {
				log.Println(err)
			}
			if bytes,err = ioutil.ReadAll(res.Body); err != nil {
				return err
			}
			ioutil.WriteFile(path, bytes, 0666)
			break
		case CMD[1]:
			f, err := os.Create(path)
			if err != nil {
				log.Println(err)
			}
			bytes = []byte(result)
			io.WriteString(f, result)
			break
		}
	}

	// command
	switch cmds[0] {
	case CMD[0]:
		//play
		err = player.Play(path)
		if err != nil {
			log.Print(err)
			os.Exit(-1)
		}
	case CMD[1]:
		//text
		fmt.Println(result)
	case CMD[2]:
		//open
		var c *exec.Cmd
		if len(cmds) == 2 {
			c = exec.Command(cmds[1])
		}else if len(cmds) > 2{
			c = exec.Command(cmds[1], cmds[2:]...)
		}else{
			return errors.New("open command error")
		}
		_, err := c.CombinedOutput()
		if err != nil {
			return err
		}
	}
	// play

	return nil
}

func main(){
	var scheme,word string
	if(len(os.Args) == 2){
		scheme = "default"
		word = os.Args[1]
	}else if(len(os.Args) == 3){
		scheme  = os.Args[1]
		word = os.Args[2]
	}

	c,ok := configs[scheme]
	if !ok {
		log.Fatal("Unknow scheme for " + scheme)
		return
	}

	if !c.Enable {
		log.Printf("the %s is not enable or dir uncreated", scheme)
	}

	err := c.exec(word)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
}
