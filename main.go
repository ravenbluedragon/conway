package main

import (
	"html/template"
	"math/rand"
	"net"
	"net/http"
)

const size = 50
const probability = 0.4

var lsn net.Listener
var t *template.Template
var board [][]bool

const tbl = `<table>{{range .}}<tr>{{range .}}<td{{if .}} class="on"{{end}}/>{{end}}</tr>{{end}}</table>`

var handlers = map[string]func(http.ResponseWriter, *http.Request){
	"/":       root,
	"/quit":   quit,
	"/random": random,
	"/step":   step,
}

func main() {
	t = template.New("table")
	t.Parse(tbl)
	for k, v := range handlers {
		http.HandleFunc(k, v)
	}
	lsn, _ = net.Listen("tcp", ":8080")
	http.Serve(lsn, nil)
}

func quit(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Shutting down..."))
	lsn.Close()
}

func random(w http.ResponseWriter, r *http.Request) {
	board = make([][]bool, size)
	for i := range board {
		r := make([]bool, size)
		board[i] = r
		for j := range r {
			if rand.Float32() < probability {
				r[j] = true
			}
		}
	}
	t.Execute(w, board)
}

func live(me bool, count int) bool {
	if count == 3 || count == 4 && me {
		return true
	}
	return false
	//Any live cell with fewer than two live neighbours dies, as if caused by under-population.
	//Any live cell with two or three live neighbours lives on to the next generation.
	//Any live cell with more than three live neighbours dies, as if by over-population.
	//Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.
}

func step(w http.ResponseWriter, r *http.Request) {
	count := make([][]int, size)
	for i, row := range board {
		count[i] = make([]int, size)
		for j := range row {
			c := 0
			if j > 0 && board[i][j-1] {
				c++
			}
			if board[i][j] {
				c++
			}
			if j < size-1 && board[i][j+1] {
				c++
			}
			count[i][j] = c
		}
	}

	for i, row := range count {
		for j := range row {
			c := 0
			if i > 0 {
				c += count[i-1][j]
			}
			c += count[i][j]
			if i < size-1 {
				c += count[i+1][j]
			}
			board[i][j] = live(board[i][j], c)
		}
	}

	t.Execute(w, board)
}

const html = `
<!DOCTYPE html>
<html>
<head>
	<title>Conway's Game of Life</title>
	<script>
		window.addEventListener('load', function() {
			var b = document.getElementById('board');
			var running;
			
			var q = document.getElementById('quit');
			q.addEventListener('click', function() {
				replace(b, "/quit");
				clearInterval(running);
			});
			
			var r = document.getElementById('random');
			r.addEventListener('click', function() {
				replace(b, "/random");
			});
			
			var s = document.getElementById('step');
			s.addEventListener('click', function() {
				replace(b, "/step");
			});
			
			s = document.getElementById('start');
			s.addEventListener('click', function() {
				running = setInterval(function(){
					replace(b, "/step");
				}, 200);
				replace(b, "/step");
			});
			
			s = document.getElementById('stop');
			s.addEventListener('click', function() {
				clearInterval(running);
			});
		});
		
		function replace(node, target) {
			var xhr = new XMLHttpRequest();
			xhr.open('GET', encodeURI(target));
			xhr.onload = function() {
				if (xhr.status === 200) {
					node.innerHTML = xhr.responseText;
				} else {
					alert('Request failed.  Returned status of ' + xhr.status);
				}
			};
			xhr.send();
		}
	</script>
	<style>
		td {
			border: solid black 1px;
			width: 6px;
			height: 6px;
		}
		td.on {
			background-color: black;
		}
		table {
			border-collapse: collapse;
			margin-top
		}
	</style>
	<link href="life.css" rel="stylesheet">
</head>
<body>

<h1>Conway's Game of Life</h1>
<button id="random">Random Seed</button>
<button id="start">Start</button>
<button id="step">Step</button>
<button id="stop">Stop</button>
<button id="quit">Quit</button>
<div id="board"/>
</body>
</html>
`

func root(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(html))
}
