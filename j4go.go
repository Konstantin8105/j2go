package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	flag.Parse()

	files := flag.Args()

	if len(files) == 0 {
		flag.Usage()
		return
	}

	for _, file := range files {
		if err := j4go(file); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v", err)
		}
	}
}

func j4go(filename string) (err error) {

	var data string
	{
		dat, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}
		data = string(dat)
	}
	data = strings.Replace(data, "\r", "", -1)
	data = strings.Replace(data, "'", "\"", -1)
	data = strings.ToLower(data)

	for {
		index := strings.Index(data, "<script")
		if index < 0 {
			break
		}
		lastIndex := strings.Index(data, "</script>")
		if lastIndex < 0 {
			break
		}
		script := data[index:lastIndex]
		data = data[lastIndex+6:]

		// remove : <script ... >
		if index := strings.Index(script, "<script"); index >= 0 {
			if lastIndex := strings.Index(script, ">"); lastIndex >= 0 {
				script = script[lastIndex+1:]
			}
		}
		// trim
		script = strings.TrimSpace(script)
		if len(script) == 0 {
			continue
		}
		// add new lines
		var parenCounter int
		bb := []byte(script)
		for i := 0; i < len(bb); i++ {
			if bb[i] == '(' {
				parenCounter++
				continue
			}
			if bb[i] == ')' {
				parenCounter--
				continue
			}
			if parenCounter == 0 && bb[i] == ';' {
				bb[i] = '\n'
				continue
			}
		}
		script = string(bb)
		// func
		script = strings.Replace(script, "math.exp", "math.Exp", -1)
		script = strings.Replace(script, "math.round", "math.Round", -1)
		script = strings.Replace(script, "math.log", "math.Log", -1)
		script = strings.Replace(script, "math.pow", "math.Pow", -1)
		script = strings.Replace(script, "math.abs", "math.Abs", -1)
		script = strings.Replace(script, "function ", "\nfunc ", -1)
		script = strings.Replace(script, "type", "typeT", -1)
		script = strings.Replace(script, "{", "{\n", -1)
		script = strings.Replace(script, "}", "}\n", -1)
		script = strings.Replace(script, "<td>", "", -1)
		script = strings.Replace(script, "</td>", "", -1)
		script = strings.Replace(script, "<tr>", "", -1)
		script = strings.Replace(script, "</tr>", "", -1)
		script = strings.Replace(script, "<b>", "", -1)
		script = strings.Replace(script, "</b>", "", -1)
		script = strings.Replace(script, "<font>", "", -1)
		script = strings.Replace(script, "</font>", "", -1)
		script = strings.Replace(script, "<center>", "", -1)
		script = strings.Replace(script, "</center>", "", -1)
		script = strings.Replace(script, "<h1>", "", -1)
		script = strings.Replace(script, "</h1>", "", -1)
		script = strings.Replace(script, "<h2>", "", -1)
		script = strings.Replace(script, "</h2>", "", -1)
		script = strings.Replace(script, "<h3>", "", -1)
		script = strings.Replace(script, "</h3>", "", -1)
		script = strings.Replace(script, "<head>", "", -1)
		script = strings.Replace(script, "</head>", "", -1)
		script = strings.Replace(script, "<title>", "", -1)
		script = strings.Replace(script, "</title>", "", -1)
		script = strings.Replace(script, "<table>", "", -1)
		script = strings.Replace(script, "</table>", "", -1)
		script = strings.Replace(script, "<body>", "", -1)
		script = strings.Replace(script, "</body>", "", -1)
		script = strings.Replace(script, "<html>", "", -1)
		script = strings.Replace(script, "</html>", "", -1)

		//var mw=new array(28.016,32,44.01,18.016,39.95,64.06,28.01)
		script = strings.Replace(script, "new array(", "[]float64{", -1)

		// for(j=1;j<=4;j++){
		{
			lines := strings.Split(script, "\n")
			for i := range lines {
				if strings.Contains(lines[i], "for(") {
					lines[i] = strings.Replace(lines[i], "for(", "for ", -1)
					lines[i] = strings.Replace(lines[i], ")", " ", 1)
				}
				if strings.Contains(lines[i], "if(") {
					lines[i] = strings.Replace(lines[i], "if(", "if ", -1)
					lines[i] = strings.Replace(lines[i], ")", " ", 1)
				}
				if strings.Contains(lines[i], "while(") {
					lines[i] = strings.Replace(lines[i], "while(", "for ", -1)
					lines[i] = strings.Replace(lines[i], ")", " ", 1)
				}
			}
			script = strings.Join(lines, "\n")
		}

		// print
		fmt.Println(script)
	}

	return nil
}
